package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// type messageResponse struct {
// 	Msg string `json:"msg"`
// }

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func main() {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery(), Cors())
	e.POST("/user/login", Posting)
	err := e.Run(":25580")
	if err != nil {
		log.Println(err)
	}
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// 发送post请求
func Posting(c *gin.Context) {
	var user User
	err := c.ShouldBindBodyWith(&user, binding.JSON)
	if err != nil {
		log.Println(err)
	}

	body, _ := json.Marshal(user)

	url := "http://192.168.192.2:25583/build"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
	}

	log.Println(req.Body)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}

	response, err := client.Do(req)
	if err != nil {

		log.Println(err)

	}

	reader := response.Body
	contentLength := response.ContentLength
	contentType := response.Header.Get("Content-Type")

	extraHeaders := map[string]string{
		"err": `json:"err"`,
		"out": `json:"out"`,
	}

	c.DataFromReader(response.StatusCode, contentLength, contentType, reader, extraHeaders)
}
