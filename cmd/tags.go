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

	subcommands := []*cobra.Command{tagsCommerceCmd, tagsNightlyCmd, tagsReleaseCmd}

	for i := 0; i < len(subcommands); i++ {
		subcommand := subcommands[i]

		tagsCmd.AddCommand(subcommand)
	}
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
		SubCommandInfo()
	},
}

var tagsCommerceCmd = &cobra.Command{
	Use:   "commerce",
	Short: "Lists all tags for Liferay Commerce Docker image",
	Long: `Lists all tags for Liferay Commerce Docker image from one of the unofficial, private repositories:
		- ` + liferay.CommercesRepository + ` (private).`,
	Run: func(cmd *cobra.Command, args []string) {
		commerce := liferay.Commerce{}

		log.Println("Sorry, but " + commerce.GetDockerHubTagsURL() + " repository is private, and we cannot access from this CLI :(")
	},
}

var tagsNightlyCmd = &cobra.Command{
	Use:   "nightly",
	Short: "Lists all tags for Liferay Portal Nightly Build Docker image",
	Long: `Lists all tags for Liferay Portal Nightly Build Docker image from one of the unofficial repository:
	- ` + liferay.NightliesRepository,
	Run: func(cmd *cobra.Command, args []string) {
		nightly := liferay.Nightly{}

		readTags(nightly)
	},
}

var tagsReleaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Lists all tags for Liferay Portal Release Docker image",
	Long: `Lists all tags for Liferay Portal Release Docker image from one of the unofficial repository:
	- ` + liferay.ReleasesRepository,
	Run: func(cmd *cobra.Command, args []string) {
		release := liferay.Release{}

		readTags(release)
	},
}

func readTags(image liferay.Image) {
	tagsPage := "https://hub.docker.com/r/" + image.GetDockerHubTagsURL() + "/tags/"

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
	availableSizes := []string{}

	maxLengthTags := 0

	// Find the review items
	doc.Find(dockerHubTagSearchQuery).Each(func(i int, selection *goquery.Selection) {
		// For each item found, get the tag
		tag := selection.Find("div .FlexTable__flexItemGrow2___3I1KN").Text()
		nodes := selection.Find("div .FlexTable__flexItemGrow1___3djP6")

		size := nodes.First().Text()

		currentTagLength := len(tag)
		if currentTagLength >= maxLengthTags {
			maxLengthTags = currentTagLength
		}

		availableTags = append(availableTags, tag)
		availableSizes = append(availableSizes, size)
	})

	if len(availableTags) > 0 {
		log.Printf("The available tags for the image are:")

		for index, tag := range availableTags {
			whitespacesCount := maxLengthTags - len(tag) + 6

			tagLine := tag

			for i := 0; i < whitespacesCount; i++ {
				tagLine += " "
			}

			tagLine += availableSizes[index]

			fmt.Println(tagLine)
		}
	} else {
		log.Printf("There are no available tags for the image")
	}
}
