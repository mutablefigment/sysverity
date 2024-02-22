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
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mutablefigment/sysverity/v2/pkg/block"
)

type Register struct {
	Pubkey         string
	HashChainToken string
	BootHash       string
}

func main() {
	fmt.Print("hello server\r\n")

	b := block.Block{}
	block.PrettyPrint(&b)

	block.CacluateCurrBlockHash(&b)
	block.CacluatePrevBlockHash(&b)

	router := gin.Default()

	router.POST("/register", register)

	router.Run(":5100")
}

func register(c *gin.Context) {
	var requestBody Register
	if err := c.BindJSON(&requestBody); err != nil {
		fmt.Println(err)
	}

	fmt.Println(requestBody.Pubkey)
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "This is a test page")
}
