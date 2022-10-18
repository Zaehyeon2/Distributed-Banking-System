package fsm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/hashicorp/raft"
)

type bankingFSM struct {
	mtx     sync.RWMutex
	ledgers map[string]uint64
}

// deposit deposit a certain amount of tokens into one's accont
func (f *bankingFSM) deposit(account string, amount uint64) error {

	// Lock FSM during depositing
	f.mtx.Lock()
	// Unlock FSM after depositing
	defer f.mtx.Unlock()

	// Check if the account exists in the ledger.
	_, exist := f.ledgers[account]
	if exist {
		// Add a certatin amount of tokens into account.
		f.ledgers[account] += amount
	} else {
		// Create new account and add a certatin amount of tokens into account.
		f.ledgers[account] = amount
	}

	return nil
}

// transfer transfer a certain amount of tokens from one's account into another
func (f *bankingFSM) transfer(sender string, receiver string, amount uint64) error {
	// Lock FSM during transferring
	f.mtx.Lock()
	// Unlock FSM after transferring
	defer f.mtx.Unlock()

	// Check if the sender's account exists in the ledger.
	_, exist := f.ledgers[sender]
	if exist {
		if f.ledgers[sender] < amount {
			// Check if the sender has enough amount of tokens.
			return fmt.Errorf("The sender '%s' does not have enough tokens.", sender)
		}
		// Subtract a certatin amount of tokens from sender's account.
		f.ledgers[sender] -= amount

		// Check if the receiver's account exists in the ledger.
		_, exist = f.ledgers[receiver]
		if exist {
			// Add a certatin amount of tokens into receiver's account.
			f.ledgers[receiver] += amount
		} else {
			// Create new account and add a certatin amount of tokens into receiver's account.
			f.ledgers[receiver] = amount
		}

		return nil
	} else {
		// Check if the sender's account exists.
		return fmt.Errorf("The sender's account '%s' does not exist.", sender)
	}
}

// get query a balance of account
func (f *bankingFSM) get(account string) (uint64, error) {

	// Check if the account exists in the ledger.
	balance, exist := f.ledgers[account]
	if exist {
		return balance, nil
	} else {
		return 0, errors.New(fmt.Sprintf("The account '%s' does not exist.", account))
	}
}

// Apply apply a command
func (f *bankingFSM) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		var payload = CommandPayload{}
		if err := json.Unmarshal(log.Data, &payload); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error marshalling store payload %s\n", err.Error())
			return nil
		}
		// Operation including deposit, transfer, and get
		op := strings.ToUpper(strings.TrimSpace(payload.Operation))
		switch op {
		case "DEPOSIT":
			return &ApplyResponse{
				Error: f.deposit(payload.Receiver, payload.Amount),
				Data:  nil,
			}
		case "TRANSFER":
			return &ApplyResponse{
				Error: f.transfer(payload.Sender, payload.Receiver, payload.Amount),
				Data:  nil,
			}
		case "GET":
			data, err := f.get(payload.Account)
			return &ApplyResponse{
				Error: err,
				Data:  data,
			}
		}
	}

	_, _ = fmt.Fprintf(os.Stderr, "not raft log command type\n")
	return nil
}

// Snapshot manually force Raft to take a snapshot
func (f *bankingFSM) Snapshot() (raft.FSMSnapshot, error) {
	return newSnapshotNoop()
}

// Restore manually force Raft to consume an external snapshot
func (f *bankingFSM) Restore(rClose io.ReadCloser) error {
	defer func() {
		if err := rClose.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[FINALLY RESTORE] close error %s\n", err.Error())
		}
	}()

	_, _ = fmt.Fprintf(os.Stdout, "[START RESTORE] read all message from snapshot\n")
	var totalRestored int

	decoder := json.NewDecoder(rClose)
	for decoder.More() {
		var data = &CommandPayload{}
		err := decoder.Decode(data)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] error decode data %s\n", err.Error())
			return err
		}

		if err := f.deposit(data.Account, data.Amount); err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] error persist data %s\n", err.Error())
			return err
		}

		totalRestored++
	}

	// read closing bracket
	_, err := decoder.Token()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] error %s\n", err.Error())
		return err
	}

	_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] success restore %d messages in snapshot\n", totalRestored)
	return nil
}

// NewBank create a new bankFSM
func NewBank() raft.FSM {
	return &bankingFSM{ledgers: make(map[string]uint64)}
}
