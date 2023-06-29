package lib

import "github.com/gin-gonic/gin"

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func Success(msg string, c *gin.Context) {
	c.JSON(200, gin.H{
		"message": msg,
	})

}
