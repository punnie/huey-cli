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

	if resp.StatusCode/100 != 2 {
		log.Fatalf("Request failed with status: %s", resp.Status)
	}

	var responseData T

	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		log.Fatalf("Error decoding response JSON: %v", err)
	}

	return responseData, nil
}

type Feed struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Uri   string `json:"uri"`
	Title string `json:"title"`
}

type ListFeedsResponse struct {
	Feeds []Feed `json:"feeds"`
}

func ListAllFeeds() (ListFeedsResponse, error) {
	result, err := RequestApi[ListFeedsResponse]("GET", "/feeds.json", nil)

	return result, err
}

type CreateFeedParameters struct {
	Uri               string `json:"uri"`
	UseGooglebotAgent bool   `json:"use_googlebot_agent,omitempty"`
}

type CreateFeedRequest struct {
	Feed CreateFeedParameters `json:"feed"`
}

func CreateFeed(uri string) (Feed, error) {
	payload := CreateFeedRequest{
		Feed: CreateFeedParameters{
			Uri: uri,
		},
	}

	result, err := RequestApi[Feed]("POST", "/feeds.json", payload)

	return result, err
}

type Stream struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type StreamResponse struct {
	Streams []Stream `json:"streams"`
}

func ListAllStreams() (StreamResponse, error) {
	result, err := RequestApi[StreamResponse]("GET", "/streams.json", nil)

	return result, err
}

type CreateStreamParameters struct {
	Name string `json:"name"`
	Permalink string `json:"permalink,omitempty"`
}

type CreateStreamRequest struct {
	Stream CreateStreamParameters `json:"stream"`
}

func CreateStream(name string, permalink string) (Stream, error) {
	payload := CreateStreamRequest{
		Stream: CreateStreamParameters{
			Name: name,
			Permalink: permalink,
		},
	}

	result, err := RequestApi[Stream]("POST", "/streams.json", payload)

	return result, err
}

func ListStreamFeeds(id string) (ListFeedsResponse, error) {
	uri := fmt.Sprintf("/streams/%s.json", id)
	result, err := RequestApi[ListFeedsResponse]("GET", uri, nil)

	return result, err
}

type StreamAssignmentParameters struct {
	StreamId string `json:"stream_id"`
	FeedId   string `json:"feed_id"`
}

type StreamAssignmentRequest struct {
	StreamAssignment StreamAssignmentParameters `json:"stream_assignment"`
}

func CreateStringAssignment(stream_id string, feed_id string) (Stream, error) {
	payload := StreamAssignmentRequest{
		StreamAssignment: StreamAssignmentParameters{
			StreamId: stream_id,
			FeedId:   feed_id,
		},
	}

	result, err := RequestApi[Stream]("POST", "/stream_assignments.json", payload)

	return result, err
}

func main() {
	// Defaults
	k.Load(confmap.Provider(map[string]interface{}{
		"api.url":   "http://localhost:3000/api/v3",
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
				Name:    "feeds",
				Aliases: []string{"f"},
				Usage:   "manage feeds",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "list feeds",
						Action: func(ctx *cli.Context) error {
							feeds, _ := ListAllFeeds()

							for _, feed := range feeds.Feeds {
								fmt.Printf("%-16s %-10s %-48s %s\n", feed.Id, feed.Type, feed.Uri, feed.Title)
							}

							return nil
						},
					},
					{
						Name:  "create",
						Usage: "create a feed",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "uri",
								Aliases:  []string{"u"},
								Usage:    "uri of the feed",
								Required: true,
							},

							&cli.BoolFlag{
								Name:    "use-googlebot-agent",
								Aliases: []string{"g"},
								Usage:   "use googlebot agent",
							},
						},
						Action: func(ctx *cli.Context) error {
							feed, _ := CreateFeed(ctx.String("uri"))

							fmt.Printf("%s\n", feed.Id)

							return nil
						},
					},
				},
			},
			{
				Name:    "streams",
				Aliases: []string{"s"},
				Usage:   "manage streams",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "list streams",
						Action: func(ctx *cli.Context) error {
							streams, _ := ListAllStreams()

							for _, stream := range streams.Streams {
								fmt.Printf("%s %s\n", stream.Id, stream.Name)
							}

							return nil
						},
					},
					{
						Name:  "create",
						Usage: "create a stream",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Aliases:  []string{"n"},
								Usage:    "name of the stream",
								Required: true,
							},

							&cli.StringFlag{
								Name:    "permalink",
								Aliases: []string{"p"},
								Usage:   "permalink for the stream",
							},
						},
						Action: func(ctx *cli.Context) error {
							stream, _ := CreateStream(ctx.String("name"), ctx.String("permalink"))

							fmt.Printf("%s\n", stream.Id)

							return nil
						},
					},
					{
						Name:  "list-feeds",
						Usage: "list feeds in stream",

						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "stream_id",
								Aliases:  []string{"s"},
								Usage:    "id of the stream",
								Required: true,
							},
						},
						Action: func(ctx *cli.Context) error {
							feeds, _ := ListStreamFeeds(ctx.String("stream_id"))

							for _, feed := range feeds.Feeds {
								fmt.Printf("%-16s %-10s %-48s %s\n", feed.Id, feed.Type, feed.Uri, feed.Title)
							}

							return nil
						},
					},
					{
						Name:  "add-feeds",
						Usage: "add feeds to a stream",

						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "stream_id",
								Aliases:  []string{"s"},
								Usage:    "id of the stream",
								Required: true,
							},
							&cli.StringSliceFlag{
								Name:     "feed_id",
								Aliases:  []string{"s"},
								Usage:    "id of the feeds to add",
								Required: true,
							},
						},
						Action: func(ctx *cli.Context) error {
							streams, _ := ListAllStreams()

							for _, stream := range streams.Streams {
								fmt.Printf("%s %s\n", stream.Id, stream.Name)
							}

							return nil
						},
					},
					{
						Name:  "remove-feeds",
						Usage: "remove feeds from a stream",

						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "stream_id",
								Aliases:  []string{"s"},
								Usage:    "id of the stream",
								Required: true,
							},
							&cli.StringSliceFlag{
								Name:     "feed_id",
								Aliases:  []string{"s"},
								Usage:    "id of the feeds to remove",
								Required: true,
							},
						},
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
