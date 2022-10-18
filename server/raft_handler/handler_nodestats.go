package raft_handler

import (
	"net/http"

	"github.com/labstack/echo"
)

// StatsRaftHandler get raft status
func (h handler) StatsNodesHandler(eCtx echo.Context) error {
	data := h.raft.Stats()
	return eCtx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Here is the nodes status in the cluster",
		"data":    data["latest_configuration"],
	})
}
