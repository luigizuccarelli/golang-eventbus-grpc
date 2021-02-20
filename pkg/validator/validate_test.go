package validator

import (
	"fmt"
	"os"
	"testing"

	"github.com/microlib/simple"
)

func TestEnvars(t *testing.T) {
	logger := &simple.Logger{Level: "trace"}

	t.Run("ValidateEnvars : should fail", func(t *testing.T) {
		os.Setenv("SERVER_PORT", "")
		err := ValidateEnvars(logger)
		if err == nil {
			t.Errorf(fmt.Sprintf("Handler %s returned with no error - got (%v) wanted (%v)", "ValidateEnvars", err, nil))
		}
	})

	t.Run("ValidateEnvars : should pass", func(t *testing.T) {
		os.Setenv("LOG_LEVEL", "info")
		os.Setenv("ADV_URL", "/test")
		os.Setenv("SERVER_PORT", "9000")
		os.Setenv("VERSION", "1.0.3")
		os.Setenv("NAME", "test")
		os.Setenv("RPCSERVER_PORT", ":7000")
		os.Setenv("GRPCSERVER_ADDRESS", "127.0.0.1:7000")
		err := ValidateEnvars(logger)
		if err != nil {
			t.Errorf(fmt.Sprintf("Handler %s returned with error - got (%v) wanted (%v)", "ValidateEnvars", err, nil))
		}
	})
}
