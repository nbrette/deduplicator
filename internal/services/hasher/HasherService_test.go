package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateHash(t *testing.T) {
	payload := "0123456"
	deviceID := "AZERTY"
	concat := deviceID + payload

	hash := sha256.Sum256([]byte(concat))
	hashHex := hex.EncodeToString(hash[:])
	var deduplicatorAttributes []string
	deduplicatorAttributes = append(deduplicatorAttributes, deviceID)
	deduplicatorAttributes = append(deduplicatorAttributes, payload)
	hashGenerated := GenerateHash(deduplicatorAttributes...)
	assert.Equal(t, hashHex, hashGenerated)

	deduplicatorAttributes[0] = payload
	deduplicatorAttributes[1] = deviceID
	hashGenerated = GenerateHash(deduplicatorAttributes...)
	assert.NotEqual(t, hashHex, hashGenerated)
}
