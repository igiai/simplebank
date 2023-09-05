package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/igiai/simplebank/db/sqlc"
	"github.com/igiai/simplebank/db/util"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	// here we configure gin for testing for the tests not to be overloaded with information
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
