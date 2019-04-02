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
	_id            string    `bson:"_id"`
	CrawledFrom    string    `bson:"crawled_from"`
	RecognizedID   string    `bson:"recognized_id"`
	GlobalID       string    `bson:"global_id"`
	Title          string    `bson:"title"`
	Category       string    `bson:"category"`
	ImageURL       string    `bson:"image_url"`
	Link           string    `bson:"link"`
	ComicUpdatedAt time.Time `bson:"comic_updated_at"`
	CreatedAt      time.Time `bson:"created_at"`
	UpdatedAt      time.Time `bson:"updated_at"`
	DeletedAt      time.Time `bson:"deleted_at"`
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
