package handler

import (
	"github.com/hashicorp/raft"
)

type Handler struct {
	raft *raft.Raft
}

func New(raft *raft.Raft) *Handler {
	return &Handler{raft: raft}
}
