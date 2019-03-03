package db

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/options"
)

// Comic is comic type
type Comic struct {
	ComicID       string
	Name          string
	Category      string
	LastUpdatedAt time.Time
}

// SaveComic is to save a comic
func SaveComic(comic Comic) (*mongo.UpdateResult, error) {
	return Collections.Comics.UpdateOne(
		context.Background(),
		bson.D{
			{"comicID", comic.ComicID},
		},
		bson.D{
			{
				"$set",
				bson.D{
					{"name", comic.Name},
					{"category", comic.Category},
					{"lastUpdatedAt", comic.LastUpdatedAt},
				},
			},
		},
		options.Update().SetUpsert(true),
	)
}
