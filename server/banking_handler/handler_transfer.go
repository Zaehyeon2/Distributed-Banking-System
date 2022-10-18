package banking_handler

import (
	"dbs/fsm"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/raft"
	"github.com/labstack/echo"
)

type requestTransfer struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   uint64 `json:"amount"`
}

// Deposit handling transfer a certain amount of tokens from one's account into one's account
func (h handler) Transfer(eCtx echo.Context) error {
	var form = requestTransfer{}
	if err := eCtx.Bind(&form); err != nil {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error binding: %s", err.Error()),
		})
	}

	if form.Sender == "" {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": "The sender's account is empty",
		})
	}

	if form.Receiver == "" {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": "The receiver's account is empty",
		})
	}

	if h.raft.State() != raft.Leader {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": "not the leader",
		})
	}

	payload := fsm.CommandPayload{
		Operation: "TRANSFER",
		Sender:    form.Sender,
		Receiver:  form.Receiver,
		Amount:    form.Amount,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error preparing saving data payload: %s", err.Error()),
		})
	}

	applyFuture := h.raft.Apply(data, 500*time.Millisecond)
	if err := applyFuture.Error(); err != nil {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error persisting data in raft cluster: %s", err.Error()),
		})
	}

	resp, ok := applyFuture.Response().(*fsm.ApplyResponse)
	if !ok {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error response is not match apply response"),
		})
	}
	if resp.Error == nil {
		return eCtx.JSON(http.StatusOK, map[string]interface{}{
			"message": "success persisting data",
			"data":    fmt.Sprintf("Transfered %d tokens from %s to %s.", form.Amount, form.Sender, form.Receiver),
		})
	} else {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": resp,
		})
	}
}
