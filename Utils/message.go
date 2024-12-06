package Utils

import "net"

type Message struct {
	JoinRequest   bool
	IsTransaction bool
	Tnsa          SignedTransaction
	Closed        net.Addr
	Peers         []string
	PeerRequest   bool
	Sender        string
}

func NewTransaction(tx SignedTransaction, sender string) *Message {
	msg := Message{
		Sender:        sender,
		Tnsa:          tx,
		IsTransaction: true,
	}
	return &msg
}

func NewJoinRequest(sender string) *Message {
	msg := Message{
		Sender:      sender,
		JoinRequest: true,
	}
	return &msg
}

func NewMessage(jr bool, tnsa SignedTransaction, cls net.Addr) *Message {
	msg := Message{
		JoinRequest: jr,
		Tnsa:        tnsa,
		Closed:      cls,
	}
	return &msg
}

func NewPeerRequest(sender string) *Message {
	msg := Message{
		PeerRequest: true,
		Sender:      sender,
	}
	return &msg
}

func NewSendPeerList(peerList []string, sender string) *Message {
	msg := Message{
		Peers:  peerList,
		Sender: sender,
	}
	return &msg
}
