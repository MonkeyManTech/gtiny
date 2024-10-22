package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	_ "net/http"
	"time"
)

type Database interface {
	Open() error
}

type SqliteDatabase struct {
}

func (db *SqliteDatabase) Open() error {
	return nil
}

type KeyService struct {
	DB Database
}

func (u *KeyService) SaveShortKey(shortKey string, originalUrl string) error {
	//stmt := fmt.Sprintf("insert into urlShortKeys(shortKey, originalUrl) values(%s, %s)", shortKey, originalUrl)
	return nil
}

type ShortReq struct {
	Url string `json:"url"`
}

type UrlHandler struct {
	urls map[string]string
}

func (h *UrlHandler) HandlePost(c *gin.Context) {
	var req ShortReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sk := generateShortKey()
	h.urls[sk] = req.Url
	c.JSON(http.StatusCreated, gin.H{"key": sk})
}

func NewUrlHandler() *UrlHandler {
	return &UrlHandler{urls: make(map[string]string)}
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sk := make([]byte, keyLength)
	for i := range sk {
		sk[i] = charset[r.Intn(len(charset))]
	}
	return string(sk)
}

func main() {
	db, err := sql.Open("sqlite3", "./gtiny.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	h := NewUrlHandler()
	r := gin.Default()
	r.POST("/shorten", h.HandlePost)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
