package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	liferay "github.com/mdelapenya/lpn/liferay"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

const dockerHubTagSearchQuery = ".FlexTable__flexRow___2mqir"

func init() {
	rootCmd.AddCommand(tagsCmd)
}

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Lists all tags for Liferay Portal Docker image",
	Long: `Lists all tags for Liferay Portal Docker image from one of the unofficial repositories:
		- ` + liferay.CommercesRepository + ` (private),
		- ` + liferay.NightliesRepository + `, and
		- ` + liferay.ReleasesRepository + `.
	For that, please run this command adding "commerce", "release" or "nightly" subcommands.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errors.New("tags requires zero or one argument representing the image tag to be pulled")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		readTags("https://hub.docker.com/r/mdelapenya/liferay-portal/tags/")
	},
}

func readTags(tagsPage string) {
	// Request the HTML page.
	res, err := http.Get(tagsPage)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	availableTags := []string{}

	// Find the review items
	doc.Find(dockerHubTagSearchQuery).Each(func(i int, selection *goquery.Selection) {
		// For each item found, get the tag
		tag := selection.Find("div .FlexTable__flexItemGrow2___3I1KN").Text()

		availableTags = append(availableTags, tag)
	})

	if len(availableTags) > 0 {
		log.Printf("The available tags for the image are:")

		for _, tag := range availableTags {
			fmt.Println(tag)
		}
	} else {
		log.Printf("There are no available tags for the image")
	}
}
