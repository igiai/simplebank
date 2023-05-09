package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// here we configure gin for testing for the tests not to be overloaded with information
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
