package main


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
