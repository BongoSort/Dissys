# TEST LIST

# PEER
- Assuming that connecting to a Peer network failed and the CreateNetwork method is called 

1. **OK** A peer can listen on a random port
1. **OK** A peer can accept a connection on port listened to 
1. **OK** A peer will find another port, if it tries to listen on a port and gets an error
1. **OK** Can create a Transaction "object"
1. **OK** A peer can connect to another peer.
1. **OK** A peer can send a Message from one peer to another peer
1. **OK** A peer can send a Transaction on the connection 
1. **OK** A peer can receive a Transaction on the connection 
1. **OK** A peer will call an update to its Ledger upon receiving a transaction
1. **OK** A peer can send a Flood using FloodTransaction, notifying all other peers on the network
1. **OK** A peer can send a Join Request
1. **OK** A Message can be marshalled 
1. **OK** A peer hosting a peernetwork can have multiple peers connected to it.
1. **OK** A peer joining a peer network can have multiple peers connected to it.
1. **OK** A peer can have multiple connections to other peers

//When connecting to an existing peer network If connection fails an exception is thrown If connection fails then CreateNetwork is called If connection is successful, the Peer should receive an update to its ledger

# LEDGER
Ledger stuff goes here


# Important shit
Can send join request
Can send a transaction
Can send a "someoneclosed" request
Can send between connections
Setup flood msg


type Message struct {
	JoinRequest bool
	Tns         Server.Transaction
	Closed      net.Addr
}