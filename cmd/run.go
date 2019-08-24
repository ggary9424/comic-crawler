package cmd

import (
	"regexp"
	"time"

	"github.com/ggary9424/comic-crawler/db"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start to crawl comic content.",
	Long:  `Start to crawl comic content.`,
	Run: func(cmd *cobra.Command, args []string) {
		startRunningCrawler()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

// HTMLElement is a HTMLElement
type HTMLElement = colly.HTMLElement

// Response is a type of colly library
type Response = colly.Response

func startRunningCrawler() {
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

	comicIDs := []string{"1152", "2504", "1878", "3583", "4085", "1698", "5200", "4982", "3899", "1221"}
	collyController := colly.NewCollector()

	collyController.OnHTML("body", func(body *HTMLElement) {
		reg, _ := regexp.Compile("https:\\/\\/www\\.cartoonmad\\.com\\/comic\\/(\\d+)")
		submatch := reg.FindStringSubmatch(body.Request.URL.String())

		var comicTitle string
		var comicCategory string
		var comicUpdatedAt time.Time
		body.ForEach("body > table > tbody > tr > td:nth-child(2) > table > tbody > tr", func(_ int, tr *HTMLElement) {
			if tr.Index == 2 {
				tr.ForEach("td:nth-child(2) a", func(i int, a *HTMLElement) {
					if i == 1 {
						// Convert from big5 to utf8
						utf8String, _, _ := transform.String(traditionalchinese.Big5.NewDecoder(), a.Text)
						comicCategory = utf8String
					}
					if i == 2 {
						// Convert from big5 to utf8
						utf8String, _, _ := transform.String(traditionalchinese.Big5.NewDecoder(), a.Text)
						comicTitle = utf8String
					}
				})
			}

			if tr.Index == 3 {
				tr.ForEach("td[background=\"/image/content_box2.gif\"] b font", func(_ int, font *HTMLElement) {
					if font.Index == 0 {
						reg, _ := regexp.Compile("\\d+/\\d+/\\d+ \\d+:\\d+:\\d+ (AM|PM)")
						dateString := reg.FindString(font.Text)
						comicUpdatedAt, _ = time.Parse("1/2/2006 3:04:05 PM", dateString)
					}
				})
			}
		})

		comic := db.Comic{
			CrawledFrom:    "動漫狂",
			RecognizedID:   submatch[1],
			Title:          comicTitle,
			Category:       comicCategory,
			ImageURL:       "https://www.cartoonmad.com/cartoonimgs/coimg/" + submatch[1] + ".jpg",
			ComicUpdatedAt: comicUpdatedAt,
			Link:           body.Request.URL.String(),
		}

		log.WithFields(log.Fields{
			"CrawledFrom":    comic.CrawledFrom,
			"RecognizedID":   comic.RecognizedID,
			"Title":          comic.Title,
			"Category":       comic.Category,
			"ComicUpdatedAt": comic.ComicUpdatedAt,
		}).Debug("Log comic content.")

		_, err := db.SaveComic(comic)
		if err != nil {
			log.Fatal(err)
		}
	})

	for _, comicID := range comicIDs {
		collyController.Visit("https://www.cartoonmad.com/comic/" + comicID + ".html")
	}
}
