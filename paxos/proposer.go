package paxos

import (
	"fmt"
)

type proposer struct {
	id int // Identificativo univoco del Proposer
	v int // valore che si desidera proporre
	n int // proposal number (detto ballot)
	nProposers int // numero totale dei Proposer nella rete
	
	receives  chan message // unico canale del Proposer
	acceptors []chan message // array di canali: uno per ogni Acceptor
}

// Funzione per creare un nuovo Proposal
func NewProposer(id, v, nProposers int, receives chan message, acceptors []chan message) *proposer {
	p := new(proposer)
	p.id = id
	p.v = v
	p.n = id // ATTENZIONE: il proposal number viene messo al suo id automaticamente (numeri disgiunti)
	p.nProposers = nProposers
	p.receives = receives
	p.acceptors = acceptors

	return p
}

func (p *proposer) prepareRequest() {
	p.n = p.n + p.nProposers // TODO: è necessario questo incremento subito qui?
	msg := NewPrepareRequestMessage(p.id, p.n)
	fmt.Println("Proposer",p.id,": sending Prepare Request",p.n)
	broadcast(p.acceptors, msg)
}

func (p *proposer) acceptRequest() {
	msg := NewAcceptRequestMessage(p.id, p.n, p.v) // crea un messaggio Accept Request
	fmt.Println("Proposer",p.id,": sending Accept Request")
	broadcast(p.acceptors, msg)
}

// Metodo che fa partire l'algoritmo di un Proposer
func (p *proposer) Run() {

		max_number_received := 0

		fmt.Println("Proposer ",p.id,": started")

			// PHASE 1: Prepare Request & Promessa
			p.prepareRequest()
			
			responded := make(map[int]bool) // mappa (dizionario) da interi a booleani
			
			for len(responded) < len(p.acceptors)/2+1 { // finchè non mi sono arrivate più della metà risposte degli Acceptor
				msg := <-p.receives // estraggo dal mio singolo canale il messaggio msg
				switch msg.t {
				case PrepareResponse:
					responded[msg.id] = true // segnala che l'acceptor ha risposto promettendo
					if msg.n != 0 && max_number_received < msg.n {
						max_number_received = msg.n
						p.v = msg.v // adeguati e sintonizzati con il nuovo valore (il più alto mandato)
					}
				default:
				}
			}
			
			// PHASE 2: Accept Request
			p.acceptRequest()

}