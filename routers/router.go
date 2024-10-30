package routers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"moderator/gpt"
	"moderator/models"
	"net/http"
	"strconv"
	"strings"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/validate", validate)

	return router
}

func validate(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "cannot read body"})
		return
	}

	var request models.Request
	err = json.Unmarshal(body, &request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	category := request.Topic
	requirements := fmt.Sprintf("проверь, что данное описание услуги подходит под категорию \"%s\". в названии должна быть только одна услуга/продукт итд. из других отраслей ничего быть не должно.\nесли не подходит, напиши 0, а на следующей строчке (важно!!) напиши причину одним предложением.\nесли подходит, напиши 1.", category)
	text := request.Text
	response, errResponse, err := gpt.Validate(requirements, text)
	if err != nil || errResponse != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "smth wrong happened while deciding",
		})
		return
	}
	result := strings.Split(response.Choices[0].Message.Content, "\n")
	result[0] = strings.TrimSpace(result[0])
	if num, err := strconv.Atoi(result[0]); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "cannot parse result. rerun request.",
			"result":  result,
		})
	} else {
		if num == 1 {
			c.JSON(http.StatusOK, gin.H{
				"relevant": 1,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"relevant": 0,
				"cause":    result[1],
			})
		}
	}
}
