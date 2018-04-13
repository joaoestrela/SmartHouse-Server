package kv

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/boltdb/bolt"
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

func TestBuildDB(t *testing.T) {
	testdb, teardown := setup(t)
	defer teardown()

	db := NewDB(testdb)
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(authBucket))
		if err.Error() != "bucket already exists" {
			t.Fatalf("bucket '%s' not created", authBucket)
		}

		_, err = tx.CreateBucket([]byte(sessionBucket))
		if err.Error() != "bucket already exists" {
			t.Fatalf("bucket '%s' not created", sessionBucket)
		}

		b := tx.Bucket([]byte(authBucket))
		err = b.Put([]byte("user"), []byte("password"))
		if err != nil && err.Error() == "database not open" {
			t.Fatal(err)
		}
		return nil
	})
}
