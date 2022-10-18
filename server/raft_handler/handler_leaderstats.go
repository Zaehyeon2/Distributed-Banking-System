package raft_handler

import (
	"net/http"

	"github.com/hashicorp/raft"
	"github.com/labstack/echo"
)

type leaderStats struct {
	Address raft.ServerAddress
	Id      raft.ServerID
}

// StatsLeaderHandler get leader's status
func (h handler) StatsLeaderHandler(eCtx echo.Context) error {
	address, ID := h.raft.LeaderWithID()
	return eCtx.JSON(http.StatusOK, map[string]interface{}{
		"message": "Here is the raft status",
		"data":    &leaderStats{Address: address, Id: ID},
	})
}
