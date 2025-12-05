package crypto

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/blake2b"
)

type ONDCCrypto struct{}

func NewONDCCrypto() *ONDCCrypto {
	return &ONDCCrypto{}
}

// SignRequest signs the request payload using ed25519
func (c *ONDCCrypto) SignRequest(privateKeyStr string, payload []byte, created int, ttl int) (string, error) {
	// Compute BLAKE2b-512 hash over the payload
	hash := blake2b.Sum512(payload)
	digest := base64.StdEncoding.EncodeToString(hash[:])

	// Create signature body
	expires := created + ttl
	signatureBody := fmt.Sprintf(
		"(created): %d\n(expires): %d\ndigest: BLAKE-512=%s",
		created,
		expires,
		digest,
	)

	// Decode private key
	decodedKey, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return "", fmt.Errorf("error decoding signing private key: %w", err)
	}

	// Sign with ed25519
	signature := ed25519.Sign(decodedKey, []byte(signatureBody))

	return base64.StdEncoding.EncodeToString(signature), nil
}

// GenerateSigningKeys generates a new ed25519 key pair for signing
func (c *ONDCCrypto) GenerateSigningKeys() (publicKey string, privateKey string, err error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", "", fmt.Errorf("error generating signing keys: %w", err)
	}

	publicKey = base64.StdEncoding.EncodeToString(pub)
	privateKey = base64.StdEncoding.EncodeToString(priv)

	return publicKey, privateKey, nil
}

// VerifyRequest verifies the signature of a request
func (c *ONDCCrypto) VerifyRequest(publicKeyStr string, payload []byte, created, expires int, signatureStr string) (bool, error) {
	// Compute BLAKE2b-512 hash over the payload
	hash := blake2b.Sum512(payload)
	digest := base64.StdEncoding.EncodeToString(hash[:])

	// Create computed message
	computedMessage := fmt.Sprintf(
		"(created): %d\n(expires): %d\ndigest: BLAKE-512=%s",
		created,
		expires,
		digest,
	)

	// Decode public key
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return false, fmt.Errorf("error decoding public key: %w", err)
	}

	// Decode signature
	receivedSignature, err := base64.StdEncoding.DecodeString(signatureStr)
	if err != nil {
		return false, fmt.Errorf("unable to base64 decode received signature: %w", err)
	}

	// Verify signature
	return ed25519.Verify(publicKeyBytes, []byte(computedMessage), receivedSignature), nil
}
