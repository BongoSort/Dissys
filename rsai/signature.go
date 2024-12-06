package rsai

import (
	"bytes"
	"fmt"
	"math/big"
)

/*
A wrapper function for the sign method taking a key structure instead of d and n from the key directly.
*/
func SignStruct(msg []byte, privateKey PrivateKey) *big.Int {
	return Sign(msg, privateKey.D, privateKey.N)
}

func Sign(msg []byte, d *big.Int, n *big.Int) *big.Int {
	msgHash := HashMessage(msg)
	bigIntMsgHash := new(big.Int)
	bigIntMsgHash.SetBytes(msgHash)
	return Decrypt(bigIntMsgHash, d, n)
}

/*
A wrapper function for the unsign method taking a key structure instead of e and n from the key directly.
*/
func UnSignStruct(msg []byte, publicKey PublicKey) *big.Int {
	var i big.Int
	i.SetBytes(msg)
	return UnSign(&i, publicKey.E, publicKey.N)
}

func UnSign(msgHash *big.Int, e *big.Int, n *big.Int) *big.Int {
	return Encrypt(msgHash, e, n)
}

/*
VerifyHashMsg verifies that a incoming message is identical to the local message

Parameters:

	signedHashedMessage : The incoming message, already signed and hashed
	msg :	The local message
	e : The public exponent for encryption.
	n : The modulus used for encryption and decryption.
*/
func VerifyHashMsg(signedMsg []byte, msg []byte, e *big.Int, n *big.Int) bool {
	//First, convert []byte to a big.Int before we can use UnSign() on it
	hashValue := new(big.Int)
	hashValue.SetBytes(signedMsg)

	//Then, we use the UnSign() to remove the signature of the private key
	incMsg := UnSign(hashValue, e, n)

	//Then we need to compute the hash of the decrypted message, so we can compare it to the signed one we recieved
	localMsg := HashMessage(msg)
	return bytes.Equal(incMsg.Bytes(), localMsg)
}

func VerifyHashMsgStruct(signedHashedMessage []byte, msg []byte, key PublicKey) bool {
	return VerifyHashMsg(signedHashedMessage, msg, key.E, key.N)
}

func verifyMsgIntern(signedMsg []byte, msg []byte, key PublicKey) bool {
	//First, convert signedmsg []byte to a big.Int before we can use UnSign() on it
	signMsgBig := new(big.Int)
	signMsgBig.SetBytes(signedMsg)

	unsMsg := UnSignStruct(signedMsg, key)

	return bytes.Equal(unsMsg.Bytes(), msg)
}

/*
Takes the signed msg in the signature and unsigns it, then it hashes the msg to check if they match.
*/
func VerifyMsg(sig *Signature, key PublicKey) bool {
	unMsg := UnSignStruct(sig.SignedMsg, key)
	hMsg := HashMessage(sig.Msg)
	fmt.Println("Hashed bytes:")
	fmt.Printf("%v\n", hMsg)
	fmt.Println("Unsigned bytes:")
	fmt.Printf("%v\n", unMsg.Bytes())
	return bytes.Equal(hMsg, unMsg.Bytes())
}
