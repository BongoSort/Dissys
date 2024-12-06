package Utils

import (
	"bytes"
	"dissys/rsai"
	"encoding/gob"
	"fmt"
	"strconv"
)

type SignedTransaction struct {
	Tx        Transaction    // The message being sent along with the signed message
	Signature rsai.Signature // Potential signaturee coded as string
}

type Transaction struct {
	ID     string
	From   string
	To     string
	Amount int
}

/*
PARAMS:
ID - id of the transaction
From - The public key of the from participant
To - The public key of the to participant
Amount - The amount to be transferred
Signature - the signed transaction
*/
func CreateTransaction(ID string, From string, To string, Amount int) *Transaction {
	transaction := Transaction{
		ID:     ID,
		From:   From,
		To:     To,
		Amount: Amount,
	}
	return &transaction
}

func ConvertTransactionToByte(tx Transaction) []byte {
	var txBuffer bytes.Buffer

	encoder := gob.NewEncoder(&txBuffer)

	err := encoder.Encode(tx)

	if err != nil {
		fmt.Println("There was a problem converting the private key to a buffer")
	}
	return txBuffer.Bytes()
}

func ConvertBytesToTransaction(txBytes []byte) Transaction {
	// Struct for the key
	var transaction Transaction

	// Setup the buffer
	var txBuff bytes.Buffer
	txBuff.WriteString(string(txBytes))
	// Fill the buffer with the potential key
	fmt.Println("Trying to decode transaction.")
	err := gob.NewDecoder(&txBuff).Decode(&transaction)
	fmt.Println("Successfully decoded transaction.")

	// Read the potential key from the buffer into private key, with bigEndian (such that it works on all machines)
	if err != nil {
		fmt.Println("Error converting bytes to transaction!!")
	}
	return transaction
}

/*
Takes a transaction and signs it using a private key
*/
func CreateSignedTransaction(tx Transaction, key *rsai.PrivateKey) SignedTransaction {
	txBytes := ConvertTransactionToByte(tx)
	txSigned := rsai.SignStruct(txBytes, *key)
	sig := rsai.NewSignature(txSigned.Bytes(), txBytes)
	SignedTransaction := SignedTransaction{
		Tx:        tx,
		Signature: *sig,
	}
	return SignedTransaction
}

/*
Dont use this as it does not match signed transaction at the moment.
*/
func UnsignTransaction(tx Transaction, key rsai.PublicKey) Transaction {
	from := rsai.UnSignStruct([]byte(tx.From), key)
	to := rsai.UnSignStruct([]byte(tx.To), key)
	ID := rsai.UnSignStruct([]byte(tx.ID), key)
	stringAmount := strconv.Itoa(tx.Amount)
	amount := rsai.UnSignStruct([]byte(stringAmount), key)
	return *CreateTransaction(ID.String(), from.String(), to.String(), int(amount.Int64()))
}
