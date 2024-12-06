package Test

import (
	"bytes"
	"dissys/Server"
	"dissys/Utils"
	"dissys/rsai"
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"golang.org/x/exp/slices"
)

/*
Test a peer can make a TCP listen on a (random in implementation) port.
Do this by connecting to the port
*/
func TestListenerOnPort(t *testing.T) {
	/*
		1. Create connection on known port.
		2. Make Peer listen on same port, and see if it creates a listener anyways bu on a different
	*/
	peer := Server.NewPeer()
	/*
		want = net.Listener 0 - kald .getPortnumber to check
		got = net.Listener p number
	*/
	trueln := peer.ListenOnPort(0)
	got := trueln.Addr().Network()
	/*
		Addr.Network() returns if connection is TCP/UDP
	*/
	fakeln, _ := net.Listen("tcp", ":0")
	want := fakeln.Addr().Network()
	//fmt.Printf("%s", got)
	if got != want {
		t.Errorf("got %q wanted %q", got, want)
	}
}

/*
If a peer opens a connection on a port that is already taken.
It should find another port and open a connection there
*/
func TestAnotherPortIsFoundIfTaken(t *testing.T) {
	/*
		1. Determine a port that is already taken
		2. Ask a peer to open a port on this and see if it still manages to create a connection
	*/
	peer := Server.NewPeer()
	lnWant, _ := net.Listen("tcp", ":0")
	/*
		Get the port number for the connection just made
	*/
	portNumber := lnWant.Addr().(*net.TCPAddr).Port
	lnGot := peer.ListenOnPort(portNumber)
	/*
		Check the connection you got can be closed.
		No error closing == nil
		An error means the error is something which is not nil
	*/
	if lnGot.Close() != nil {
		t.Errorf("got %q wanted %q", lnGot.Addr().String(), lnWant.Addr().Network())
	}
}

/*
DUMB TEST
Test a peer can create a network and updates its field
*/
func TestDumb(t *testing.T) {
	hostPeer := Server.NewPeer()
	hostPeer.ListenOnPort(0)
	//testPeer.Connect(hostPeer.IpAddr, hostPeer.Port)
	if hostPeer.IpAddr == "" || hostPeer.Port == 0 {
		t.Errorf("Host ip %s or port %d has not changed", hostPeer.IpAddr, hostPeer.Port)
	}
}

/*
A peer can connect to another peer on the network
*/
func TestAPeerCreatesNetworkAndAPeerCanConnect(t *testing.T) {
	/*
		1. Take a peer and make it connect to nonexisting peer
		2. Test this by making another peer connect to it.
	*/
	hostPeer := Server.NewPeer()
	testPeer := Server.NewPeer()

	/*
		Host peer creates network
	*/
	go hostPeer.CreateNetwork()
	//TODO: fix sleep solution below
	// The intention of the sleep call is to make sure the connection is formed
	time.Sleep(2 * time.Second)
	testPeer.Connect(hostPeer.IpAddr, hostPeer.Port)

	/*
		1. test that the testPeer has a connection which is successfull
		2. testPeer has a successfull connection if it can close without getting err
	*/
	conn := testPeer.ConnectedTo[0]

	if conn.Close() != nil {
		t.Errorf("There was no connection available")
	}
}

func TestPeerCanHaveMultipleConnections(t *testing.T) {
	// One creates the network
	hostPeer := Server.NewPeer()
	testPeer1 := Server.NewPeer()
	testPeer2 := Server.NewPeer()
	testPeer3 := Server.NewPeer()
	hostPeer.CreateNetwork()
	ip := hostPeer.IpAddr
	port := hostPeer.Port

	/*
		Test both peers can connect
	*/
	testPeer1.Connect(ip, port)
	testPeer2.Connect(ip, port)
	testPeer3.Connect(ip, port)
	//TODO: fix sleep solution below
	// The intention of the sleep call is to make sure the connection is formed
	time.Sleep(2 * time.Second)
	/*
		Peers are connected if either the testPeer or host peer
		can close the connection (both are to be tested)
	*/
	totalConnections := len(hostPeer.ConnectedPeers)
	conn2 := hostPeer.ConnectedPeers[2]
	conn1 := testPeer1.ConnectedTo[0]
	if conn1.Close() != nil && conn2.Close() != nil && totalConnections != 3 {
		t.Errorf("Either the connected peer or host peer could not close the connection")
	}
}

