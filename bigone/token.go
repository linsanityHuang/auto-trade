package bigone

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// SignAuthenticationToken sign token
func SignAuthenticationToken(key, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"type":        "OpenAPIV2",
		"sub":         key,
		"nonce":       strconv.FormatInt(time.Now().UnixNano(), 10),
		"recv_window": "50",
	})
	block := []byte(secret)
	return token.SignedString(block)
}
