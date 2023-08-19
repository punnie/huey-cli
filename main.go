package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"encoding/json"
	"net/http"

	"github.com/knadh/koanf/providers/confmap"
	// "github.com/knadh/koanf/parsers/toml"
	// "github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/urfave/cli/v2"
)

var k = koanf.New(".")

func RequestApi[T any](verb string, path string, payloadData interface{}) (T, error) {
	jsonPayload, err := json.Marshal(payloadData)

	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	url := fmt.Sprintf("%s%s", k.String("api.url"), path)

	req, err := http.NewRequest(verb, url, bytes.NewBuffer(jsonPayload))

	if err != nil {
		log.Fatalf("Error creating Request: %v", err)
	}

	authorizationString := fmt.Sprintf("Bearer %s", k.String("api.token"))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorizationString)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error requesting stuff")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status: %s", resp.Status)
	}

	var responseData T

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		log.Fatalf("Error decoding response JSON: %v", err)
	}

	return responseData, nil
}

type Feed struct {
	Id string `json:"id"`
	Uri string `json:"uri"`
}

type ListFeedsResponse struct {
	Feeds []Feed `json:"feeds"`
}

func ListAllFeeds() (ListFeedsResponse, error) {
	result, err := RequestApi[ListFeedsResponse]("GET", "/feeds.json", nil)

	return result, err
}

type CreateFeedRequest struct {
	Feed Feed `json:"feed"`
}

func CreateFeed(uri string) (Feed, error) {
	payload := CreateFeedRequest{
		Feed: Feed{
			Uri: uri,
		},
	}

	result, err := RequestApi[Feed]("POST", "/feeds.json", payload)

	return result, err
}

type Stream struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type StreamResponse struct {
	Streams []Stream `json:"streams"`
}

func ListAllStreams() (StreamResponse, error) {
	result, err := RequestApi[StreamResponse]("GET", "/streams.json", nil)

	return result, err
}

func main() {
	// Defaults
	k.Load(confmap.Provider(map[string]interface{}{
		"api.url": "http://localhost:3000/api/v3",
		"api.token": "123123123",
	}, "."), nil)

	// Override from configuration file
	// if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
	// 	log.Printf("Error loading config: %v", err)
	// }

	// CLI interface
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "feeds",
				Aliases: []string{"f"},
				Usage: "manage feeds",
				Subcommands: []*cli.Command{
					{
						Name: "list",
						Usage: "list feeds",
						Action: func(ctx *cli.Context) error {
							feeds, _ := ListAllFeeds()

							for _, feed := range feeds.Feeds {
								fmt.Printf("%s %s\n", feed.Id, feed.Uri)
							}

							return nil
						},
					},
					{
						Name: "create",
						Usage: "create a feed",
						Flags: []cli.Flag{
							&cli.StringFlag {
								Name: "uri",
									Aliases: []string{"u"},
									Usage: "uri of the feed",
									Required: true,
								},
						},
						Action: func(ctx *cli.Context) error {
							feed, _ := CreateFeed(ctx.String("uri"))

							fmt.Printf("%s", feed.Id)

							return nil
						},
					},
				},
			},
			{
				Name: "streams",
				Aliases: []string{"s"},
				Usage: "manage streams",
				Subcommands: []*cli.Command{
					{
						Name: "list",
						Usage: "list streams",
						Action: func(ctx *cli.Context) error {
							streams, _ := ListAllStreams()

							for _, stream := range streams.Streams {
								fmt.Printf("%s %s\n", stream.Id, stream.Name)
							}

							return nil
						},
					},
				},
			},
		},
	}

	// Run program
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