func TestPeerJoiningCanHaveMltpConn(t *testing.T) {
	/*
		Host peer -> peer1 connecting to host
		-> peer2 connecting to peer1
		-> peer3 connecting to peer1
	*/
	hostPeer := Server.NewPeer()
	testPeer1 := Server.NewPeer()
	testPeer2 := Server.NewPeer()
	testPeer3 := Server.NewPeer()

	hostPeer.CreateNetwork()

	testPeer1.Connect(hostPeer.IpAddr, hostPeer.Port)

	ip := testPeer1.IpAddr
	port := testPeer1.Port

	go testPeer2.Connect(ip, port)
	go testPeer3.Connect(ip, port)

	time.Sleep(2 * time.Second)
	if len(testPeer2.ConnectedTo) != 1 && len(testPeer3.ConnectedTo) != 1 {
		t.Error("Somehow either p3 or p2 is not connected to p1")
	}
}

func TestPeerCanConnectToMultiple(t *testing.T) {
	/*
		1. Make a peer host a network (hostPeer)
		2. Make a peer join the hosted network (peer1)
		3. Make another peer connect to the above two peers (peer2)
	*/
	hostPeer := Server.NewPeer()
	peer1 := Server.NewPeer()
	peer2 := Server.NewPeer()

	hostPeer.CreateNetwork()
	peer1.Connect(hostPeer.IpAddr, hostPeer.Port)
	peer2.Connect(hostPeer.IpAddr, hostPeer.Port)
	peer2.Connect(peer1.IpAddr, peer1.Port)

	/*
		Peer2 now has 2 connections
		Check that peer2 have added both connections in list of current conn
		Check connection exists by closing them.
	*/
	//TODO: fix sleep solution below
	// The intention of the sleep call is to make sure the connection is formed
	time.Sleep(4 * time.Second)
	totalConnections := len(peer2.ConnectedTo)
	if totalConnections != 2 {
		t.Errorf("There was no connection available")
	}
}

var KeySize = 256

func TestMsgMarshalling(t *testing.T) {
	sk, pk := rsai.KeyGenStruct(KeySize)
	tx := Utils.CreateTransaction(rsai.MakePublicKeyString(pk), "peere2", "peer1", 100)
	sTd := Utils.CreateSignedTransaction(*tx, sk)
	msgBefore := Utils.Message{
		JoinRequest: true,
		Tnsa:        sTd,
		Closed:      nil,
	}

	jsonData := Utils.MustMarshalMsgToBytes(&msgBefore)
	msgAfter := Utils.MustDeMarshallBytesToMsg(&jsonData)

	//TODO: fix below method of checking
	if msgAfter.JoinRequest != msgBefore.JoinRequest {
		fmt.Println("msgBefore and msgAfter are not the same. Marshalling failed")
	}
}

