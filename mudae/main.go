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

	"github.com/dustin/go-humanize"
)

const (
	MudaeID               = "432610292342587392"
	MudaeGachaInteraction = "wa"
	Location              = "America/Chicago"
)

var (
	claimRankRegexp   = regexp.MustCompile(`Claims: #([0-9]+)`)
	likesRankRegexp   = regexp.MustCompile(`Likes: #([0-9]+)`)
	kakeraPriceRegexp = regexp.MustCompile(`\*\*([0-9]+)\*\*:kakera`)
	showTitleRegexp   = regexp.MustCompile(`([^/]+)\n`)
)

var (
	showsFrequencies = map[string]map[string]int{}
)

func main() {
	exportFilename := flag.String("export", "", "discord export file")
	outputFilename := flag.String("output", "output.csv", "where to output")
	flag.Parse()

	if len(*exportFilename) < 1 {
		log.Fatalln("need a discord export file")
	}

	exportFile, err := os.Open(*exportFilename)
	if err != nil {
		log.Fatalf("opening %s: %v", *exportFilename, err)
	}

	log.Printf("starting to deserialize: %s", path.Base(*exportFilename))
	decoder := json.NewDecoder(exportFile)
	export := DiscordExport{}
	if err := decoder.Decode(&export); err != nil {
		log.Fatalf("deserealizing json: %v", err)
	}

	log.Printf("found %s messages to comb through", humanize.Comma(int64(len(export.Messages))))

	output, err := os.Create(*outputFilename)
	if err != nil {
		log.Fatalf("couldn't create output: %v", err)
	}
	csvWriter := csv.NewWriter(output)

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

	loc, err := time.LoadLocation(Location)
	if err != nil {
		log.Fatalf("loading location %s: %v", Location, err)
	}

	log.Println("starting to process...")
	mudae, hadClaim, hadLikes, hadKakera := int64(0), int64(0), int64(0), int64(0)
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

		mudae++

		embed := message.Embeds[0]
		description := embed.Description
		user := message.Interaction.User.Name
		claimRank := extractProperty(claimRankRegexp, description, &hadClaim)
		likesRank := extractProperty(likesRankRegexp, description, &hadLikes)
		kakeraPrice := extractProperty(kakeraPriceRegexp, description, &hadKakera)

		showTitle := extractProperty(showTitleRegexp, description, nil)
		if _, ok := showsFrequencies[user]; !ok {
			showsFrequencies[user] = map[string]int{}
		}
		if len(strings.TrimSpace(showTitle)) > 0 {
			showsFrequencies[user][showTitle]++
		}

		centralTime := message.Timestamp.In(loc)
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
	log.Printf("processed %s mudae messages", humanize.Comma(mudae))
	log.Printf("only %s had claim rank information", humanize.Comma(hadClaim))
	log.Printf("only %s had likes rank information", humanize.Comma(hadLikes))
	log.Printf("only %s had kakera price information", humanize.Comma(hadKakera))
	csvWriter.Flush()
	output.Close()

	for user, shows := range showsFrequencies {
		showsArray := make([]string, 0, len(showsFrequencies))
		for key := range shows {
			showsArray = append(showsArray, key)
		}
		sort.SliceStable(showsArray, func(i, j int) bool {
			return shows[showsArray[i]] > shows[showsArray[j]]
		})

		fmt.Println(user)
		for i := 0; i < 10; i++ {
			fmt.Println(showsArray[i], shows[showsArray[i]])
		}
		fmt.Println("------------")
	}

}

func extractProperty(pattern *regexp.Regexp, what string, counter *int64) string {
	extracted := extractFirstGroup(pattern, what)
	if counter != nil && extracted != "" {
		(*counter)++
	}
	return extracted
}

func extractFirstGroup(pattern *regexp.Regexp, what string) string {
	matches := pattern.FindStringSubmatch(what)
	if matches == nil || len(matches) < 2 {
		return ""
	}
	return matches[1]
}
