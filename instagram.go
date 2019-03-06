//
// Feed images from RSS feed into REDIS for output
//

package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis"
	"github.com/mmcdole/gofeed"
	// "github.com/mmcdole/gofeed/rss"
)

type InstagramImage struct {
	url   string
	image string
}

// Grab the RSS feed for instagram images
func main() {
	// file, _ := os.Open("feed.rss")
	// defer file.Close()

	redisClient := newClient()

	feedURL := "https://feeds.pinboard.in/rss/u:robballou/t:instagram/"
	res, _ := http.Get(feedURL)
	defer res.Body.Close()

	fp := gofeed.NewParser()
	feed, _ := fp.Parse(res.Body)
	for _, item := range feed.Items {
		document, err := goquery.NewDocument(item.Link)
		if err != nil {
			os.Exit(1)
		}
		image, imageErr := getImage(document, item.Link)
		if imageErr == nil {
			fmt.Println(image.image)
			redisClient.Set(image.url, image.image, 0)
		}
	}
}

func getImage(document *goquery.Document, url string) (InstagramImage, error) {
	tag := document.Find("meta[property=\"og:image\"]").First()
	href, exists := tag.Attr("content")
	if !exists {
		return InstagramImage{}, errors.New("Does not exist")
	}

	return InstagramImage{url, href}, nil
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
