package main

import (
	"beehiveAI/reviews"
	"context"
	"testing"
	"time"
)

type mockReviewServer struct{}

func (m *mockReviewServer) Search(ctx context.Context, filter *reviews.ReviewFilter) (*reviews.ReviewResponse, error) {
	// test cases
	review1 := &reviews.Review{
		ReviewerName: "John Doe",
		Title:        "Test Review 1",
		Text:         "This is a test review 1.",
		Rating:       4,
		Timestamp:    time.Now().Unix(),
	}

	review2 := &reviews.Review{
		ReviewerName: "Jane Smith",
		Title:        "Test Review 2",
		Text:         "This is a test review 2.",
		Rating:       5,
		Timestamp:    time.Now().Add(-time.Hour).Unix(),
	}

	reviewList := []*reviews.Review{review1, review2}

	return &reviews.ReviewResponse{
		Reviews: reviewList,
	}, nil
}

func TestSearch(t *testing.T) {
	server := &mockReviewServer{}
	// Test case 1: Filter by reviewer name
	filter1 := &reviews.ReviewFilter{ReviewerName: "J. McDonald"}
	response1, err := server.Search(context.Background(), filter1)
	if err != nil {
		t.Errorf("Search with reviewer name filter failed: %v", err)
	}
	if len(response1.Reviews) == 0 {
		t.Errorf("Expected non-empty review list for reviewer name filter")
	}

	// Filter by rating range
	filter2 := &reviews.ReviewFilter{MinRating: 3, MaxRating: 5}
	response2, err := server.Search(context.Background(), filter2)
	if err != nil {
		t.Errorf("Search with rating range filter failed: %v", err)
	}
	if len(response2.Reviews) == 0 {
		t.Errorf("Expected non-empty review list for rating range filter")
	}

	// Filter by timestamp range
	minTimestamp := time.Now().Add(-time.Hour).Unix()
	maxTimestamp := time.Now().Unix()
	filter3 := &reviews.ReviewFilter{MinTimestamp: minTimestamp, MaxTimestamp: maxTimestamp}
	response3, err := server.Search(context.Background(), filter3)
	if err != nil {
		t.Errorf("Search wit	h timestamp range filter failed: %v", err)
	}
	if len(response3.Reviews) == 0 {
		t.Errorf("Expected non-empty review list for timestamp range filter")
	}
}
