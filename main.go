package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, PUT")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

type Token struct {
	Access_token string `json:"access_token"`
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token Token
		c.ShouldBindBodyWith(&token, binding.JSON)
		if token.Access_token != con.Access_token {
			log.Println()
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		c.Next()
	}
}

func main() {
	log.Println(con)
	Init(&con)
	log.Println(con)
	log.SetFlags(log.Ldate | log.Lshortfile)
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery(), Cors(), Auth())
	e.POST("/redeploy", Posting)
	err := e.Run(":" + con.Port)
	if err != nil {
		log.Println(err)
	}
}

func Init(c interface{}) {
	con.Access_key = os.Getenv("ACCESS_KEY")
	con.Access_token = os.Getenv("ACCESS_TOKEN")
	con.Port = os.Getenv("PORT")
	con.Url = os.Getenv("URL")
}

var con Config

type Config struct {
	Port         string `env:"PORT"`
	Access_key   string `env:"ACCESS_KEY"`
	Access_token string `env:"ACCESS_TOKEN"`
	Url          string `env:"URL"`
}

type ClientMassage struct {
	ExecOut string `json:"execOut"`
	ExecErr string `json:"execErr"`
}

// 发送post请求
func Posting(c *gin.Context) {
	err := c.ShouldBindBodyWith(&con, binding.JSON)
	if err != nil {
		log.Println(err)
	}

	payload := strings.NewReader("{\"access_key\":\"" + con.Access_key + "\"}")
	req, err := http.NewRequest("PUT", con.Url, payload)
	if err != nil {
		log.Println(err)
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer response.Body.Close()
	var resp []byte
	resp, err = io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}

	var msg ClientMassage
	err = json.Unmarshal(resp, &msg)
	if err != nil {
		log.Println(err)
	}

	c.JSON(200, gin.H{
		"out": resp,
	})
	log.Println(resp)
}
