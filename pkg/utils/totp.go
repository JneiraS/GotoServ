package utils

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	// "golang.org/x/tools/go/analysis/passes/printf"
)

func NormalizeSecret(secret string) (string, error) {
	s := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(secret), " ", ""))

	// base32
	if _, err := base32.StdEncoding.DecodeString(addBase32Padding(s)); err == nil {
		return addBase32Padding(s), nil
	}

	// hex -> base32
	if raw, err := hex.DecodeString(s); err == nil {
		return base32.StdEncoding.EncodeToString(raw), nil
	}

	// base64 -> base32
	if raw, err := base64.StdEncoding.DecodeString(addBase64Padding(s)); err == nil {
		return base32.StdEncoding.EncodeToString(raw), nil
	}

	// Si rien ne marche, erreur
	return "", errors.New("clé secrète invalide")
}

func GenerateCurrentTOTP(secret string) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", errors.New("SECRET_KEY vide : vérifiez la variable d'environnement ou le .env")
	}
	key, err := NormalizeSecret(secret)
	if err != nil {
		return "", err
	}
	return totp.GenerateCodeCustom(key, time.Now(), totp.ValidateOpts{
		Period:    60,
		Digits:    otp.DigitsEight,
		Algorithm: otp.AlgorithmSHA256,
	})
}

func ValidateTOTP(secret, code string) (bool, error) {
	key, err := NormalizeSecret(secret)
	if err != nil {
		return false, err
	}

	return totp.ValidateCustom(code, key, time.Now(), totp.ValidateOpts{
		Period:    60,
		Skew:      1,
		Digits:    otp.DigitsEight,
		Algorithm: otp.AlgorithmSHA256,
	})
}

func addBase32Padding(s string) string {
	if m := len(s) % 8; m != 0 {
		return s + strings.Repeat("=", 8-m)
	}
	return s
}

func addBase64Padding(s string) string {
	if m := len(s) % 4; m != 0 {
		return s + strings.Repeat("=", 4-m)
	}
	return s
}
