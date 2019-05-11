package auth

import (
	"errors"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/philip-bui/fruits-service/config"
	"github.com/rs/zerolog/log"
)

const (
	base = 32
	sub  = "sub"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

// GetClaimsFromJWT checks, parses and returns Claims from JWT.
func GetClaimsFromJWT(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		log.Error().Err(err).Str("token", tokenString).Msg("invalid signature")
		return nil, ErrInvalidToken
	}
	return token.Claims.(jwt.MapClaims), nil
}

// GetUserIDFromJWT checks, parses and retrieves UserID from JWT.
func GetUserIDFromJWT(tokenString string) (int64, error) {
	claims, err := GetClaimsFromJWT(tokenString)
	if err != nil {
		return 0, err
	}
	userID, ok := claims[sub].(string)
	if !ok {
		log.Error().Fields(claims).Msg("error getting userID from JWT")
		return 0, ErrInvalidToken
	}
	i, err := strconv.ParseInt(userID, base, 64)
	if err != nil {
		log.Error().Str("userID", userID).Msg("error parsing userID to int64")
		return 0, ErrInvalidToken
	}
	return i, nil
}

// keyFunc checks the tokens signing method and signature matching.
func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, ErrInvalidToken
	}
	return config.JWT, nil
}

// SignInToken creates a JWT Token with UserID
func SignInToken(userID int64) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		sub: strconv.FormatInt(userID, base),
	}).SignedString(config.JWT)
	if err != nil {
		log.Error().Err(err).Int64("userID", userID).Msg("error creating sign in token")
		return "", err
	}
	log.Info().Int64("userID", userID).Msg("created sign in token")
	return token, nil
}
