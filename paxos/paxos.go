package paxos

type network struct {
	proposers []*proposer
	acceptors []*acceptor
	learners  []*learner
}

// funzione che prende un numero intero e restituisce un array di taglia n di canali
func makeChannels(n int) []chan message {
	chans := make([]chan message, n)
	for i := range chans {
		chans[i] = make(chan message, 1024) // ogni canale dell'array è grande 1024
	}
	return chans
}

func NewNetwork(nProposers, nAcceptors, nLearners int, vs []int) *network {
	cProposers := makeChannels(nProposers) // crea un array di canali: uno per ogni Proposer
	cAcceptors := makeChannels(nAcceptors) // crea un array di canali: uno per ogni Acceptor
	cLearners := makeChannels(nLearners) // crea un array di canali: uno per ogni Learner

	n := new(network) // n è un puntatore ad un oggetto network

	n.proposers = make([]*proposer, nProposers) // inizializza la slice con puntatori ai Proposer	
	for i := range n.proposers {
		n.proposers[i] = NewProposer(i, vs[i], nProposers, cProposers[i], cAcceptors)
	}

	n.acceptors = make([]*acceptor, nAcceptors) // inizializza la slice con puntatori agli Acceptor
	for i := range n.acceptors {
		n.acceptors[i] = NewAcceptor(i, cAcceptors[i], cProposers, cLearners)
	}

	n.learners = make([]*learner, nLearners) // inizializza la slice con puntatori ai Learner
	for i := range n.learners {
		n.learners[i] = NewLearner(i, nAcceptors, cLearners[i])
	}

	return n
}

// data una slice di canali e un messaggio, trasmette il messaggio su ogni canale
func broadcast(channels []chan message, msg message) {
	for _, channel := range channels {
		channel <- msg
	}
}

// Questa funzione fa partire l'algoritmo
// fa partire una goroutine per ogni processo
func (n *network) Start() {

	for _, l := range n.learners {
		go l.Run()
	}

	for _, a := range n.acceptors {
		go a.Run()
	}

	for _, p := range n.proposers {
		go p.Run()
	}
	
}