func TestPeerCanSendMessageToAnotherPeer(t *testing.T) {
	//TODO: find alternative method of testing this quality
	t.Skip()
	sk, pk := rsai.KeyGenStruct(KeySize)
	tx := Utils.CreateTransaction(rsai.MakePublicKeyString(pk), "peere2", "peer1", 100)
	sTd := Utils.CreateSignedTransaction(*tx, sk)
	msgBefore := Utils.Message{
		JoinRequest: true,
		Tnsa:        sTd,
		Closed:      nil,
	}

	//Make peer1 and peer2
	peer1 := Server.NewPeer()
	peer2 := Server.NewPeer()

	//Create a network using peer1
	fmt.Println("Creating network")
	peer1.CreateNetwork()
	peer1Ip := peer1.IpAddr
	peer1Port := peer1.Port //Connect from peer2 to peer1

	//Connect from peer2 to peer1
	fmt.Printf("Connecting from peer2 to peer1's network on ip: %s with addr: %d\n", peer1Ip, peer1Port)
	peer2.Connect(peer1Ip, peer1Port)
	//TODO: Fix sleep solution
	time.Sleep(2 * time.Second)
	//Get connection that peer2 has in connectedTo
	fmt.Printf("Fetching connection from peer2 log, log has size %d\n", len(peer2.ConnectedTo))
	connSnd := peer2.ConnectedTo[0]    //This is the connection from peer1 to peer2
	connRcv := peer1.ConnectedPeers[0] //This is the connection from peer2 to peer1
	//Rvc a message on peer1 on the conn equivalent with peer2
	go peer1.RcvMsg(connRcv)
	//TODO: Fix sleep solution
	time.Sleep(2 * time.Second)

	//Message to be sent
	//msg := "Hello World!\n"
	go peer2.SendMsg(connSnd, msgBefore)

	//TODO: Fix sleep solution
	time.Sleep(2 * time.Second)

	//Check that msg that was rcved is the same as the sent msg.
	msgAfter := peer1.MsgLog[0]

	if msgBefore.JoinRequest != msgAfter.JoinRequest {
		t.Errorf("Error sending/recieving Message. MsgAfter was not the same as MsgBefore")
	}
}

func TestGetMsgFromLog(t *testing.T) {
	//Create a test Message
	sk, pk := rsai.KeyGenStruct(KeySize)
	tx := Utils.CreateTransaction(rsai.MakePublicKeyString(pk), "peere2", "peer1", 100)
	sTd := Utils.CreateSignedTransaction(*tx, sk)
	msgBefore := Utils.Message{
		JoinRequest: true,
		Tnsa:        sTd,
		Closed:      nil,
	}

	//Create a test peer
	testPeer := Server.NewPeer()

	//Add the testMsg to the MsgLog of testPeer
	testPeer.AddMsgToLog(msgBefore)

	_, err := testPeer.GetMsgFromLog(0)

	if err != nil {
		t.Errorf("Could not get Message from log: %v", err)
	}
}

func TestGetMsgFromLogOutOfBounds(t *testing.T) {
	//Create a test Message
	sk, pk := rsai.KeyGenStruct(KeySize)
	tx := Utils.CreateTransaction(rsai.MakePublicKeyString(pk), "peere2", "peer1", 100)
	sTd := Utils.CreateSignedTransaction(*tx, sk)
	msgBefore := Utils.Message{
		JoinRequest: true,
		Tnsa:        sTd,
		Closed:      nil,
	}

	//Create a test peer
	testPeer := Server.NewPeer()

	//Add the testMsg to the MsgLog of testPeer
	testPeer.AddMsgToLog(msgBefore)

	_, gotErr := testPeer.GetMsgFromLog(1)

	//Error should be not be nil. We are out of bounds. If err is nil, something is wrong
	if gotErr == nil {
		t.Errorf("Something went wrong, index should be out of bounds: %v", gotErr)
	}
}

func TestWhenAPeerJoinsANetWorkItHasItselfAsAPeer(t *testing.T) {
	//Make hostPeer
	//Make peer1 -> hostPeer
	//Check peer1.peers has the addr of itself with port
	hostPeer := Server.NewPeer()
	peer1 := Server.NewPeer()

	hostPeer.CreateNetwork()
	peer1.Connect(hostPeer.IpAddr, hostPeer.Port)

	//TMP Solution appending list and such
	//TODO find another method of solving this
	//Since it adds itself first to list of peers
	//It must hence be the last place in the list
	pList := peer1.PeersInNewtwork
	totalPeerAddr := peer1.IpAddr + ":" + strconv.Itoa(peer1.Port)
	if !slices.Contains(pList, totalPeerAddr) {
		t.Error("Has addr: ", totalPeerAddr)
	}
}

