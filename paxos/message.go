package paxos

type messageType int
const (
	PrepareRequest messageType = iota // Prepare Request (Proposer -> Acceptor)
	PrepareResponse // Prepare Response (Acceptor -> Proposer)
	AcceptRequest // Accept Request (Proposer -> Acceptor)
	AcceptResponse // Accept Response (Acceptor -> Learner)
)

type message struct {
	t  messageType // tipologia del messaggio
	id int  // id del mittente del messaggio
	n int  // proposal number: ovvero il numero di proposta
	v int // proposal value: ovvero il valore associato alla proposta
}


func NewPrepareRequestMessage(id, n int) message {
	return message{t: PrepareRequest, id: id, n: n} // Non porta con sé un valore v
}

func NewPrepareResponseMessage(id, n, v int) message {
	return message{t: PrepareResponse, id: id, n: n, v: v}
}

func NewAcceptRequestMessage(id, n, v int) message {
	return message{t: AcceptRequest, id: id, n: n, v: v}
}

// L'Accept response da Acceptor a Learner porta con se solo i valori accettati (no numero di proposta)
func NewAcceptResponseMessage(id, v int) message {
	return message{t: AcceptResponse, id: id, v: v}
}