package kv

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

type Storer interface {
	PutUser(user, creds string) error
	GetUser(user string) []byte
	GetSession(token string) []byte
	PutSession(token, created string) error
	Close() error
}

type Storage struct {
	*bolt.DB
}

func NewDB(file string) Storer {
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

	return &Storage{db}
}

func (s *Storage) PutUser(user, creds string) error {
	err := s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(authBucket))
		err := b.Put([]byte(user), []byte(creds))
		return err
	})
	return err
}

func (s *Storage) GetUser(user string) []byte {
	var creds []byte
	_ = s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(authBucket))
		creds = b.Get([]byte(user))
		return nil
	})
	return creds
}

func (s *Storage) PutSession(token, created string) error {
	err := s.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionBucket))
		err := b.Put([]byte(token), []byte(created))
		return err
	})
	return err
}

func (s *Storage) GetSession(token string) []byte {
	var created []byte
	_ = s.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionBucket))
		created = b.Get([]byte(token))
		return nil
	})
	return created
}
