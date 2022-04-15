package v2

import (
	"crypto/rand"
	"crypto/subtle"
	"sync"
	"time"
)

type ApiKey struct {
	k         []byte
	expiresAt time.Time
}

type ApiKeyManager struct {
	userKeyMap map[string]*ApiKey

	eventChan chan time.Time
	lock      sync.Mutex
}

func NewApiKey(len int, d time.Duration) *ApiKey {
	return &ApiKey{
		k:         randBytes(len),
		expiresAt: time.Now().Add(d),
	}
}

func (ak *ApiKey) Key() string {
	return string(ak.k)
}

func NewApiKeyManager() *ApiKeyManager {
	return &ApiKeyManager{
		userKeyMap: make(map[string]*ApiKey),
	}
}

func (akm *ApiKeyManager) RegisterKey(uid string, key *ApiKey) {
	akm.lock.Lock()
	defer akm.lock.Unlock()

	// If the key already exists and hasn't been revoked, return
	if _, ok := akm.userKeyMap[uid]; ok {
		if time.Now().Before(akm.userKeyMap[uid].expiresAt) {
			return
		}
	}

	akm.userKeyMap[uid] = key

	timer := time.NewTimer(key.expiresAt.Sub(time.Now()))
	go func() {
		<-timer.C
		akm.RevokeKey(uid)
	}()
}

func (akm *ApiKeyManager) RetrieveKey(uid string) *ApiKey {
	akm.lock.Lock()
	defer akm.lock.Unlock()

	if v, ok := akm.userKeyMap[uid]; ok {
		return v
	}

	return nil
}

func (akm *ApiKeyManager) LookupUser(key string) string {
	var user string

	for k, v := range akm.userKeyMap {
		if subtle.ConstantTimeCompare(v.k, []byte(key)) == 1 {
			user = k
		}
	}

	return user
}

// TODO we probably will need to cancel the timer thread
func (akm *ApiKeyManager) RevokeKey(uid string) {
	akm.lock.Lock()
	defer akm.lock.Unlock()

	if _, ok := akm.userKeyMap[uid]; ok {
		delete(akm.userKeyMap, uid)
	}
}

func randBytes(length int) []byte {
	const validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)

	for index, rbyte := range bytes {
		bytes[index] = validChars[rbyte%byte(len(validChars))]
	}

	return bytes
}
