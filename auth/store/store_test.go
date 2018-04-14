package store

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func setup(t *testing.T) (string, func()) {
	t.Parallel()

	const testdb = "test.db"

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	teardown := func() { os.RemoveAll(dir) }

	return filepath.Join(dir, testdb), teardown
}

func TestNewDB(t *testing.T) {
	testdb, teardown := setup(t)
	defer teardown()

	_ = NewAuthDB(testdb)

	if _, err := os.Stat(testdb); os.IsNotExist(err) {
		t.Fatalf("failed to create db: %v", err)
	}
}
