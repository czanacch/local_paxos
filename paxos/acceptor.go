package paxos

import "fmt"

type acceptor struct {
	id int // identificativo univoco dell'Acceptor
	n int // ultimo proposal number relativo ad una ACCETTAZIONE
	v int // ultimo valore in custodia (accettato)
	max_n int // massimo proposal number che è stato promesso
	receives  chan message // canale attraverso il quale si rivevono messaggi
	proposers []chan message // array di canali attraverso i quali si inviano messaggi (ai Proposer)
	learners  []chan message // array di canali: uno per ogni Learner
}

func NewAcceptor(id int, receives chan message,	proposers, learners []chan message) *acceptor {
	a := new(acceptor)
	a.id = id
	a.n = 0 // ATTENZIONE: il numero di proposta accettata viene messo a 0 automaticamente
	a.v = 0 // ATTENZIONE: il valore viene messo a 0 automaticamente
	a.max_n = 0 // ATTENZIONE: il max_n viene messo a 0 automaticamente
	a.receives = receives
	a.proposers = proposers
	a.learners = learners
	return a
}

func (a *acceptor) Run() {
	fmt.Println("Acceptor ",a.id,": started")
	
	for {
		msg := <-a.receives
		switch msg.t {
		
		case PrepareRequest: // Nel caso in cui il messaggio sia una Prepare request
			if msg.n > a.max_n { // Se il numero è maggiore di qualsiasi cosa promessa
				a.max_n = msg.n // Nuovo massimo numero di proposta
				a.proposers[msg.id] <- NewPrepareResponseMessage(a.id, a.n, a.v)
				fmt.Println("Acceptor ",a.id,": sending Prepare Request")
			}
			
		case AcceptRequest: // Nel caso in cui il messaggio sia una Accept request
			if msg.n >= a.max_n { // se il proposal number è maggiore o uguale a ciò che ho promesso
				a.max_n = msg.n // aggiorna il max proposal number promesso
				a.n = msg.n // aggiorna il proposal number (relativo a un accettazione)
				a.v = msg.v // aggiorna il valore, finalmente

				broadcast(a.learners, NewAcceptResponseMessage(a.id, a.v))

				fmt.Println("Acceptor",a.id,": sending Accept Response")
			}
		default:
		}
	}

}