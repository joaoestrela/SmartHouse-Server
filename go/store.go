package server

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

// AuthStore supports the creation and retrieval of user and session data
type AuthStore struct {
	*bolt.DB
}

const (
	authBucket    = "auth"
	sessionBucket = "sessions"
)

var db *AuthStore

// NewAuthDB returns a new and initialized db
func NewAuthDB(file string) error {
	// Start fresh every time for now
	err := os.RemoveAll(file)
	if err != nil {
		return fmt.Errorf("failed to delete existing db: %v", err)
	}

	storage, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("failed to open new db: %v", err)
	}

	if err := storage.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(authBucket))
		_, err = tx.CreateBucketIfNotExists([]byte(sessionBucket))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
		return nil
	}); err != nil {
		return err
	}

	db = &AuthStore{storage}

	return nil
}

// PutUser persists a user name and its credentials (key and salt)
func (s *AuthStore) PutUser(user, creds string) error {
	err := s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(authBucket))
		err := b.Put([]byte(user), []byte(creds))
		return err
	})
	return err
}

// UserCredentials retrieves a user's credentials (key and salt)
func (s *AuthStore) UserCredentials(user string) []byte {
	var creds []byte
	_ = s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(authBucket))
		creds = b.Get([]byte(user))
		return nil
	})
	return creds
}

// PutSession persists a session token and its creation date
func (s *AuthStore) PutSession(token, created string) error {
	err := s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionBucket))
		err := b.Put([]byte(token), []byte(created))
		return err
	})
	return err
}

// SessionCreation retrieves the creation date of a session token
func (s *AuthStore) SessionCreation(token string) (int64, error) {
	var created int64
	if err := s.View(func(tx *bolt.Tx) (err error) {
		b := tx.Bucket([]byte(sessionBucket))
		c := b.Get([]byte(token))
		created, err = strconv.ParseInt(string(c), 10, 64)
		return err
	}); err != nil {
		return 0, err
	}
	return created, nil
}

// DeleteSession deletes a session token
func (s *AuthStore) DeleteSession(token string) error {
	if err := s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionBucket))
		return b.Delete([]byte(token))
	}); err != nil {
		return err
	}
	return nil
}
