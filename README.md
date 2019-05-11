# Fruits Service

## Setup

Run `make` to get a list of commands.

- `make setup` performs downloading dependencies, generating files and running Postgres on Docker.
- `make godoc` creates a HTML API documentation.
- `make mod` installs dependencies.
- `make proto` installs [gogoprotobuf](https://github.com/gogo/protobuf) and compiles Protobuf files.
- `make json` installs [gojay](https://github.com/francoispqt/gojay) and generates JSON Serializers and Deserializers.
- `make postgres` installs and runs Postgres 11 Docker on port 5437.
- `make db` creates various Postgres databases.
- `make tables` creates various Postgres tables.
- `make coverage` creates a HTML test coverage report.

## Deployment

- `make deploy` compiles an executable for Linux 64-bit, zips and deploys through ElasticBeanstalk.

## License

```
Copyright (©) Philip Bui - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by Philip Bui <philip.bui.developer@gmail.com>
```
