package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Comic is comic type
type Comic struct {
	_id            string
	CrawledFrom    string
	RecognizedID   string
	GlobalID       string
	Title          string
	Category       string
	ImageURL       string
	Link           string
	ComicUpdatedAt time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

// SaveComic is to save a comic
func SaveComic(comic Comic) (*mongo.UpdateResult, error) {
	return Collections.Comics.UpdateOne(
		context.Background(),
		bson.D{
			{"crawled_from", comic.CrawledFrom},
			{"recognized_id", comic.RecognizedID},
		},
		bson.D{
			{
				"$set",
				bson.D{
					{"title", comic.Title},
					{"crawled_from", comic.CrawledFrom},
					{"category", comic.Category},
					{"image_url", comic.ImageURL},
					{"link", comic.Link},
					{"comic_updated_at", comic.ComicUpdatedAt},
					{"updated_at", time.Now()},
				},
			},
			{
				"$setOnInsert",
				bson.D{
					{"global_id", nil},
					{"created_at", time.Now()},
					{"deleted_at", nil},
				},
			},
		},
		options.Update().SetUpsert(true),
	)
}
