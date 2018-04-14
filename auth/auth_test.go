package auth

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

	_, err := NewAuthDB(testdb)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}

	if _, err := os.Stat(testdb); os.IsNotExist(err) {
		t.Fatalf("failed to create db: %v", err)
	}
}

func TestEndToEnd(t *testing.T) {
	testdb, teardown := setup(t)
	defer teardown()

	user := "Bob"
	pw := "password"

	db, err := NewAuthDB(testdb)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
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
			err:       "unregistered user: Alice",
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
}
