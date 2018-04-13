package main

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
	db := buildDB(testdb)

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
	teardown()
}

func TestEndToEnd(t *testing.T) {
	testdb, teardown := setup(t)
	db := buildDB(testdb)

	user := "Bob"
	pw := "password"

	Register(db, user, pw)

	tt := []struct {
		desc      string
		loginName string
		loginPW   string
		err       string
	}{
		{
			desc:      "happy path",
			loginName: user,
			loginPW:   pw,
		},
		{
			desc:      "wrong pw",
			loginName: "Bob",
			loginPW:   "notpassword",
			err:       "incorrect password",
		},
		{
			desc:      "user does not exist",
			loginName: "Alice",
			loginPW:   "password",
			err:       "failed to get: unregistered user: Alice",
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := Authenticate(db, tc.loginName, tc.loginPW)
			if tc.err != "" && err.Error() != tc.err {
				t.Fatalf("unexpected error: '%v', expected: '%s'", err, tc.err)
			}
		})
	}
	teardown()
}
