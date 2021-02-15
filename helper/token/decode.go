package token

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/xor"
)

// DecodeRootToken will decode the root token returned by the Vault API
// The algorithm was initially used in the generate root command
func DecodeRootToken(encoded, otp string, otpLength int) (string, error) {
	if otpLength == 0 {
		// Backwards compat
		tokenBytes, err := xor.XORBase64(encoded, otp)
		if err != nil {
			return "", fmt.Errorf("Error xoring token: %s", err)
		}

		uuidToken, err := uuid.FormatUUID(tokenBytes)
		if err != nil {
			return "", fmt.Errorf("Error formatting base64 token value: %s", err)
		}
		return strings.TrimSpace(uuidToken), nil
	}

	tokenBytes, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("Error decoding base64'd token: %v", err)
	}

	tokenBytes, err = xor.XORBytes(tokenBytes, []byte(otp))
	if err != nil {
		return "", fmt.Errorf("Error xoring token: %v", err)
	}
	return string(tokenBytes), nil
}
