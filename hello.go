package main

import (
	"example/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	n := 5
	fmt.Println(utils.MagicSum(n))
	fmt.Println(utils.MagicPow(n))
	fmt.Println(utils.MagicOdd(n))
	fmt.Println(utils.MagicGrade(n))
	fmt.Println(utils.MagicName(n))
	fmt.Println(utils.MagicTria(n))
	utils.MagicChange(&n)
	fmt.Println(n)
	test()
}

func test() {
	r := gin.Default()
	r.GET("/try", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "hahahihi",
		})
	})
	r.Run()
}
