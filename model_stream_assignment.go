package main

import (
	"fmt"
)

type StreamAssignmentParameters struct {
	StreamId string `json:"stream_id"`
	FeedId   string `json:"feed_id"`
}

type StreamAssignmentRequest struct {
	StreamAssignment StreamAssignmentParameters `json:"stream_assignment"`
}

type StreamAssignmentFeed struct {
	Id        string `json:"id"`
	FeedId    string `json:"feed_id"`
	FeedUri   string `json:"uri"`
	FeedTitle string `json:"title"`
	FeedType  string `json:"type"`
}

func CreateStreamAssignment(stream_id string, feed_id string) (StreamAssignmentFeed, error) {
	payload := StreamAssignmentRequest{
		StreamAssignment: StreamAssignmentParameters{
			StreamId: stream_id,
			FeedId:   feed_id,
		},
	}

	result, err := RequestApi[StreamAssignmentFeed]("POST", "/stream_assignments.json", payload)

	return result, err
}

type ListStreamFeedsResponse struct {
	Feeds []StreamAssignmentFeed
}

func ListStreamFeeds(id string) (ListStreamFeedsResponse, error) {
	uri := fmt.Sprintf("/streams/%s.json", id)
	result, err := RequestApi[ListStreamFeedsResponse]("GET", uri, nil)

	return result, err
}

func DestroyStreamAssignment(stream_id string, feed_id string) ([]interface{}, error) {
	stream_assignments, _ := ListStreamFeeds(stream_id)

	var results []interface{}

	for _, assignment := range stream_assignments.Feeds {
		if assignment.FeedId == feed_id {
			uri := fmt.Sprintf("/stream_assignments/%s.json", assignment.Id)
			result, _ := RequestApi[interface{}]("DELETE", uri, nil)

			results = append(results, result)
		}
	}

	return results, nil
}
