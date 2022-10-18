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

// Deposit handling query a balance of one's account
func (h handler) Get(eCtx echo.Context) error {
	Account := eCtx.Param("account")

	if Account == "" {
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
		Operation: "GET",
		Account:   Account,
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

	if resp.Error != nil {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("The account '%s' does not exist.", Account),
		})
	} else {
		return eCtx.JSON(http.StatusOK, map[string]interface{}{
			"message": "success persisting data",
			"data":    fmt.Sprintf("The account '%s' has %d tokens.", Account, resp.Data),
		})
	}
}
