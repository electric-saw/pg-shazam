package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/raft"
)

type requestJoin struct {
	NodeID      string `json:"node_id"`
	RaftAddress string `json:"raft_address"`
}

func (h *Handler) JoinRaftHandler(c *gin.Context) {
	var from = requestJoin{}
	if err := c.Bind(from); err != nil {
		c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error binding %s", err.Error()),
		})
		return
	}

	if h.raft.State() != raft.Leader {
		c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": "not the leader",
		})
		return
	}

	config := h.raft.GetConfiguration()
	if err := config.Error(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("failed to get raft configuration: %s", err.Error()),
		})
		return
	}
	f := h.raft.AddVoter(raft.ServerID(from.NodeID), raft.ServerAddress(from.RaftAddress), 0, 0)
	if f.Error() != nil {
		c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"error": fmt.Sprintf("error add voter: %s", f.Error().Error()),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("node %s at %s joined successfully", from.NodeID, from.RaftAddress),
		"data":    h.raft.Stats(),
	})
}
