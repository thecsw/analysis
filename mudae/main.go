package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
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
	kakeraPriceRegexp = regexp.MustCompile(`\*\*([0-9]+)\*\*`)
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
		"kakera_price",
		"weekday",
		"hour",
	})

	loc, err := time.LoadLocation(Location)
	if err != nil {
		log.Fatalf("loading location %s: %v", Location, err)
	}

	log.Println("starting to process...")
	mudae, eligible := int64(0), int64(0)
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

		claimRank := extractFirstGroup(
			claimRankRegexp,
			message.Embeds[0].Description,
		)

		kakeraPrice := extractFirstGroup(
			kakeraPriceRegexp,
			message.Embeds[0].Description,
		)

		if claimRank == "" || kakeraPrice == "" {
			continue
		}

		centralTime := message.Timestamp.In(loc)

		eligible++
		csvWriter.Write([]string{
			message.Interaction.User.Name,
			claimRank,
			kakeraPrice,
			centralTime.Weekday().String(),
			strconv.Itoa(centralTime.Hour()),
		})

	}
	log.Printf("processed %s mudae messages", humanize.Comma(mudae))
	log.Printf("only %s were eligible for export", humanize.Comma(eligible))
	csvWriter.Flush()
	output.Close()

}

func extractFirstGroup(pattern *regexp.Regexp, what string) string {
	matches := pattern.FindStringSubmatch(what)
	if matches == nil || len(matches) < 2 {
		return ""
	}
	return matches[1]
}
