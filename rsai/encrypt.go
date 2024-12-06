package rsai

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"math/big"
	"os"
)

/*
Encrypt performs RSA encryption on the given plaintext using the provided
public exponent (e) and modulus (n).

Parameters:

	plainText : The plaintext to encrypt.
	e : The public exponent for encryption.
	n : The modulus used for encryption and decryption.

Returns:

	float64: The encrypted ciphertext.
*/
func Encrypt(plainText *big.Int, e *big.Int, n *big.Int) *big.Int {
	ciphertext := new(big.Int)
	ciphertext.Exp(plainText, e, n)
	return ciphertext
}

/*
EncryptToFile encrypts a file using aes in counter mode
This code was inspired from a snippet on the internet

Parameters

	key : The public key used to Encrypt. Must be 16, 24, or 32 byte
	filename : The name of the returned, encoded file
	plaintext : The data that is to be encoded
*/
func EncryptToFile(key []byte, filename string, plaintext []byte) {
	ciph, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(ciph)

	nonce := make([]byte, gcm.NonceSize())
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Write the ciphertext to the file.
	os.WriteFile(filename, ciphertext, 0644)
}

/*
HashMessage takes a message in bytes and returns the sha256 on the msg

Parameters

	msg : The message in bytes to be hashed with sha256

Returns

	The sha256 hash of the msg
*/
func HashMessage(msg []byte) []byte {
	h := sha256.New()
	h.Write(msg)
	return h.Sum(nil)
}
