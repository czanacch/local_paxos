package paxos

import "fmt"

type learner struct {
	id int // identificativo univoco del Learner
	receives chan message // canale per la ricezione dei messaggi (dagli Acceptor)
	nAcceptors int // numero totale degli acceptors
}

func NewLearner(id, nAcceptors int, receives chan message) *learner {
	l := new(learner)
	l.id = id
	l.nAcceptors = nAcceptors
	l.receives = receives
	return l
}

func (l *learner) Run() {
	number := 1
	for { // Fa questo lavoro continuamente
		guard := true
		responded := make(map[int]int) // mappatura da valori a n° di volte in cui è arrivato il valore

		for guard == true {
			msg := <-l.receives // estraggo dal mio singolo canale il messaggio msg
			switch msg.t {
			case AcceptResponse:
				responded[msg.v] = responded[msg.v] + 1 // Attenzione alla creazione mappa
				
				for v,num := range responded {
					if num >= l.nAcceptors {
						fmt.Println("Learner ",l.id,": Chosen ",v, "n° della scelta:",number)
						number++
						guard = false
						break
					}
				}
			}

		}
	}
}