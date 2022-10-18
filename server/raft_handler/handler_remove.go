package raft_handler

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/raft"
	"github.com/labstack/echo"
)

// requestRemove request payload for removing node from raft cluster

// RemoveRaftHandler handling removing raft
func (h handler) RemoveRaftHandler(eCtx echo.Context) error {
	nodeID := eCtx.Param("node_id")

	if h.raft.State() != raft.Leader {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": "not the leader",
		})
	}

	configFuture := h.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("failed to get raft configuration: %s", err.Error()),
		})
	}

	future := h.raft.RemoveServer(raft.ServerID(nodeID), 0, 0)
	if err := future.Error(); err != nil {
		return eCtx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error removing existing node %s: %s", nodeID, err.Error()),
		})
	}

	return eCtx.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("node %s removed successfully", nodeID),
	})
}
