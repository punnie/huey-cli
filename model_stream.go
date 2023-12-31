package main

type Stream struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Permalink string `json:"permalink"`
}

type StreamResponse struct {
	Streams []Stream `json:"streams"`
}

func ListAllStreams() (StreamResponse, error) {
	result, err := RequestApi[StreamResponse]("GET", "/streams.json", nil)

	return result, err
}

type CreateStreamParameters struct {
	Name      string `json:"name"`
	Permalink string `json:"permalink,omitempty"`
}

type CreateStreamRequest struct {
	Stream CreateStreamParameters `json:"stream"`
}

func CreateStream(name string, permalink string) (Stream, error) {
	payload := CreateStreamRequest{
		Stream: CreateStreamParameters{
			Name:      name,
			Permalink: permalink,
		},
	}

	result, err := RequestApi[Stream]("POST", "/streams.json", payload)

	return result, err
}
