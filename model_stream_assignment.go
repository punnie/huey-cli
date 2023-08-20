package main

type StreamAssignmentParameters struct {
	StreamId string `json:"stream_id"`
	FeedId   string `json:"feed_id"`
}

type StreamAssignmentRequest struct {
	StreamAssignment StreamAssignmentParameters `json:"stream_assignment"`
}

func CreateStreamAssignment(stream_id string, feed_id string) (ListFeedsResponse, error) {
	payload := StreamAssignmentRequest{
		StreamAssignment: StreamAssignmentParameters{
			StreamId: stream_id,
			FeedId:   feed_id,
		},
	}

	result, err := RequestApi[ListFeedsResponse]("POST", "/stream_assignments.json", payload)

	return result, err
}
