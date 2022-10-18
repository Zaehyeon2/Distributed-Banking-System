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

type requestDeposit struct {
	Account string `json:"account"`
	Amount  uint64 `json:"amount"`
}

func (h handler) Deposit(eCtx echo.Context) error {
	var form = requestDeposit{}
	if err := eCtx.Bind(&form); err != nil {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error binding: %s", err.Error()),
		})
	}

	if form.Account == "" {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": "account is empty",
		})
	}

	if h.raft.State() != raft.Leader {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": "not the leader",
		})
	}

	payload := fsm.CommandPayload{
		Operation: "DEPOSIT",
		Receiver:  form.Account,
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

	_, ok := applyFuture.Response().(*fsm.ApplyResponse)
	if !ok {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error response is not match apply response"),
		})
	}

	return eCtx.JSON(http.StatusOK, map[string]interface{}{
		"message": "success persisting data",
		"data":    fmt.Sprintf("Deposited %d tokens to %s.", form.Amount, form.Account),
	})
}
