package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/freddygv/SmartHouse-Server/auth/store"
	"github.com/hashicorp/go-uuid"
	"golang.org/x/crypto/pbkdf2"
)

const (
	maxRetries        = 6
	iterations        = 4096
	keyLength         = 64
	storage           = "auth.db"
	expirationSeconds = 60 * 60 * 24 // 1 day
)

// Register registers a new house member
func Register(db *store.AuthStore, user, pw string) {
	salt, err := uuid.GenerateUUID()
	if err != nil {
		log.Fatalf("failed to generate uuid: %v", err)
	}
	key := hash(pw, salt)

	buf, err := json.Marshal(credential{Key: key, Salt: salt})
	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}

	err = db.PutUser(user, string(buf))
	if err != nil {
		log.Fatalf("failed to put new user: %v", err)
	}
}

// Authenticate validates a username and password then returns a session token
func Authenticate(db *store.AuthStore, user, pw string) (token string, err error) {
	stored := db.UserCredentials(user)
	if stored == nil {
		return "", fmt.Errorf("unregistered user: %s", user)
	}

	var creds credential
	if err := json.Unmarshal(stored, &creds); err != nil {
		return "", fmt.Errorf("failed to get: %v", err)
	}

	key := hash(pw, creds.Salt)
	if key != creds.Key {
		return "", fmt.Errorf("incorrect password")
	}

	token, err = newSession(db)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

func hash(pw, salt string) string {
	key := pbkdf2.Key([]byte(pw), []byte(salt), iterations, keyLength, sha256.New)
	return base64.StdEncoding.EncodeToString(key)

}

// newSession persists and returns a new session token
func newSession(db *store.AuthStore) (token string, err error) {
	for i := 0; i < maxRetries; i++ {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return "", fmt.Errorf("failed to generate uuid: %v", err)
		}

		if created := db.SessionCreation(uuid); len(created) != 0 {
			continue
		}

		time := strconv.FormatInt(time.Now().Unix(), 10)
		err = db.PutSession(uuid, time)
		if err == nil {
			token = uuid
			break
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
