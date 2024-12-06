package main

import (
	"dissys/Server"
	"dissys/Utils"
	"dissys/rsai"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

/*
The purpose of this struct is to have a collection of peers in the network
This is used such that a random peer to peer network can be generated each time it is run
*/
type pair struct {
	IpAddr string
	Port   int
}

func NewPair() *pair {
	pr := pair{
		IpAddr: "",
		Port:   0,
	}
	return &pr
}

func main() {
	peerCreatingNetWork := Server.NewPeer()
	peerCreatingNetWork.Connect("doesn't matter", 6969)
	prCN := NewPair()
	addPeerToPair(prCN, peerCreatingNetWork)
	// Spawn peers connecting to other random peers
	n := 15
	var count int = 0
	var pList []*Server.Peer
	var posNames []string

	// Generate accounts and sk keys to create transactions
	sk1, pk1 := rsai.KeyGenStruct(256 * 4)
	sk2, pk2 := rsai.KeyGenStruct(256 * 4)
	sk3, pk3 := rsai.KeyGenStruct(256 * 4)
	sk4, pk4 := rsai.KeyGenStruct(256 * 4)
	sk5, pk5 := rsai.KeyGenStruct(256 * 4)
	pk1String := rsai.MakePublicKeyString(pk1)
	pk2String := rsai.MakePublicKeyString(pk2)
	pk3String := rsai.MakePublicKeyString(pk3)
	pk4String := rsai.MakePublicKeyString(pk4)
	pk5String := rsai.MakePublicKeyString(pk5)
	posNames = append(posNames, pk1String, pk2String, pk3String, pk4String, pk5String)
	pList = append(pList, peerCreatingNetWork)

	var skList []rsai.PrivateKey
	skList = append(skList, *sk1, *sk2, *sk3, *sk4, *sk5)

	for count < n {
		p := Server.NewPeer()
		randomConn := rand.Intn(len(pList))
		p.Connect(pList[randomConn].IpAddr, pList[randomConn].Port)
		pList = append(pList, p)
		count++
	}
	//Making sure peer 2 peer network is fully connected
	//Otherwise because of lazy flooding ledgers arent really alike
	time.Sleep(10 * time.Second)
	count = 0
	msgAmount := 10
	for count < msgAmount {
		tx := generateRandomTx(posNames, skList)
		go pList[rand.Intn(len(pList))].FloodTransaction(&tx)
		count++
	}
	time.Sleep(2 * time.Second)
	count = 0
	for count < len(pList) {
		innerCounter := 0
		for innerCounter < len(posNames) {
			printLedger(posNames[innerCounter], pList[count], count)
			innerCounter++
		}
		count++
	}
}

func generateRandomTx(pkList []string, skList []rsai.PrivateKey) Utils.SignedTransaction {
	amount := rand.Intn(51)
	pkAndSkIndex := rand.Intn(len(skList))
	randIndex := rand.Intn(len(skList))
	from := pkList[pkAndSkIndex]
	to := pkList[randIndex]
	// Doing it this way, so not all transactions are verified.
	sk := skList[pkAndSkIndex]
	tx := Utils.CreateTransaction(MakeRandomId(), from, to, amount)
	signedTx := Utils.CreateSignedTransaction(*tx, &sk)
	return signedTx
}

func MakeRandomId() string {
	// Specify the length of the random string
	length := 12

	// Generate a random string
	randomString := generateRandomString(length)

	return randomString
}

func generateRandomString(length int) string {
	// Calculate the number of bytes needed for the given string length
	numBytes := length / 2 // Hex encoding is two chars pr byte

	// Create a byte slice to hold the random bytes
	randomBytes := make([]byte, numBytes)

	// Convert the random bytes to a hexadecimal string
	randomString := hex.EncodeToString(randomBytes)

	return randomString
}

func printLedger(acc string, p *Server.Peer, pNumber int) {
	accEasy := acc[:4]
	fmt.Printf("PID %d has account: %s with value: %d\n", pNumber, accEasy, p.Ldg.GetAccAmount(acc))
}

func addPeerToPair(pr *pair, peer *Server.Peer) {
	pr.IpAddr = peer.IpAddr
	pr.Port = peer.Port
}
