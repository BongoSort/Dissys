/*
Exercise 9.11 (Software Wallet)

Use your solution from Exercise 6.10 to create a software
wallet for an RSA secret key. It should have these functions.

1. Generate(filename string, password string) string which generates a public key
and secret key and places the secret key on disk in the file with filename in filename
encrypted under the password in password. The function returns the public key.

2. Sign(filename string, password string, msg []byte) Signature which if the
filename and password match that of Generate will sign msg and return the signature.
You pick what the type Signature is.

Your solution should:

1. Make measures that make it costly for an adversary which gets its hands on the keyfile to
bruteforce the password.
2. Describe clearly what measure have been taken.
3. Explain why the system was designed the way it was and, in particular, argue why the
system achieves the desired security properties.
4. Test the system and describe how it was tested.
5. Describe how your TA can run the system and how to run the test

*/

package rsai

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/big"
	"strings"
)

type Signature struct {
	SignedMsg []byte
	Msg       []byte
}

func NewSignature(signedMsg []byte, msg []byte) *Signature {
	s := Signature{
		SignedMsg: signedMsg,
		Msg:       msg,
	}
	return &s
}

/*
Generates a public key and secret key and places the secret key on disk in the
file with "filename" encrypted under the password.
The function returns the public key.

Parameters:

	filename: name of the file the secret key is placed in
	password: password required to sign, has to be at least 12 chars long, and not all the same char

Returns:

	publicKey: The public key
*/
func Generate(filename string, password string) string {
	/* Reject password if it is not 12 chars long or larger than that */
	if len(password) < 12 {
		panic("Password must be more that 12 characters long")
	}

	/* For improved password security, require password to be complex, so no 000000000000  or aaaaaaaaaaaa */
	if isTooSimplePassword(password) {
		panic("Password is too simple, please use a more complex password.")
	}

	bytePassword := []byte(password)
	hashPass := HashMessage(bytePassword)

	privateKey, publicKey := KeyGenStruct(len(hashPass))

	var privateKeyBuffer bytes.Buffer

	encoder := gob.NewEncoder(&privateKeyBuffer)

	err := encoder.Encode(privateKey)

	if err != nil {
		fmt.Println("There was a problem converting the private key to a buffer")
	}

	publicKeyString := MakePublicKeyString(publicKey)

	EncryptToFile(hashPass, filename, privateKeyBuffer.Bytes()) // Encrypts the file with the hashed password.

	return publicKeyString
}

/*
Sign msg and return the signature if the filename match a file, and the password is the correct one for that file.

Parameters:

	filename: Name of the file to be signed
	password: The required password to be able to sign the message
	msg: The message we want to sign with
*/
func SignWithPass(filename string, password string, msg []byte) *Signature {

	//1. Hash input password
	passwdHash := HashMessage([]byte(password))

	//2. Decrypt filenae
	privateKeyBytes := DecryptFromFile(passwdHash, filename)

	// Struct for the key
	var privateKey PrivateKey

	// Setup the buffer
	var pkBuff bytes.Buffer
	pkBuff.WriteString(string(privateKeyBytes))
	// Fill the buffer with the potential key
	fmt.Println("Trying to decode private key.")
	err := gob.NewDecoder(&pkBuff).Decode(&privateKey)
	fmt.Println("Successfully decoded privatekey-")

	// Read the potential key from the buffer into private key, with bigEndian (such that it works on all machines)
	if err != nil {
		fmt.Println("Wrong password used!")
	}

	//3. Sign msg with private key
	n := privateKey.N
	d := privateKey.D

	signedMsg := Sign(msg, d, n)
	signedMsgBytes := signedMsg.Bytes()

	return NewSignature(signedMsgBytes, msg)

}

/*
Makes a string corresponding to a PrivateKey (n,d)

Parameters:

	n: The modulus, a part of both the public and private keys.
	d: The private exponent used for decryption

Returns:

	The string corresponding to a publickey or privatekey
*/
func MakePrivateKeyString(key *PrivateKey) string {

	nStr := key.N.String()
	dStr := key.D.String()

	return nStr + "-" + dStr
}

/*
Makes a string corresponding to a PublicKey (n,e)

Parameters:

	n: The modulus, a part of both the public and private keys.
	e: The public exponent used for encryption.

Returns:

	The string corresponding to a PublicKey
*/
func MakePublicKeyString(key *PublicKey) string {

	nStr := key.N.String()
	eStr := key.E.String()

	return nStr + "-" + eStr
}

/*
Makes a PublicKey (n,e) from a string

Parameters:

	publicKey: string

Returns:

	(n, e)
	n: The modulus,
	e: The public exponent used for encryption.
*/
func GetPublicKeyParts(publicKey string) (*big.Int, *big.Int) {
	keyParts := strings.Split(publicKey, "-")

	var n big.Int
	N, successN := n.SetString(keyParts[0], 10)

	var e big.Int
	E, successE := e.SetString(keyParts[1], 10)

	if !successE || !successN {
		fmt.Println("Error converting from string to bigInt")
	}
	return N, E
}

/*
Checks if the password is too simple (i.e. all characters are the same)

Parameters:

	password: The password to be checked

Returns:

	bool, dependent on if the password is too simple or not
*/
func isTooSimplePassword(password string) bool {
	if len(password) != 12 {
		return false
	}

	for i := 1; i < len(password); i++ {
		if password[i] != password[0] {
			return false
		}
	}
	return true
}
