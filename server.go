package beehiveAI

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

// prerequisite
//go get google.golang.org/grpc
//go get github.com/jackc/pgx/v4
//go get github.com/golang/protobuf/proto



package main

import (
"context"
"database/sql"
"fmt"
"log"
"net"
"time"

"github.com/golang/protobuf/ptypes"
"google.golang.org/grpc"
"reviews"
"github.com/jackc/pgx/v4"
)

const (
	databaseHost     = "localhost"
	databasePort     = 5432
	databaseName     = "your_database_name"
	databaseUser     = "your_username"
	databasePassword = "your_password"
)

type reviewServer struct{}

func (s *reviewServer) Search(ctx context.Context, filter *reviews.ReviewFilter) (*reviews.ReviewResponse, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", databaseUser, databasePassword, databaseHost, databasePort, databaseName)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	defer conn.Close(context.Background())

	var query string
	var args []interface{}

	query = "SELECT reviewer_name, title, text, rating, timestamp FROM reviews WHERE 1=1"

	if filter.ReviewerName != "" {
		query += " AND reviewer_name = $1"
		args = append(args, filter.ReviewerName)
	}

	if filter.MinRating > 0 {
		query += " AND rating >= $2"
		args = append(args, filter.MinRating)
	}

	if filter.MaxRating < 5 {
		query += " AND rating <= $3"
		args = append(args, filter.MaxRating)
	}

	if filter.MinTimestamp > 0 {
		query += " AND timestamp >= $4"
		args = append(args, time.Unix(filter.MinTimestamp, 0))
	}

	if filter.MaxTimestamp < time.Now().Unix() {
		query += " AND timestamp <= $5"
		args = append(args, time.Unix(filter.MaxTimestamp, 0))
	}

	rows, err := conn.Query(context.Background(), query, args...)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
	defer rows.Close()

	var reviewList []*reviews.Review

	for rows.Next() {
		var reviewerName, title, text string
		var rating int32
		var timestamp time.Time

		err := rows.Scan(&reviewerName, &title, &text, &rating, &timestamp)
		if err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}

		review := &reviews.Review{
			ReviewerName: reviewerName,
			Title:        title,
			Text:         text,
			Rating:       rating,
			Timestamp:    ptypes.TimestampNow(),
		}

		reviewList = append(reviewList, review)
	}

	return &reviews.ReviewResponse{
		Reviews: reviewList,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	reviews.RegisterReviewServiceServer(grpcServer, &reviewServer{})
	log.Println("gRPC server listening on port 50051...")
	grpcServer.Serve(lis)
}
