package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type sbTask struct {
	Task    task `json:"task"`
	SpaceID int  `json:"space_id"`
}

type task struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	mux := httprouter.New()
	mux.GET("/", hello)
	mux.POST("/go", handleTask)

	http.ListenAndServe(":"+os.Getenv("PORT"), mux)
}

func hello(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	io.WriteString(w, "Hello from Docker container")
}

func handleTask(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var task sbTask

	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		log.Println(err)
	}

	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}

	client, err := getClient(&creds)
	if err != nil {
		log.Println("Error getting twitter client:", err)
	}

	tweet, res, err := client.Statuses.Update("Hello from my automated tweet bot", nil)
	if err != nil {
		log.Print(err)
	}
	log.Printf("%+v\n", res)
	log.Printf("%+v\n", tweet)

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Println(err)
	}
}
