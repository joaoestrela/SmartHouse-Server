package main

import (
	"testing"

	"github.com/boltdb/bolt"
)

const testFile = "test.db"

func TestBuildDB(t *testing.T) {
	db := buildDB(testFile)

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

func TestEndToEnd(t *testing.T) {
	tt := []struct {
		desc      string
		loginName string
		loginPW   string
		err       string
	}{
		{
			desc:      "happy path",
			loginName: "Bob",
			loginPW:   "password",
		},
		{
			desc:      "wrong pw",
			loginName: "Bob",
			loginPW:   "notpassword",
			err:       "incorrect password",
		},
		{
			desc:      "wrong pw",
			loginName: "Alice",
			loginPW:   "password",
			err:       "failed to get: unregistered user: Alice",
		},
	}

	db := buildDB(testFile)
	Register(db, "Bob", "password")

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := Authenticate(db, tc.loginName, tc.loginPW)
			if tc.err != "" && err.Error() != tc.err {
				t.Fatalf("unexpected error: '%v', expected: '%s'", err, tc.err)
			}
		})
	}
}
