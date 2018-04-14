package store

import (
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
)

const (
	authBucket    = "auth"
	sessionBucket = "sessions"
)

// AuthStorer enables the registration and validation of users and sessions
type AuthStorer interface {
	PutUser(user, creds string) error
	GetUser(user string) []byte
	PutSession(token, created string) error
	GetSession(token string) []byte
	Close() error
}

type kvStorage struct {
	*bolt.DB
}

// NewAuthDB returns a new and initialized db
func NewAuthDB(file string) AuthStorer {
	// Start fresh every time for now
	err := os.RemoveAll(file)
	if err != nil {
		log.Fatalf("failed to delete existing db: %v", err)
	}

	db, err := bolt.Open(file, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("failed to open new db: %v", err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(authBucket))
		_, err = tx.CreateBucketIfNotExists([]byte(sessionBucket))
		if err != nil {
			log.Fatalf("failed to create bucket: %v", err)
		}
		return nil
	})

	return &kvStorage{db}
}

func (s *kvStorage) PutUser(user, creds string) error {
	err := s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(authBucket))
		err := b.Put([]byte(user), []byte(creds))
		return err
	})
	return err
}

func (s *kvStorage) GetUser(user string) []byte {
	var creds []byte
	_ = s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(authBucket))
		creds = b.Get([]byte(user))
		return nil
	})
	return creds
}

func (s *kvStorage) PutSession(token, created string) error {
	err := s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionBucket))
		err := b.Put([]byte(token), []byte(created))
		return err
	})
	return err
}

func (s *kvStorage) GetSession(token string) []byte {
	var created []byte
	_ = s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionBucket))
		created = b.Get([]byte(token))
		return nil
	})
	return created
}
