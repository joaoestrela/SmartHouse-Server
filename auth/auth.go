package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

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
// TODO: Check if user is already registered first
func Register(db *AuthStore, user, pw string) error {
	salt, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("failed to generate uuid: %v", err)
	}
	key := hash(pw, salt)

	buf, err := json.Marshal(credential{Key: key, Salt: salt})
	if err != nil {
		return fmt.Errorf("failed to marshal (key:%s, salt: %s): %v", key, salt, err)
	}

	err = db.PutUser(user, string(buf))
	if err != nil {
		return fmt.Errorf("failed to put new user '%s': %v", user, err)
	}
	return nil
}

// Authenticate validates a username and password then returns a session token
func Authenticate(db *AuthStore, user, pw string) (token string, err error) {
	stored := db.UserCredentials(user)
	if stored == nil {
		return "", fmt.Errorf("unregistered user: %s", user)
	}

	var creds credential
	if err := json.Unmarshal(stored, &creds); err != nil {
		return "", fmt.Errorf("failed to unmarshal creds: %v", err)
	}

	key := hash(pw, creds.Salt)
	if key != creds.Key {
		return "", fmt.Errorf("incorrect password")
	}

	token, err = newSession(db)
	if err != nil || token == "" {
		return "", fmt.Errorf("failed to generate session token: %v", err)
	}

	return token, nil
}

func hash(pw, salt string) string {
	key := pbkdf2.Key([]byte(pw), []byte(salt), iterations, keyLength, sha256.New)
	return base64.StdEncoding.EncodeToString(key)

}

// newSession persists and returns a new session token
func newSession(db *AuthStore) (token string, err error) {
	for i := 0; i < maxRetries; i++ {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return "", fmt.Errorf("failed to generate uuid: %v", err)
		}

		created, err := db.SessionCreation(uuid)
		if err != nil {
			log.Printf("failed to parse creation. deleting token '%s': %v\n", uuid, err)
			if err = db.DeleteSession(uuid); err != nil {
				log.Printf("failed to delete token '%s': %v\n", uuid, err)
				continue
			}
		}
		// Continue to generate a new token if the current one exists
		if created != 0 {
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