func TestAPeerCanSendPeerListAsMsg(t *testing.T) {
	//Make hostPeer
	//Make peer1 -> hostPeer
	//Make peer2 -> hostPeer
	//Check peer2 has peer1 and hostPeer in list of peers

	hostPeer := Server.NewPeer()
	peer1 := Server.NewPeer()
	peer2 := Server.NewPeer()

	hostPeer.CreateNetwork()
	peer1.Connect(hostPeer.IpAddr, hostPeer.Port)
	peer2.Connect(hostPeer.IpAddr, hostPeer.Port)

	peer1Addr := peer1.IpAddr + ":" + strconv.Itoa(peer1.Port)
	hostAddr := hostPeer.IpAddr + ":" + strconv.Itoa(hostPeer.Port)

	time.Sleep(2 * time.Second)
	time.Sleep(2 * time.Second)
	time.Sleep(2 * time.Second)
	containsPeer1 := slices.Contains(peer2.PeersInNewtwork, peer1Addr)
	containsHost := slices.Contains(peer2.PeersInNewtwork, hostAddr)

	if !(containsHost && containsPeer1) {
		t.Error("whaddafack")
	}
}

/*
Test FloodJoinMessage is sent when joining a network
*/
func TestFloodJoinIsSent(t *testing.T) {
	/*
		1. hostPeer
		2. p1 -> hostPeer
		3. p2 -> p1
		Check that hostPeer has p2 as a peer in the peer to peer network
	*/

	hostPeer := Server.NewPeer()
	p1 := Server.NewPeer()
	p2 := Server.NewPeer()

	hostPeer.CreateNetwork()
	p1.Connect(hostPeer.IpAddr, hostPeer.Port)
	p2.Connect(p1.IpAddr, p1.Port)

	time.Sleep(2 * time.Second)
	time.Sleep(2 * time.Second)
	if len(hostPeer.PeersInNewtwork) != 3 {
		t.Error("P2 did not send join msg to hostPeer")
	}
}

func TestHashIsTheSameBeforeAndAfterSignUnsignProcess(t *testing.T) {
	msg := "This is very funny"
	hashedMsg := rsai.HashMessage([]byte(msg))
	sk, pk := rsai.KeyGenStruct(KeySize)
	signedMsg := rsai.SignStruct([]byte(msg), *sk)
	unSignedMsg := rsai.UnSignStruct(signedMsg.Bytes(), *pk)
	byteUnsignedMsg := unSignedMsg.Bytes()
	if bytes.Equal(byteUnsignedMsg, hashedMsg) == false {
		t.Errorf("The Hashed Message is: %s, but unSigned Message was %s .\n", string(hashedMsg), string(byteUnsignedMsg))
	}
}

func TestSignedIsTheSameAsUnsignedAndHashedMsg(t *testing.T) {
	sk, pk1 := rsai.KeyGenStruct(KeySize)
	pkString := rsai.MakePublicKeyString(pk1)
	hpk := rsai.HashMessage([]byte(pkString))
	signedPk := rsai.SignStruct([]byte(pkString), *sk)
	unsignedPk := rsai.UnSignStruct(signedPk.Bytes(), *pk1)

	if bytes.Equal(unsignedPk.Bytes(), hpk) != true {
		t.Error("Sign and unsigned hash is not the same.")
	}
}

func TestDummyTest2(t *testing.T) {
	sk, pk1 := rsai.KeyGenStruct(KeySize * 2)
	pkString := rsai.MakePublicKeyString(pk1)
	tx := Utils.CreateTransaction("111", pkString, "", 100)
	//testTrans := Utils.CreateSignedTransaction(*tx, sk)

	//txSigned := Utils.ConvertBytesToTransaction(testTrans.Signature.SignedMsg)
	signedFrom := rsai.SignStruct([]byte(tx.From), *sk)
	sig := rsai.NewSignature(signedFrom.Bytes(), []byte(tx.From))

	if rsai.VerifyMsg(sig, *pk1) != true {
		t.Error("Failed to verify from")
	}
}

