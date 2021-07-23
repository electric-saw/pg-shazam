package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) StatsRaftHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Here is the raft status",
		"data":    h.raft.Stats(),
	})
}
