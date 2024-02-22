/*
sysverity
Copyright (C) 2024  mutablefigment

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
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
