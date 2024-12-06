# Readme

## How to run the handin.go file

1. Use the terminal
2. cd dissys/handin
3. go run handin.go
4. See the result in the terminal

We have only taken the first 4 chars from the public key as it is unreadable otherwise.
We kept the tests and updated them from the previous exercise to make sure everything had the same properties.
We had one where we utilized floodtransaction.
This still does floodtransaction however it does so with a verified transaction.
We also have a single test testing that when a peer recieves a transaction from an advesary that it is not verified and the transaction is not performed.
That transaction can be viewed below.
```go
func TestTransactionIsRejected(t *testing.T) {
	sk, pk := rsai.KeyGenStruct(KeySize)
	skFake, pkAdvesary := rsai.KeyGenStruct(KeySize)
	// Make sure the private key are different
	for skFake == sk {
		skFake, pkAdvesary = rsai.KeyGenStruct(KeySize)
	}
	pkString := rsai.MakePublicKeyString(pk)
	pkAString := rsai.MakePublicKeyString(pkAdvesary)
	tx := Utils.CreateTransaction("1", pkString, pkAString, 100)
	signedTx := Utils.CreateSignedTransaction(*tx, skFake)
	p1 := Server.NewPeer()
	p1.Ldg.Update(&signedTx)

	if p1.Ldg.GetAccAmount(pkString) != 0 {
		t.Error("Transaction should have been rejected and no money should have been taken.")
	}
}
```