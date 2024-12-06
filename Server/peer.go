package Server

import (
	"dissys/Utils"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"sync"
)

type Peer struct {
	IpAddr          string
	Port            int
	ConnectedPeers  []net.Conn
	ConnectedTo     []net.Conn
	PeersInNewtwork []string
	MsgLog          []Utils.Message
	wg              sync.WaitGroup
	Ldg             Ledger
}

func NewPeer() *Peer {
	peer := Peer{
		IpAddr: "",
		Port:   0,
		Ldg:    *NewLedger(),
	}
	return &peer
}

/*
Taken from internet.
*/
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println("error: ", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func (p *Peer) ListenOnPort(port int) net.Listener {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		for err != nil {
			port = rand.Intn(65535 + 1)
			ln, err = net.Listen("tcp", ":"+strconv.Itoa(port))
		}
	}
	/*
		Update the field
	*/
	p.IpAddr = GetOutboundIP().String()
	p.Port = ln.Addr().(*net.TCPAddr).Port

	return ln
}
func (p *Peer) AckIncConn(ln net.Listener) {
	for {
		conn, _ := ln.Accept()
		// Add conn to the PeerList
		p.ConnectedPeers = append(p.ConnectedPeers, conn)
		go p.RcvMsg(conn)
		fmt.Println("Got a connection...")
	}
}

func (p *Peer) CreateNetwork() {
	randomPortNumber := rand.Intn(65535 + 1)
	ln := p.ListenOnPort(randomPortNumber)

	//When creating a peer network add itself as a peer in it
	tAddr := p.IpAddr + ":" + strconv.Itoa(p.Port)
	p.PeersInNewtwork = append(p.PeersInNewtwork, tAddr)

	// Go and manage peers trying to connect to port
	go p.AckIncConn(ln)
}

func (p *Peer) JoinNetWorkProcedure(conn net.Conn) {
	/*
		Step 1. Listen on a port so other can connect to this peer
		Step 2. Add itself to list of peers
		Step 3. Send request to recieve a list of peers
		Step 4. Send join request
	*/

	//Waitgroup is used to synchronize when peerlist is updated

	//Step 1
	if p.Port == 0 {
		ln := p.ListenOnPort(0)
		go p.AckIncConn(ln)
	}

	tAddr := p.IpAddr + ":" + strconv.Itoa(p.Port)
	//Step 1
	p.PeersInNewtwork = append(p.PeersInNewtwork, tAddr)

	//Step 3 send peer list request
	senderAddr := p.IpAddr + ":" + strconv.Itoa(p.Port)
	pRequest := Utils.NewPeerRequest(senderAddr)
	fmt.Printf("Peer request sent from %d to %s\n", p.Port, conn.RemoteAddr().String())
	p.SendMsg(conn, *pRequest)

	//Step 4 flood join request to all peers in network
	pJoinRg := Utils.NewJoinRequest(senderAddr)
	//p.SendMsg(conn, *pJoinRg)
	p.wg.Add(1)
	//TODO mby find different method of synchronizing?
	p.wg.Wait()
	p.FloodMessage(*pJoinRg)

	fmt.Println(strconv.Itoa(p.Port) + " Connected to Peer to peer...")
}

func (p *Peer) FloodMessage(msg Utils.Message) {
	pList := p.PeersInNewtwork
	fmt.Printf("FloodMSG called with pAmonut: %d\n", len(pList))
	ownAddr := p.IpAddr + ":" + strconv.Itoa(p.Port)
	for _, addr := range pList {
		// Check if its trying to send msg to send itself and stop it
		if addr == ownAddr {
			continue
		}
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Printf("Couldn't connect to: %s in floodmsg\n", addr)
		}
		defer conn.Close()
		fmt.Printf("Sending message to: %s by flood.\n", addr)
		p.SendMsg(conn, msg)
	}
}

