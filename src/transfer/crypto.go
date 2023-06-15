package transfer

import (
	"crypto/rand"
	"errors"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
)

const PublicKeySize = chacha20poly1305.KeySize

type Crypto struct {
	publicKey  *[PublicKeySize]byte
	secretKey  *[PublicKeySize]byte
	sessionKey *[PublicKeySize]byte
}

func NewCrypto() *Crypto {
	crypto := &Crypto{
		publicKey:  new([PublicKeySize]byte),
		secretKey:  new([PublicKeySize]byte),
		sessionKey: new([PublicKeySize]byte),
	}

	_, _ = rand.Read(crypto.secretKey[:])
	curve25519.ScalarBaseMult(crypto.publicKey, crypto.secretKey)
	return crypto
}

func (c *Crypto) PublicKeySize() int {
	return PublicKeySize
}

func (c *Crypto) LocalPublicKey() []byte {
	return c.publicKey[:]
}

// SetRemotePublicKey 同时将生成sessionKey 已修正
func (c *Crypto) SetRemotePublicKey(remotePublicKey []byte) error {
	if len(remotePublicKey) != PublicKeySize {
		return errors.New("invalid remote public key length")
	}
	x25519, err := curve25519.X25519(c.secretKey[:], remotePublicKey)
	if err != nil {
		return err
	}
	copy(c.sessionKey[:], x25519)
	return nil
}

// SessionKeyDigest 生成验证码
func (c *Crypto) SessionKeyDigest() string {
	hash, _ := blake2b.New(16, nil)
	hash.Write(c.sessionKey[:])
	h := hash.Sum(nil)
	hashValue := uint64(0)
	for i := 0; i < 8; i++ {
		hashValue |= uint64(h[i]) << (uint(i) * 8)
	}
	return fmt.Sprintf("%06d", hashValue%1000000)
}

// Encrypt
// 加密数据并将nonce放到数据体前方
func (c *Crypto) Encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, chacha20poly1305.NonceSize)
	_, err := rand.Read(nonce[:])
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	block, err := chacha20poly1305.New(c.sessionKey[:])
	if err != nil {
		return nil, errors.New("sessionKey is invalid")
	}
	encrypted := block.Seal(nil, nonce, data, nil)
	bytes := append(nonce, encrypted...)
	return bytes, nil
}

// Decrypt  解密数据
func (c *Crypto) Decrypt(data []byte) ([]byte, error) {
	if len(data) < chacha20poly1305.NonceSize {
		return nil, errors.New("cipher text too short")
	}

	var nonce = data[:chacha20poly1305.NonceSize]
	var cipherText = data[chacha20poly1305.NonceSize:]
	block, err := chacha20poly1305.New(c.sessionKey[:])
	if err != nil {
		return nil, errors.New("sessionKey is invalid")
	}
	decrypted, err := block.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, errors.New("decryption failed")
	}

	return decrypted, nil
}
