package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

// Image a cool instagram image
type Image struct {
	URL   string `json:"url"`
	Image string `json:"image"`
}

// GetImages get 20 random images
func GetImages(writer http.ResponseWriter, request *http.Request) {
	redisClient := newClient()

	var images []Image
	count := 20
	for i := 0; i < count; i++ {
		key, err := redisClient.RandomKey().Result()
		if err != nil {
			panic(err)
		}

		value, err := redisClient.Get(key).Result()
		if err != nil {
			panic(err)
		}
		images = append(images, Image{URL: key, Image: value})
	}

	// json.NewEncoder(writer).Encode(images)
	// response, _ := json.Marshal(images)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(writer).Encode(images); err != nil {
		panic(err)
	}
}

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return client
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/images", GetImages).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