func (p *Peer) Connect(ipAddr string, port int) {
	/*Concat iAddr and port*/
	totalAddr := ipAddr + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", totalAddr)

	if err != nil {
		/*Setup server*/
		p.CreateNetwork()
		return
	}

	//Joins network and updates connectedTo list
	p.ConnectedTo = append(p.ConnectedTo, conn)
	p.JoinNetWorkProcedure(conn)
}

func (p *Peer) SendMsg(conn net.Conn, msg Utils.Message) {
	btsMsg := Utils.MustMarshalMsgToBytes(&msg)
	_, err := conn.Write(btsMsg)
	if err != nil {
		fmt.Printf("There was a problem sending the msg with err: %s\n", err)
	}
}

func (p *Peer) RcvMsg(conn net.Conn) {
	// Create a buffer to read data from the connection
	fmt.Printf("Rcvmsg called from %d\n", p.Port)
	buffer := make([]byte, 4096*4) // Adjust the buffer size as needed
	for {
		// Read data from the connection into the buffer
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Someone closed the connection ", err)
			return
		}

		// Create a slice containing only the received data (remove any unused buffer space)
		receivedData := buffer[:bytesRead]

		// Unmarshal the received data into a Message struct
		originalMsg := Utils.MustDeMarshallBytesToMsg(&receivedData)
		//Add the message to the log --- or face errors ;)
		p.AddMsgToLog(originalMsg)

		// Where is the message comming from?
		go p.handleMessage(&originalMsg)
	}
}

func (p *Peer) handleMessage(msg *Utils.Message) {
	conn, err := net.Dial("tcp", msg.Sender)
	if err != nil {
		fmt.Printf("There was a problem connecting to the sender: %s, err: %s\n", msg.Sender, err)
	}
	defer conn.Close()
	sender := p.IpAddr + ":" + strconv.Itoa(p.Port)
	if msg.PeerRequest {
		//Send peer list to sender
		fmt.Printf("Sending peer list to sender. List size is %d\n", len(p.PeersInNewtwork))
		msg := Utils.NewSendPeerList(p.PeersInNewtwork, sender)
		p.SendMsg(conn, *msg)
		return
	}
	fmt.Println("Checking if it should update peer list")
	if len(msg.Peers) > 0 {
		fmt.Println("Recieved a peer list")
		peersSent := msg.Peers
		//Append all peers to the peer list
		p.PeersInNewtwork = append(p.PeersInNewtwork, peersSent...)
		p.wg.Done()
		return
	}
	fmt.Println("Checking if the message is a joinmessage")
	if msg.JoinRequest {
		fmt.Println("Updated peer list with new peer")
		p.PeersInNewtwork = append(p.PeersInNewtwork, msg.Sender)
		return
	}
	fmt.Println("Checking if the message is a transaction")
	if msg.IsTransaction {
		p.Ldg.Update(&msg.Tnsa)
		fmt.Println("Recieved a transaction")
		tx := Utils.ConvertBytesToTransaction(msg.Tnsa.Signature.Msg)
		newBalance := p.Ldg.GetAccAmount(tx.From)
		fmt.Printf("%s is now %d\n", tx.From, newBalance)
	}
}

func (p *Peer) AddMsgToLog(msg Utils.Message) {
	p.MsgLog = append(p.MsgLog, msg)
}

func (p *Peer) GetMsgFromLog(index int) (Utils.Message, error) {
	indexLimit := len(p.MsgLog)
	if 0 <= index && index < indexLimit {
		return p.MsgLog[index], nil
	}
	return Utils.Message{}, errors.New("INDEX OUT OF BOUNDS")
}

func (p *Peer) FloodTransaction(tx *Utils.SignedTransaction) {
	p.Ldg.Update(tx)
	senderAddr := p.IpAddr + ":" + strconv.Itoa(p.Port)
	msg := Utils.NewTransaction(*tx, senderAddr)
	p.FloodMessage(*msg)
}
