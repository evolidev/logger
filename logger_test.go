package logger

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/evolidev/filesystem"
	"github.com/stretchr/testify/assert"
)

func TestLogging(t *testing.T) {
	t.Run("should log message", func(t *testing.T) {
		// create pipe
		r, w, _ := os.Pipe()
		defer r.Close()
		defer w.Close()

		logger := NewLogger(&Config{
			Stdout: w,
		})

		logger.Info("test")

		// read from pipe
		buf := make([]byte, 1024)
		n, _ := r.Read(buf)

		// check if message is logged
		assert.True(t, strings.Contains(string(buf[:n]), "test"), "message should be logged")
	})

	t.Run("should log message with prefix", func(t *testing.T) {
		// create pipe
		r, w, _ := os.Pipe()
		defer r.Close()
		defer w.Close()

		logger := NewLogger(&Config{
			Stdout: w,
			Name:   "prefix",
		})

		logger.Debug("test")
		logger.Error("test")
		logger.Success("test")
		logger.Log("test")
		logger.Info("test")

		// read from pipe with timeout
		buf := make([]byte, 1024)
		timeout := time.After(5 * time.Second)

		var result string

		select {
		case <-timeout:
			// handle timeout
			fmt.Println("Read operation timed out")
		case <-time.After(1 * time.Second):
			// read from pipe
			n, _ := r.Read(buf)
			// process the read data
			result = string(buf[:n])
			fmt.Println(result)
		}

		// check if message is logged
		assert.True(t, strings.Contains(result, "prefix"), "message should be logged with prefix")
	})

	t.Run("should log to a file", func(t *testing.T) {
		logger := NewLogger(&Config{
			Stdout: os.Stdout,
			Path:   "test.log",
		})

		logger.Info("test loggger")

		// check if message is logged to file
		assert.FileExists(t, "test.log", "message should be logged to file")

		content := filesystem.Read("test.log")
		assert.True(t, strings.Contains(content, "test loggger"), "message should be logged to file")

		// remove file
		defer os.Remove("test.log")
	})
}
