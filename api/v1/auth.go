package v1

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

type ApiKeyProperties struct {
	UID       string
	ExpiresAt time.Time
	cancel    context.CancelFunc
}

type ApiKeyManager struct {
	// m contains the apikey string as the key and the properties as the value
	m map[string]ApiKeyProperties

	eventChan chan time.Time
	lock      sync.Mutex
}

func NewApiKey(len int, uid string, d time.Duration) (string, ApiKeyProperties) {
	key := string(randBytes(len))
	properties := ApiKeyProperties{
		UID:       uid,
		ExpiresAt: time.Now().Add(d),
	}

	return key, properties
}

func NewApiKeyManager() *ApiKeyManager {
	return &ApiKeyManager{
		m: make(map[string]ApiKeyProperties),
	}
}

// RegisterKey
func (akm *ApiKeyManager) RegisterKey(key string, properties ApiKeyProperties) error {
	akm.lock.Lock()
	defer akm.lock.Unlock()

	// Check the passed in key for expiration
	if time.Now().After(properties.ExpiresAt) {
		return fmt.Errorf("failed to register expired key")
	}

	// Check for an already existing key and that it is unexpired
	if _, ok := akm.m[key]; ok && time.Now().After(akm.m[key].ExpiresAt) {
		return nil
	}

	// Add a cancellation channel to the properties
	ctx, cancel := context.WithCancel(context.Background())
	properties.cancel = cancel

	// Create the new key
	akm.m[key] = properties

	// Queue up a revoke of the key at the expiration time
	timer := time.NewTimer(properties.ExpiresAt.Sub(time.Now()))
	go func() {
		select {
		case <-timer.C:
			akm.RevokeKey(key)
		case <-ctx.Done():
			return
		}
	}()

	return nil
}

func (akm *ApiKeyManager) RetrieveKey(key string) (ApiKeyProperties, bool) {
	akm.lock.Lock()
	defer akm.lock.Unlock()

	if v, ok := akm.m[key]; ok {
		return v, true
	}

	return ApiKeyProperties{}, false
}

func (akm *ApiKeyManager) RevokeKey(key string) {
	akm.lock.Lock()
	defer akm.lock.Unlock()

	if _, ok := akm.m[key]; ok {
		akm.m[key].cancel()
		delete(akm.m, key)
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
