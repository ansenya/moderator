package routers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"log"
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
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "cannot parse body"})
		return
	}
	if request.Title == "" || request.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "title or description is empty"})
		return
	}

	messages := []models.Message{
		{
			"system",
			[]models.Content{
				{
					Type: "text",
					Text: "ты проверяешь услугу по критериям\n\n" +
						"критерий для заголовка: заголовок должен соответствовать написанию\n" +
						"критерии для описания услуги:\n" +
						"1. в описании услуги не может быть цена отличная от той, что указана выше\n" +
						"2. в описании нет ссылок на сайты, мессенджеры, соцсети\n" +
						"3. текст описания отражает услуги, которые можно заказать. описание не должно быть предложением продать работу и/или запросом на выполнение задач\n" +
						"4. в описании нет публикаций услуг, направленных на нарушение законодательства, нормы морали и этики, а также законодательные ограничения\n" +
						"5. описание услуги соответствует заголовку\n\n" +
						"если услуга совсем не подходит по критериям, напиши \"-1\" и следующей строчкой напиши причину\n" +
						"если услуга частично не подходит по критериям и/или нарушает критерии не сильным образом (не критично), напиши \"0\" и следующей строчкой напиши причину\n" +
						"если услуга соответствует критерям, напиши \"1\"\n\n" +
						"причину пиши только одним предложением, не нужно писать \"Причина:...\"",
				},
				{
					Type: "text",
					Text: fmt.Sprintf(
						"цена: %d\n"+
							"заголовок: %s\n"+
							"описание:%s", request.Price, request.Title, request.Description),
				},
			},
		},
	}

	response, errResponse, err := gpt.Validate(messages)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "smth wrong happened while getting response",
		})
		return
	}
	if errResponse != nil {
		log.Println(errResponse)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "smth wrong happened getting response",
		})
		return
	}

	result := strings.Split(response.Choices[0].Message.Content, "\n")
	result[0] = strings.TrimSpace(result[0])
	if num, err := strconv.Atoi(result[0]); err != nil {
		log.Println(err)
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
				"relevant": num,
				"cause":    result[1],
			})
		}
	}
}
