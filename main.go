package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/ggary9424/comic-crawler/db"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

// HTMLElement is a HTMLElement
type HTMLElement = colly.HTMLElement

// Response is a type of colly library
type Response = colly.Response

func getEnv() string {
	if os.Getenv("APP_ENV") == "" {
		return "local"
	}
	return os.Getenv("APP_ENV")
}

func init() {
	// Initialize environment variables
	viper.SetConfigName(getEnv())
	viper.AddConfigPath("./configs/")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// Initialize logger system
	debug := viper.GetBool("system.debug")

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if getEnv() != "local" {
		// Some log event collection services
	}
	log.SetOutput(os.Stdout)
}

func main() {
	var (
		dbHost = viper.GetString("database.host")
		dbPort = viper.GetString("database.port")
		dbURL  = "mongodb://" + dbHost + ":" + dbPort
		dbName = viper.GetString("database.name")
	)

	_, e := db.Connect(dbURL, dbName)

	if e != nil {
		log.WithFields(log.Fields{
			"dbURL":  dbURL,
			"dbName": dbName,
		}).Fatal("Connect to DB fail. Please check your configuration or DB status.")
		return
	}

	comicIDs := []string{"1878", "3583"}
	collyController := colly.NewCollector()

	collyController.OnHTML("body", func(body *HTMLElement) {
		reg, _ := regexp.Compile("https:\\/\\/www\\.cartoonmad\\.com\\/comic\\/(\\d+)")
		submatch := reg.FindStringSubmatch(body.Request.URL.String())
		comic := db.Comic{ComicID: submatch[1]}
		body.ForEach("body > table > tbody > tr > td:nth-child(2) > table > tbody > tr", func(_ int, tr *HTMLElement) {
			if tr.Index == 2 {
				tr.ForEach("td:nth-child(2) a", func(i int, a *HTMLElement) {
					if i == 1 {
						// Convert from big5 to utf8
						utf8String, _, _ := transform.String(traditionalchinese.Big5.NewDecoder(), a.Text)
						comic.Category = utf8String
					}
					if i == 2 {
						// Convert from big5 to utf8
						utf8String, _, _ := transform.String(traditionalchinese.Big5.NewDecoder(), a.Text)
						comic.Name = utf8String
					}
				})
			}

			if tr.Index == 3 {
				tr.ForEach("td[background=\"/image/content_box2.gif\"] b font", func(_ int, font *HTMLElement) {
					if font.Index == 0 {
						reg, _ := regexp.Compile("\\d+/\\d+/\\d+ \\d+:\\d+:\\d+ (AM|PM)")
						dateString := reg.FindString(font.Text)
						lastUpdatedAt, _ := time.Parse("1/2/2006 3:04:05 PM", dateString)
						comic.LastUpdatedAt = lastUpdatedAt
					}
				})
			}
		})

		log.WithFields(log.Fields{
			"comicID":       comic.ComicID,
			"name":          comic.Name,
			"category":      comic.Category,
			"lastUpdatedAt": comic.LastUpdatedAt,
		}).Debug("Log comic content.")

		db.SaveComic(comic)
	})

	for _, comicID := range comicIDs {
		collyController.Visit("https://www.cartoonmad.com/comic/" + comicID + ".html")
	}
}
