package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker creates a new JWTMaker
// Here we return Maker interface not the JWTMaker struct itself to ensure that
// JWTMaker implements Maker interface
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken creates new token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	// keyFunc receives parsed (split) but unverified token
	// It is responsible for providing a Parse method with a key to validate token
	// In our case, before providing the key we want to check whether or not the signing algorithm set in the header of
	// the token is the appropriate algorithm, used in a system to sign tokens
	// This prevents from someone hacking token by enforcing system to use algorithm the hacker wants
	// If it matches, keyFunc will return a key that can be used to verify the token
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// check if the algorithm specified in the token matches the one used in our system
		// we are checking for SigningMethodHMAC because SigningMethodHS256 used in a system to sign tokens is an instance of SigningMethodHMAC
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	// Parsing is a process of extracting information from a string, so in this case a token string is split into 3 parts
	// as JWT token is composed of 3 parts, that is already parsing and then information will be extracted from these parts
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		// Inside the ParseWithClaims method a Valid method on Payload is called and if it returns an error it is embedded in
		// jwt.ValidationError and not returned explicitly, so we have to extract it here
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
