package main

import (
	"log"
	"os"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/urfave/cli/v2"

	"github.com/fatih/color"
	"github.com/punnie/table"
)

var k = koanf.New(".")

var headerFmt = color.New(color.FgGreen, color.Underline).SprintfFunc()
var columnFmt = color.New(color.FgYellow).SprintfFunc()

func isNotPiped() bool {
	o, _ := os.Stdout.Stat()
	return (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice
}

func buildTable(headers ...interface{}) table.Table {
	tbl := table.New(headers...)

	if isNotPiped() {
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	} else {
		tbl.WithPrintHeaders(false)
	}

	return tbl
}

func main() {
	// Defaults
	k.Load(confmap.Provider(map[string]interface{}{
		"api.url":   "http://localhost:3000/api/v3",
		"api.token": "123123123",
	}, "."), nil)

	// Override from configuration file
	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		log.Printf("Error loading config: %v", err)
	}

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

							tbl := buildTable("Id", "Type", "Uri", "Title")

							for _, feed := range feeds.Feeds {
								tbl.AddRow(feed.Id, feed.Type, feed.Uri, feed.Title)
							}

							tbl.Print()

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

							tbl := buildTable("Id", "Uri")

							tbl.AddRow(feed.Id, feed.Uri)

							tbl.Print()

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

							tbl := buildTable("Id", "Name", "Permalink")

							for _, stream := range streams.Streams {
								tbl.AddRow(stream.Id, stream.Name, stream.Permalink)
							}

							tbl.Print()

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

							tbl := buildTable("Id", "Name", "Permalink")

							tbl.AddRow(stream.Id, stream.Name, stream.Permalink)

							tbl.Print()

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
							assignments, _ := ListStreamFeeds(ctx.String("stream_id"))

							tbl := buildTable("Id", "Feed Id", "Feed Type", "Feed Uri", "Feed Title")

							for _, assignment := range assignments.Feeds {
								tbl.AddRow(assignment.Id, assignment.FeedId, assignment.FeedType, assignment.FeedUri, assignment.FeedTitle)
							}

							tbl.Print()

							return nil
						},
					},
					{
						Name:  "add-feed",
						Usage: "add feed to a stream",

						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "stream_id",
								Aliases:  []string{"s"},
								Usage:    "id of the stream",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "feed_id",
								Aliases:  []string{"f"},
								Usage:    "id of the feed to add",
								Required: true,
							},
						},
						Action: func(ctx *cli.Context) error {
							assignment, _ := CreateStreamAssignment(ctx.String("stream_id"), ctx.String("feed_id"))

							tbl := buildTable("Id", "Feed Id", "Feed Type", "Feed Uri", "Feed Title")

							tbl.AddRow(assignment.Id, assignment.FeedId, assignment.FeedType, assignment.FeedUri, assignment.FeedTitle)

							tbl.Print()

							return nil
						},
					},
					{
						Name:  "remove-feed",
						Usage: "remove feed from a stream",

						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "stream_id",
								Aliases:  []string{"s"},
								Usage:    "id of the stream",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "feed_id",
								Aliases:  []string{"f"},
								Usage:    "id of the feed to remove",
								Required: true,
							},
						},
						Action: func(ctx *cli.Context) error {
							DestroyStreamAssignment(ctx.String("stream_id"), ctx.String("feed_id"))

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
