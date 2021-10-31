package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func GenerateHash(values ...string) string {
	var sb strings.Builder
	for _, value := range values {
		sb.WriteString(value)
	}
	// Generate Hash and convert into hex string
	hash := sha256.Sum256([]byte(sb.String()))

	return hex.EncodeToString(hash[:])
}
