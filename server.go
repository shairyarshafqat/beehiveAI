package main

import (
	"beehiveAI/reviews"
	"context"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

const (
	//databaseURL = "postgres://your_username:your_password@localhost:5432/your_database_name" // Update with your database URL
	databaseURL = "postgres://shairyar:kickme1@localhost:5432/amazon_reviews"
)

type ReviewServiceServer interface {
	Search(ctx context.Context, filter *reviews.ReviewFilter) (*reviews.ReviewResponse, error)
}

type reviewServer struct {
	reviews.UnimplementedReviewServiceServer
}

func (s *reviewServer) Search(ctx context.Context, filter *reviews.ReviewFilter) (*reviews.ReviewResponse, error) {
	//connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", databaseUser, databasePassword, databaseHost, databasePort, databaseName)

	var conn, err = pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	var query string
	var args []interface{}

	query = "SELECT reviewer_name, title, text, rating, timestamp FROM reviews WHERE 1=1"

	if filter.ReviewerName != "" {
		query += " AND reviewer_name = $1"
		args = append(args, filter.ReviewerName)
	}

	if filter.MinRating > 0 {
		query += " AND rating >= $2"
		args = append(args, pgtype.Int4{Int: int32(filter.MinRating)})
	}

	if filter.MaxRating < 5 {
		query += " AND rating <= $3"
		args = append(args, pgtype.Int4{Int: int32(filter.MaxRating)})
	}

	if filter.MinTimestamp > 0 {
		query += " AND timestamp >= $4::timestamp"
		args = append(args, filter.MinTimestamp)
	}

	if filter.MaxTimestamp < time.Now().Unix() {
		query += " AND timestamp <= $5::timestamp"
		args = append(args, filter.MaxTimestamp)
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
			Timestamp:    timestamp.Unix(),
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
	reflection.Register(grpcServer)
	reviews.RegisterReviewServiceServer(grpcServer, &reviewServer{})
	log.Println("gRPC server listening on port 50051...")
	grpcServer.Serve(lis)
}
