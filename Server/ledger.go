package Server

import (
	"dissys/Utils"
	"dissys/rsai"
	"fmt"
	"sync"
)

type Ledger struct {
	Log []string //Might be usefull later??

	AccountMu sync.Mutex
	Accounts  map[string]int
}

func NewLedger() *Ledger {
	ldg := Ledger{
		//Empty accounts from the beginning
		Accounts: make(map[string]int),
	}
	return &ldg
}

func (ldg *Ledger) AddAccount(accName string) {
	ldg.Accounts[accName] = 0
}

func (ldg *Ledger) AddAmountToAccountAndCreateAcc(accName string, amount int) {
	ldg.Accounts[accName] = amount
}

func (ldg *Ledger) GetAccAmount(accName string) int {
	return ldg.Accounts[accName]
}

func VerifyTransaction(tx *Utils.SignedTransaction) bool {
	/*
		1. Verify the t.signature is a valid rsa signature on the fields
		- under the public key t.From.
		- the public key is the t.from
	*/
	transaction := tx.Tx
	n, e := rsai.GetPublicKeyParts(transaction.From)
	publicKey := rsai.NewPublicKey(n, e)

	return rsai.VerifyMsg(&tx.Signature, *publicKey)
}

func (ldg *Ledger) Update(tx *Utils.SignedTransaction) {
	//Remove money from "from"
	ldg.AccountMu.Lock()
	defer ldg.AccountMu.Unlock()
	if VerifyTransaction(tx) {
		verifiedTransaction := Utils.ConvertBytesToTransaction(tx.Signature.Msg)
		from := verifiedTransaction.From
		to := verifiedTransaction.To
		id := verifiedTransaction.ID
		amount := verifiedTransaction.Amount

		ldg.Accounts[from] -= amount
		ldg.Accounts[to] += amount
		ldg.Log = append(ldg.Log, id)
	} else {
		ldg.Log = append(ldg.Log, "Failed to verify sender.")
		fmt.Println("Failed to verify transaction.")
	}
}