func TestConvertTxToAndFromByte(t *testing.T) {
	tx := Utils.CreateTransaction("111", "hunde", "", 100)
	byteTx := Utils.ConvertTransactionToByte(*tx)
	convertedTx := Utils.ConvertBytesToTransaction(byteTx)
	if tx.From != convertedTx.From {
		t.Errorf("%s is not the same as %s\n", tx.From, convertedTx.From)
	}
}

func TestTransactionVerificationPositiveTest(t *testing.T) {
	sk, pk := rsai.KeyGenStruct(KeySize)
	pkString := rsai.MakePublicKeyString(pk)
	tx := Utils.CreateTransaction("111", pkString, "delulu", 100)
	signedTx := Utils.CreateSignedTransaction(*tx, sk)

	sig := signedTx.Signature

	if rsai.VerifyMsg(&sig, *pk) != true {
		t.Error("it also seemed to easy.")
	}
}

func TestSignPublicKeyAndUnsignHashMatches(t *testing.T) {
	sk, pk := rsai.KeyGenStruct(KeySize)
	pkString := rsai.MakePublicKeyString(pk)

	pkSigned := rsai.SignStruct([]byte(pkString), *sk)

	sig := rsai.NewSignature(pkSigned.Bytes(), []byte(pkString))

	if rsai.VerifyMsg(sig, *pk) != true {
		t.Errorf("%s and %s are not a like", pkString, pkSigned.String())
	}
}

func TestYouCanGetPkFromSignedTransaction(t *testing.T) {
	sk, pk := rsai.KeyGenStruct(KeySize)
	pkString := rsai.MakePublicKeyString(pk)
	tx := Utils.CreateTransaction("111", pkString, "delulu", 100)
	txToBeSent := Utils.CreateSignedTransaction(*tx, sk)

	// Get public key from txToBeSent
	pkStringSent := Utils.ConvertBytesToTransaction(txToBeSent.Signature.Msg).From

	if pkStringSent != pkString {
		t.Errorf("%s is not the same as %s.\n", pkString, pkStringSent)
	}
}

func TestDummyTest3(t *testing.T) {
	sk, pk := rsai.KeyGenStruct(KeySize)
	pkString := rsai.MakePublicKeyString(pk)
	tx := Utils.CreateTransaction("111", pkString, "", 100)
	signedTrans := Utils.CreateSignedTransaction(*tx, sk)

	if Server.VerifyTransaction(&signedTrans) != true {
		t.Errorf("Unsuccessfull in verifying transaction.")
	}
}

func TestVerifyTransaction(t *testing.T) {
	sk, pk1 := rsai.KeyGenStruct(KeySize * 2)
	_, pk2 := rsai.KeyGenStruct(KeySize * 2)
	pkString := rsai.MakePublicKeyString(pk1)
	pk2String := rsai.MakePublicKeyString(pk2)
	tx := Utils.CreateTransaction("111", pkString, pk2String, 100)
	testTrans := Utils.CreateSignedTransaction(*tx, sk)

	signedTrans := Utils.CreateSignedTransaction(*tx, sk)

	if Server.VerifyTransaction(&signedTrans) != true {
		t.Errorf("The from field was not signed properly.")
	}
	if Server.VerifyTransaction(&testTrans) != true {
		t.Error("failed to verify transaction.")
	}

}

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

