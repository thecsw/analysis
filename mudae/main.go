package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/thecsw/analysis/models/discord"

	"github.com/dustin/go-humanize"
)

const (
	// MudaeID is the ID of the Mudae bot.
	MudaeID = "432610292342587392"

	// MudaeGachaInteraction is the name of the interaction that Mudae uses to
	// respond to gacha requests.
	MudaeGachaInteraction = "wa"

	// Location is the location of the discord server.
	Location = "America/Chicago"
)

var (
	// Regexp to extract the claim rank from the embed description.
	claimRankRegexp = regexp.MustCompile(`(?U)Claims: #([0-9]+)`)
	// Regexp to extract the likes rank from the embed description.
	likesRankRegexp = regexp.MustCompile(`(?U)Likes: #([0-9]+)`)
	// Regexp to extract the kakera price from the embed description.
	kakeraPriceRegexp = regexp.MustCompile(`(?U)\*\*([0-9]+)\*\*:kakera`)
	// Regexp to extract the show title from the embed description.
	showTitleRegexp = regexp.MustCompile(`(?U)([^/]+)\n`)
)

func main() {
	exportFilename := flag.String("export", "", "discord export file")
	outputFilename := flag.String("output", "output.csv", "where to output")
	printTopShows := flag.Int("top", 10, "how many top shows to print")
	shouldPrintTopShows := flag.Bool("print", false, "whether to print top shows for each user")
	flag.Parse()

	// Make sure we have an export file.
	if len(*exportFilename) < 1 {
		log.Fatalln("need a discord export file")
	}

	// Open the export file.
	exportFile, err := os.Open(*exportFilename)
	if err != nil {
		log.Fatalf("opening %s: %v", *exportFilename, err)
	}

	// Deserialize the export.
	log.Printf("starting to deserialize: %s", path.Base(*exportFilename))
	decoder := json.NewDecoder(exportFile)
	export := discord.DiscordExport{}
	if err := decoder.Decode(&export); err != nil {
		log.Fatalf("deserealizing json: %v", err)
	}

	// Log some stats.
	log.Printf("found %s messages to comb through", humanize.Comma(int64(len(export.Messages))))

	// Open the output file.
	csvWriter, file, err := getCsvWriter(*outputFilename)
	if err != nil {
		log.Fatalf("getting csv writer: %v", err)
	}

	// Write the header.
	csvWriter.Write([]string{
		"user",
		"claim_rank",
		"likes_rank",
		"kakera_price",
		"show_title",
		"wife",
		"weekday",
		"hour",
	})

	// Load the location.
	loc, err := time.LoadLocation(Location)
	if err != nil {
		log.Fatalf("loading location %s: %v", Location, err)
	}

	// Process the messages.
	log.Println("starting to process...")

	// Keep track of some stats.
	mudae, hadClaim, hadLikes, hadKakera := int64(0), int64(0), int64(0), int64(0)

	// Keep track of the shows.
	showsFrequencies := map[string]map[string]int{}

	// Loop through the messages.
	for _, message := range export.Messages {
		// Only look for Mudae's responses.
		if message.Author.ID != MudaeID {
			continue
		}
		// Only look for "/wa" interactions.
		if message.Interaction.Name != MudaeGachaInteraction {
			continue
		}
		// Need embeds.
		if len(message.Embeds) < 1 {
			continue
		}

		// Increment the number of mudae messages.
		mudae++

		// Extract the properties.
		embed := message.Embeds[0]
		description := embed.Description
		user := message.Interaction.User.Name
		claimRank := extractProperty(claimRankRegexp, description, &hadClaim)
		likesRank := extractProperty(likesRankRegexp, description, &hadLikes)
		kakeraPrice := extractProperty(kakeraPriceRegexp, description, &hadKakera)

		// Extract the show title.
		showTitle := extractProperty(showTitleRegexp, description, nil)
		// Keep track of the shows.
		if _, ok := showsFrequencies[user]; !ok {
			showsFrequencies[user] = map[string]int{}
		}
		// Increment the show frequency for the user.
		if len(strings.TrimSpace(showTitle)) > 0 {
			showsFrequencies[user][showTitle]++
		}

		// Convert the timestamp to central time.
		centralTime := message.Timestamp.In(loc)

		// Write the row.
		csvWriter.Write([]string{
			user,
			claimRank,
			likesRank,
			kakeraPrice,
			showTitle,
			embed.Author.Name,
			centralTime.Weekday().String(),
			strconv.Itoa(centralTime.Hour()),
		})

	}

	// Log some stats.
	log.Printf("processed %s mudae messages", humanize.Comma(mudae))
	log.Printf("only %s had claim rank information", humanize.Comma(hadClaim))
	log.Printf("only %s had likes rank information", humanize.Comma(hadLikes))
	log.Printf("only %s had kakera price information", humanize.Comma(hadKakera))

	// Flush the writer and close the file.
	csvWriter.Flush()
	file.Close()

	log.Printf("wrote output to %s", *outputFilename)

	// Print the top shows.
	if *shouldPrintTopShows {
		printTopShowsFunc(showsFrequencies, *printTopShows)
	}
}

func printTopShowsFunc(showsFrequencies map[string]map[string]int, printTopShows int) {
	fmt.Printf("\nTop %d shows by user:\n", printTopShows)
	// Print the shows.
	fmt.Println("------------")
	for user, shows := range showsFrequencies {
		showsArray := make([]string, 0, len(showsFrequencies))
		for key := range shows {
			showsArray = append(showsArray, key)
		}
		sort.SliceStable(showsArray, func(i, j int) bool {
			return shows[showsArray[i]] > shows[showsArray[j]]
		})

		fmt.Println("USER:", user)
		for i := 0; i < printTopShows; i++ {
			//fmt.Printf(showsArray[i], shows[showsArray[i]])
			fmt.Printf("%d. %s (%d)\n", i+1, showsArray[i], shows[showsArray[i]])
		}
		fmt.Println("------------")
	}
}

// getCsvWriter returns a csv writer and the file it's writing to.
func getCsvWriter(filename string) (*csv.Writer, *os.File, error) {
	output, err := os.Create(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("creating file: %v", err)
	}
	csvWriter := csv.NewWriter(output)
	return csvWriter, output, nil
}

// extractProperty extracts the first group from the given pattern in the given
// string. If the group is found, the counter is incremented.
func extractProperty(pattern *regexp.Regexp, what string, counter *int64) string {
	extracted := extractFirstGroup(pattern, what)
	if counter != nil && extracted != "" {
		(*counter)++
	}
	return extracted
}

// extractFirstGroup extracts the first group from the given pattern in the
// given string.
func extractFirstGroup(pattern *regexp.Regexp, what string) string {
	matches := pattern.FindStringSubmatch(what)
	if matches == nil || len(matches) < 2 {
		return ""
	}
	return matches[1]
}
