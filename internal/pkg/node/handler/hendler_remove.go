package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/raft"
)

type requestRemove struct {
	NodeID string `json:"node_id"`
}

func (h *Handler) RemoveRaftHandler(c *gin.Context) {
	var form = requestRemove{}
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error binding: %s", err.Error()),
		})
		return
	}

	var nodeID = form.NodeID

	if h.raft.State() != raft.Leader {
		c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": "not the leader",
		})
	}

	configFuture := h.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("failed to get raft configuration: %s", err.Error()),
		})
		return
	}

	future := h.raft.RemoveServer(raft.ServerID(nodeID), 0, 0)
	if err := future.Error(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error removing existing node %s: %s", nodeID, err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("node %s removed successfully", nodeID),
		"data":    h.raft.Stats(),
	})
}
