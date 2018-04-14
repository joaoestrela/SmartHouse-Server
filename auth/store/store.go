package store

import (
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

const (
	authBucket    = "auth"
	sessionBucket = "sessions"
)

// AuthStore supports the creation and retrieval of users and sessions
type AuthStore struct {
	*bolt.DB
}

// NewAuthDB returns a new and initialized db
func NewAuthDB(file string) (*AuthStore, error) {
	// Start fresh every time for now
	err := os.RemoveAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to delete existing db: %v", err)
	}

	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open new db: %v", err)
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(authBucket))
		_, err = tx.CreateBucketIfNotExists([]byte(sessionBucket))
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &AuthStore{db}, nil
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
func (s *AuthStore) SessionCreation(token string) []byte {
	var created []byte
	_ = s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionBucket))
		created = b.Get([]byte(token))
		return nil
	})
	return created
}
