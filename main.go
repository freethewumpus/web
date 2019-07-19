package main

import (
	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var RedisConnection *redis.Client
var Host string
var IndexHTML []byte
var RethinkConnection *r.Session

func main() {
	UtilsInit()
	rand.Seed(time.Now().Unix())

	data, err := ioutil.ReadFile("./index.html")
	if err != nil {
		panic(err)
	}
	IndexHTML = data

	Host = os.Getenv("HOST")
	if Host == "" {
		Host = "freethewump.us"
	}
	RedisHost := os.Getenv("REDIS_HOST")
	if RedisHost == "" {
		RedisHost = "localhost:6379"
	}
	RedisPassword := os.Getenv("REDIS_PASSWORD")
	if RedisPassword == "" {
		RedisPassword = ""
	}
	RedisConnection = redis.NewClient(&redis.Options{
		Addr: RedisHost,
		Password: RedisPassword,
		DB: 0,
	})

	RethinkHost := os.Getenv("RETHINK_HOST")
	if RethinkHost == "" {
		RethinkHost = "127.0.0.1:28015"
	}
	RethinkPass := os.Getenv("RETHINK_PASSWORD")
	RethinkUser := os.Getenv("RETHINK_USER")
	if RethinkUser == "" {
		RethinkUser = "admin"
	}
	s, err := r.Connect(r.ConnectOpts{
		Address: RethinkHost,
		Password: RethinkPass,
		Username: RethinkUser,
		Database: "freethewumpus",
	})
	if err != nil {
		panic(err)
	}
	RethinkConnection = s

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/", Upload)
	router.GET("/:image", View)

	log.Print("Always listening.")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", router))
}
