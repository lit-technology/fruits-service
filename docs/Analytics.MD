# Fruitlytics

- PostgreSQL BRIN indexed by Timestamp for small to mid scale surveys.
- Time Histogram, data charts are supported.
- Premium customers move into BigData, allowing filtering, searching at the Petabyte level.

## Answer Database Design

There were 5 possible outcomes concluded to scale our data needs.

- NoSQL. Overhead of document-based storage. Extra storage due to additional indexing in case of Elasticsearch etc.
- EAV (Entity-Attribute-Value). In SQL, this introduces massive overhead of JOINS and per rows storage overhead, and for KV databases this introduces extra operational costs.
- Millions of tables. Unscalable due to every table creating a file in FS, introducing overhead to every SQL operation. Imagine searching 1 million files for your table metadata every query.
- JSONB. Simple design, extra query, write, storage overhead.
- Hundreds of columns. Storing null columns costs nothing, denormalizing means no joins, no ids = most performant in speed and storage. Huge operational costs.

We chose the hundreds of columns, primarily to save storage costs. This means there is an upper limit of 500 questions.

### Setup

We partitioned manually instead of using PostgreSQL native partitioning.

- 2048 partitions. Instagram uses 2048 shards as their most optimal performance.
- Performance. Skip the query partition planner. Insert directly into partition for more write performance.
- Debugging. PostgreSQL mutates IDs before partitioning. E.g ID=1 44th partition, ID=2 48th etc.

### Further reading

- https://stackoverflow.com/questions/13570613/making-sense-of-postgres-row-sizes
- https://stackoverflow.com/questions/870808/entity-attribute-value-database-vs-strict-relational-model-ecommerce
- https://www.compose.com/articles/faster-operations-with-the-jsonb-data-type-in-postgresql/
- https://heap.io/blog/engineering/when-to-avoid-jsonb-in-a-postgresql-schema

## PostgreSQL instead of MySQL OLAP

Several reasons were chosen for PostgreSQL instead of MySQL.

- BRIN Indexing. Known as block-range indexing, we only store sumarries of block ranges for the query planner to determine whether to query in the block. A B-Tree will store additional values in a tree index, increasing storage and write costs.
- Array support. Given we chose duplicate-column tables, arrays help reduce operational costs and querying costs.
- Incredibly powerful query functionality. CTE, Grouping Sets etc.
- Query planner is more optimized.

On the contrary, MySQL did have upsides.

- TinyInt support. Most of our answers are INT, we theoretically on average can save 1/6 of storage per row.
- RocksDB Storage Engine. Incredibly performant write and read from recent data.

### Costs

AWS

- $106 Small Instance Yearly
- $0.115 GB SDD
- $0.1 GB HDD

### Further reading

- https://info.crunchydata.com/blog/postgresql-brin-indexes-big-data-performance-with-minimal-storage

## Athena / S3-Parquet instead of BigQuery

AWS Athena was chosen due to the open-source community and fine-tuning opportunities.

BigQuery

- Simplicity. Insert answers as encountered. Simple to create partitioned and clustering through SQL DDL. Arrays and Map is simple.
- Automatic cold-storage pricing. BigQuery data not accessed in 3 months downgrades pricing from $0.02 to $0.01.

Athena

- Message queue complexity. To reduce pricing, it is recommended to store as many rows (partitioned) in one file, reducing S3 GET and PUT operations.
- Fine-tuning INTs. BigQuery stores 8 Bytes instead of 1 Byte in Parquet.
- Fine-tuning Clustering. BigQuery only supports clustering RANGE(survey_id), whereas S3 only query only one file.
- Fine-tuning Partitioning. BigQuerry only supports partitioning RANGE(timestamp), whereas S3 can sub-divide into year, month, day folders.
- Simple to upgrade to Presto to save future costs, or query against multiple storage locations.

### Costs

AWS

- $0.019 GB Infrequent Access
- $0.01 1K Puts
- $0.001 1K Gets
- Data Transfers between S3 buckets or from Amazon S3 to any service(s) within the same AWS Region are free.
- Athena Lambda costs.

Note: [Snappy](https://github.com/google/snappy) compresses our data in S3 by roughly 50%.

### Snappy vs gzip

Snappy is fast (250MB/s), ideal for decompressing large amounts of data as opposed to gzip. Snappy also supports Parquet column splitting.

- https://www.cloudera.com/documentation/enterprise/5-3-x/topics/admin_data_compression_performance.html

### Further Reading

- https://github.com/apache/parquet-format/blob/master/LogicalTypes.md
