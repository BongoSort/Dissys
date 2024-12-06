package rsai

import (
	"crypto/aes"
	"crypto/cipher"
	"math/big"
	"os"
)

/*
Decrypt performs RSA decryption on the given cipherText using the provided
private exponent (d) and modulus (n).

Parameters:

	cipherText : The encrypted ciphertext to decrypt.
	d : The private exponent for decryption.
	n : The modulus used for encryption and decryption.

Returns:

	float64: The decrypted plaintext.
*/

func Decrypt(ciphertext *big.Int, d *big.Int, n *big.Int) *big.Int {
	plaintext := new(big.Int)
	plaintext.Exp(ciphertext, d, n)
	return plaintext
}

/*
Decrypts using AES count clock and gives it to a file
This code was inspired from a snippet on the internet

Parameters:

	key ([]byte) : The key to be used fro the AES cipher
	filename : The name of the file to decrypt
*/
func DecryptFromFile(key []byte, filename string) []byte {
	// Read the ciphertext from the file.
	cipherText, _ := os.ReadFile(filename)

	// Create a new AES cipher using key.
	ciph, _ := aes.NewCipher(key)
	// Create a new GCM stream.
	gcm, _ := cipher.NewGCM(ciph)

	// Decrypt the ciphertext.
	nonce := make([]byte, gcm.NonceSize())
	plaintext, _ := gcm.Open(nil, nonce, cipherText, nil)

	return plaintext
}
