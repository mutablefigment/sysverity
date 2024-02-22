package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mutablefigment/sysverity/v2/pkg/block"
)

type BlockHandler struct {
	store blockStore
}

type blockStore interface {
	Genisis(b *block.Block) error
	LastHash(b *block.Block) (string, error)
}

func (h BlockHandler) CreateGensisiBlock(c *gin.Context) {}

func NewGenisisBlock(s blockStore) *BlockHandler {
	return &BlockHandler{
		store: s,
	}
}
