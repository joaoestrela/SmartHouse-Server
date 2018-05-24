package server

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	uuid "github.com/hashicorp/go-uuid"
	"golang.org/x/crypto/pbkdf2"
)

const (
	maxRetries        = 6
	iterations        = 4096
	keyLength         = 64
	expirationSeconds = 60 * 60 * 24 * 7 // 7 days
	secret            = "esperta"
)

// Login validates a username and password then returns a session token
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var in loginInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		msg := fmt.Sprintf("failed to decode request: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Login failed: %s"}`, msg)))
		return
	}

	stored := db.UserCredentials(in.Username)
	if stored == nil {
		msg := fmt.Sprintf("unregistered user: %s", in.Username)
		log.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "Login failed: %s"}`, msg)))
		return
	}

	var creds credential
	if err := json.Unmarshal(stored, &creds); err != nil {
		msg := fmt.Sprintf("failed to unmarshal json: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Login failed: %s"}`, msg)))
		return
	}

	key := hash(in.Password, creds.Salt)
	if key != creds.Key {
		msg := fmt.Sprintf("incorrect password")
		log.Println(msg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"message": "Login failed: %s"}`, msg)))
		return
	}

	token, err := newSession()
	if err != nil || token == "" {
		msg := fmt.Sprintf("failed to generate session token: %v", err)
		log.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message": "Login failed: %s"}`, msg)))
		return
	}

	buf, err := json.Marshal(StatusResponse{Message: token})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

// Register registers a new house member
// TODO: Check if user is already registered
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var in regInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		log.Printf("failed to decode req: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Registration failed."))
		return
	}

	salt, err := uuid.GenerateUUID()
	if err != nil {
		log.Printf("failed to generate uuid: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Registration failed."))
		return
	}
	key := hash(in.Password, salt)

	buf, err := json.Marshal(credential{Key: key, Salt: salt})
	if err != nil {
		log.Printf("failed to marshal (key:%s, salt: %s): %v\n", key, salt, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Registration failed."))
		return
	}

	err = db.PutUser(in.Username, string(buf))
	if err != nil {
		log.Printf("failed to put new user '%s': %v\n", in.Username, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Registration failed."))
		return
	}

	buf, err = json.Marshal(StatusResponse{Message: "OK"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(buf)
}

func hash(pw, salt string) string {
	key := pbkdf2.Key([]byte(pw), []byte(salt), iterations, keyLength, sha256.New)
	return base64.StdEncoding.EncodeToString(key)

}

// newSession persists and returns a new session token
func newSession() (token string, err error) {
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

type regInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Secret   string `json:"secret"`
}

type loginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type credential struct {
	Key  string
	Salt string
}

type session struct {
	Token   string
	Created int64
}
