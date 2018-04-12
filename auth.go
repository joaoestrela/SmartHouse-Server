package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/go-uuid"
	"golang.org/x/crypto/pbkdf2"
)

const (
	iterations        = 4096
	keyLength         = 64
	authBucket        = "auth"
	sessionBucket     = "sessions"
	storage           = "auth.db"
	expirationSeconds = 60 * 60 * 24 // 1 day
)

func main() {
	db := buildDB()
	defer db.Close()

	user := "Bob"
	pw := "p4$$w0rd"

	Register(db, user, pw)
	t := Login(db, user, pw)
	fmt.Println(t)
}

// TODO: Move error handling to responsewriter
func Register(db *bolt.DB, user, pw string) {
	salt, err := uuid.GenerateUUID()
	if err != nil {
		log.Fatalf("failed to generate uuid: %v", err)
	}
	key := hash(pw, salt)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(authBucket))
		buf, err := json.Marshal(credential{Key: key, Salt: salt})
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
func Login(db *bolt.DB, user, pw string) (token string) {
	var creds credential

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(authBucket))
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

	key := hash(pw, creds.Salt)
	if key != creds.Key {
		log.Fatalf("incorrect password")
	}

	token, err = newSession(db)
	if err != nil {
		log.Fatalf("failed to generate token: %v", err)
	}

	return token
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
		_, err := tx.CreateBucketIfNotExists([]byte(authBucket))
		_, err = tx.CreateBucketIfNotExists([]byte(sessionBucket))
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

// newSession persists and returns a new session token
func newSession(db *bolt.DB) (token string, err error) {
	invalid := true
	for invalid {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return "", fmt.Errorf("failed to generate uuid: %v", err)
		}

		err = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(sessionBucket))
			v := b.Get([]byte(token))
			if len(v) != 0 {
				return errors.New("token exists")
			}

			time := strconv.FormatInt(time.Now().Unix(), 10)
			err = b.Put([]byte(uuid), []byte(time))
			return err
		})
		if err == nil {
			token = uuid
			invalid = false
		}
	}

	return token, nil
}

type credential struct {
	Key  string
	Salt string
}

type session struct {
	Token   string
	Created int64
}
