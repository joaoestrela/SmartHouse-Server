package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	server "github.com/freddygv/SmartHouse-Server/go"
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

	err := server.NewAuthDB(testdb)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}

	if _, err := os.Stat(testdb); os.IsNotExist(err) {
		t.Fatalf("failed to create db: %v", err)
	}
}

func TestAuth(t *testing.T) {
	testdb, teardown := setup(t)
	defer teardown()

	err := server.NewAuthDB(testdb)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}

	srv := httptest.NewServer(server.NewRouter())
	defer srv.Close()

	input := server.RegInput{
		Username: "Bob",
		Password: "password",
		Secret:   server.Secret,
	}

	buf, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("failed to marshal creds: %v", err)
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/SmartHouse/1.0.2/register", srv.URL),
		"application/json",
		bytes.NewBuffer(buf),
	)
	if err != nil {
		t.Fatalf("failed to post: %v", err)
	}
	defer resp.Body.Close()

	tt := []struct {
		desc      string
		loginName string
		loginPW   string
		err       string
		code      int
	}{
		{
			desc:      "happy path",
			loginName: input.Username,
			loginPW:   input.Password,
			code:      http.StatusOK,
		},
		{
			desc:      "wrong pw",
			loginName: input.Username,
			loginPW:   "notpassword",
			err:       "Login failed: incorrect password",
			code:      http.StatusBadRequest,
		},
		{
			desc:      "user does not exist",
			loginName: "Alice",
			loginPW:   "password",
			err:       "Login failed: unregistered user: Alice",
			code:      http.StatusBadRequest,
		},
	}

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			user := server.LoginInput{
				Username: tc.loginName,
				Password: tc.loginPW,
			}
			buf, err := json.Marshal(user)
			if err != nil {
				t.Fatalf("failed to marshal creds: %v", err)
			}

			resp, err := http.Post(
				fmt.Sprintf("%s/SmartHouse/1.0.2/login", srv.URL),
				"application/json",
				bytes.NewBuffer(buf),
			)
			if err != nil {
				t.Fatalf("failed to post: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.code {
				t.Errorf("expected status: %d, got: %d", tc.code, resp.StatusCode)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("failed to read body: %v", err)
			}

			var payload string
			if err := json.Unmarshal(
				body,
				&struct {
					Message *string `json:"message"`
				}{
					&payload,
				},
			); err != nil {
				t.Fatalf("failed to unmarshal response")
			}

			if tc.err != "" && payload != tc.err {
				t.Fatalf("unexpected error: '%v', expected: '%s'", err, tc.err)
			}
		})
	}
}
