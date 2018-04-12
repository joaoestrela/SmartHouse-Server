package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-uuid"
	"golang.org/x/crypto/pbkdf2"
)

const (
	iterations = 4096
	keyLength  = 64
	bucket     = "auth"
	storage    = "auth.db"
)

func main() {
	db := buildDB()
	defer db.Close()

	user := "Bob"
	pw := "p4$$w0rd"

	Register(db, user, pw)
	Login(db, user, pw)
}

// TODO: Move error handling to responsewriter
func Register(db *bolt.DB, user, pw string) {
	salt, err := uuid.GenerateUUID()
	if err != nil {
		log.Fatalf("failed to generate uuid: %v", err)
	}
	token := hash(pw, salt)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		buf, err := json.Marshal(credentials{Token: token, Salt: salt})
		if err != nil {
			log.Fatalf("failed to marshal: %v", err)
		}

		err = b.Put([]byte(user), buf)
		return err
	})
	if err != nil {
		log.Fatalf("failed to put new user: %v", err)
	}
	fmt.Println("OK")
}

// TODO: Move error handling to responsewriter
func Login(db *bolt.DB, user, pw string) {
	var creds credentials

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		value := b.Get([]byte(user))
		if value == nil {
			log.Fatalf("unregistered user: %s", user)
		}
		err := json.Unmarshal(value, &creds)
		return err
	})
	if err != nil {
		log.Fatalf("failed to get: %v", err)
	}

	h := hash(pw, creds.Salt)
	if h != creds.Token {
		log.Fatalf("incorrect password")
	}

	fmt.Println("OK")
}

func buildDB() *bolt.DB {
	err := os.RemoveAll(storage)
	if err != nil {
		log.Fatalf("failed to delete existing db: %v", err)
	}

	db, err := bolt.Open(storage, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("failed to open new db: %v", err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			log.Fatalf("failed to create bucket: %v", err)
		}
		return nil
	})

	return db
}

func hash(pw, salt string) string {
	key := pbkdf2.Key([]byte(pw), []byte(salt), iterations, keyLength, sha256.New)
	return base64.StdEncoding.EncodeToString(key)

}

type credentials struct {
	Token string
	Salt  string
}