// TODO: Test a Peer.Transaction executes the transaction on the peer.
func TestPeerUpdatesLedgerWhenTransactionIsRecieved(t *testing.T) {
	/*
		1. HostPeer creates network
		2. p1 -> HostPeer
		2. p2 -> HostPeer
		3. Wait 1 sec for connection to be fully established
		4. FloodTransaction	from p1
		5. Check ledgers are the same for all peers
	*/

	sk1, pk1 := rsai.KeyGenStruct(KeySize)
	sk2, pk2 := rsai.KeyGenStruct(KeySize)
	for pk2 == pk1 {
		_, pk2 = rsai.KeyGenStruct(KeySize)
	}
	pk1String := rsai.MakePublicKeyString(pk1)
	pk2String := rsai.MakePublicKeyString(pk2)
	tx := Utils.CreateTransaction("111", pk1String, pk2String, 100)
	testTrans := Utils.CreateSignedTransaction(*tx, sk1)

	hostPeer := Server.NewPeer()
	p1 := Server.NewPeer()
	p2 := Server.NewPeer()

	hostPeer.CreateNetwork()
	p1.Connect(hostPeer.IpAddr, hostPeer.Port)
	p2.Connect(p1.IpAddr, p1.Port)
	time.Sleep(2 * time.Second)
	//testTrans := Utils.CreateTransaction("111", from, to, 100) //ID, From, To, Amount
	p1.FloodTransaction(&testTrans)
	time.Sleep(2 * time.Second)
	var hostLdg *Server.Ledger = &hostPeer.Ldg
	var p1Ldg *Server.Ledger = &p1.Ldg
	var p2Ldg *Server.Ledger = &p2.Ldg
	time.Sleep(2 * time.Second)

	//isFromSame := hostLdg.GetAccAmount(from) == p1Ldg.GetAccAmount(from) &&  p1Ldg.GetAccAmount(from) == p2Ldg.GetAccAmount(from)
	//isToSame := hostLdg.GetAccAmount(to) == p1Ldg.GetAccAmount(to) &&  p1Ldg.GetAccAmount(to) == p2Ldg.GetAccAmount(to)
	if p1Ldg.GetAccAmount(pk2String) != 100 || hostLdg.GetAccAmount(pk1String) != -100 {
		fmt.Printf("Do the private keys match: %v\n", sk1 == sk2)
		t.Error("Somehow somethin was wrong with the stuff. Check prints under and use to debug.")
	}
	fmt.Printf("Amount in p1 is %d amount in p2 is %d amount in host is %d\n", p1Ldg.GetAccAmount(pk1String), p2Ldg.GetAccAmount(pk2String), hostLdg.GetAccAmount(pk1String))
	fmt.Printf("Amount in p1 to is %d\n", p1.Ldg.GetAccAmount(pk2String))
}

//In case we need it ;)

//func TestPeerCanWriteAndRecieveOnConnection(t *testing.T) {
////Make peer1 and peer2
//peer1 := Server.NewPeer()
//peer2 := Server.NewPeer()

////Create a network using peer1
//fmt.Println("Creating network")
//peer1.CreateNetwork()
//peer1Ip := peer1.IpAddr
//peer1Port := peer1.Port //Connect from peer2 to peer1

////Connect from peer2 to peer1
//fmt.Printf("Connecting from peer2 to peer1's network on ip: %s with addr: %d\n", peer1Ip, peer1Port)
//peer2.Connect(peer1Ip, peer1Port)
////TODO: Fix sleep solution
//time.Sleep(2 * time.Second)
////Get connection that peer2 has in connectedTo
//fmt.Printf("Fetching connection from peer2 log, log has size %d\n", len(peer2.ConnectedTo))
//connSnd := peer2.ConnectedTo[0]    //This is the connection from peer1 to peer2
//connRcv := peer1.ConnectedPeers[0] //This is the connection from peer2 to peer1
////Rvc a message on peer1 on the conn equivalent with peer2
//go peer1.RcvMsg(connRcv)
////TODO: Fix sleep solution
//time.Sleep(2 * time.Second)

////Message to be sent
//msg := "Hello World!\n"
//go peer2.SendMsg(connSnd, msg)

////TODO: Fix sleep solution
//time.Sleep(2 * time.Second)

////Check that msg that was rcved is the same as the sent msg.
//rcvMsg := peer1.MsgLog[0]

//if rcvMsg != msg {
//t.Errorf("Error sending/recieving message: 'Hello World!'")
//}
//}
