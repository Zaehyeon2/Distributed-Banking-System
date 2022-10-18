package fsm

type CommandPayload struct {
	Operation string
	Account   string
	Sender    string
	Receiver  string
	Amount    uint64
}
