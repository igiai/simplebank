package token

import (
	"testing"
	"time"

	"github.com/igiai/simplebank/db/util"
	"github.com/o1egl/paseto"
	"github.com/stretchr/testify/require"
)

// Test for paseto maker is almost identical as the one for JWT maker as they both implement Maker interface
func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t, err)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidPasetoToken(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	invalidKey := util.RandomString(32)
	validKey := util.RandomString(32)

	paseto := paseto.NewV2()
	token, err := paseto.Encrypt([]byte(invalidKey), payload, nil)
	require.NoError(t, err)

	maker, err := NewPasetoMaker(validKey)
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestInvalidKeyLengthForPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(16))
	require.Error(t, err)
	require.Nil(t, maker)
}
