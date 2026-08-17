package mfa

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"image"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TOTPSecret struct {
	Secret string
	URL    string
}

func GenerateTOTPSecret() (*TOTPSecret, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "HinduRitualPlatform",
		AccountName: "user",
		Period:      30,
		Digits:      6,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	return &TOTPSecret{
		Secret: key.Secret(),
		URL:    key.URL(),
	}, nil
}

func ValidateTOTPCode(secret, code string) bool {
	secret = strings.ToUpper(secret)
	if !strings.HasSuffix(secret, "=") {
		secret = padBase32(secret)
	}
	valid, err := totp.ValidateCustom(
		code,
		secret,
		time.Now(),
		totp.ValidateOpts{
			Period:    30,
			Digits:    6,
			Algorithm: otp.AlgorithmSHA1,
		},
	)
	if err != nil {
		return false
	}
	return valid
}

func padBase32(s string) string {
	padding := 8 - len(s)%8
	if padding == 8 {
		return s
	}
	return s + strings.Repeat("=", padding)
}

func GenerateRecoveryCodes() ([]string, error) {
	codes := make([]string, 10)
	for i := 0; i < 10; i++ {
		bytes := make([]byte, 8)
		if _, err := rand.Read(bytes); err != nil {
			return nil, fmt.Errorf("failed to generate recovery code: %w", err)
		}
		code := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)
		codes[i] = fmt.Sprintf("%s-%s-%s-%s", code[0:4], code[4:8], code[8:12], code[12:16])
	}
	return codes, nil
}

func GetTOTPQRCode(secret, email string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "HinduRitualPlatform",
		AccountName: email,
		Secret:      []byte(secret),
		Period:      30,
		Digits:      6,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP URL: %w", err)
	}

	img, err := totpImage(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	return img, nil
}

func totpImage(key *otp.Key) (string, error) {
	img, err := key.Image(200, 200)
	if err != nil {
		return "", fmt.Errorf("failed to create QR image: %w", err)
	}

	return imageToBase64(img)
}

func imageToBase64(img image.Image) (string, error) {
	return "", fmt.Errorf("QR code generation requires external library: install github.com/skip2/go-qrcode or use the key.URL() directly")
}

func GenerateTOTPWithOpts(email, issuer string) (*TOTPSecret, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: email,
		Period:      30,
		Digits:      6,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP: %w", err)
	}

	return &TOTPSecret{
		Secret: key.Secret(),
		URL:    key.URL(),
	}, nil
}
