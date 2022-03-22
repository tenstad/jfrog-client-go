package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tenstad/jfrog-client-go/utils/io/content"

	"github.com/tenstad/jfrog-client-go/utils/log"
)

func init() {
	log.SetLogger(log.NewLogger(log.DEBUG, nil))
}

func getBaseTestDir(t *testing.T) string {
	pwd, err := os.Getwd()
	assert.NoError(t, err, "Failed to get current dir")
	if err != nil {
		return ""
	}
	return filepath.Join(pwd, "tests", "testdata")
}

func readerCloseAndAssert(t *testing.T, reader *content.ContentReader) {
	assert.NoError(t, reader.Close(), "Couldn't close reader")
}
