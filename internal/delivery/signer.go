package delivery

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// SignPayload generates HMAC-SHA256 signature for webhook payload
func SignPayload(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	signature := hex.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("sha256=%s", signature)
}

// VerifySignature verifies the HMAC-SHA256 signature
func VerifySignature(payload []byte, signature, secret string) bool {
	expectedSignature := SignPayload(payload, secret)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
