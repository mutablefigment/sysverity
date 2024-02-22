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
