package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BlueMonday/go-scryfall"
	log "github.com/sirupsen/logrus"
)

const DefaultConfigFile = "Config.yaml"
const DefaultSpoiledCardsFile = "SpoiledCards.json"
const DefaultLogFile = "ScryfallSpoilerBot.log"

func main() {
	ConfigFile := flag.String("config", DefaultConfigFile, "path to yaml config file")
	SpoiledCardsFile := flag.String("spoiledCards", DefaultSpoiledCardsFile, "path to json spoiled cards file")
	LogFile := flag.String("logFile", DefaultLogFile, "path to log file")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	f, err := os.OpenFile(*LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)

	spoiledCards, err := ReadSpoiledCards(*SpoiledCardsFile)
	if err != nil {
		log.Error(err)
	}

	conf, err := NewConfig(*ConfigFile)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	client, err := scryfall.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	sco := scryfall.SearchCardsOptions{
		Unique: scryfall.UniqueModePrints,
		Order:  "spoiled",
		Dir:    scryfall.DirDesc,
	}

	for {
		sets, err := client.ListSets(ctx)
		if err != nil {
			log.Error(err)
		}

		for _, set := range sets {
			if contains(conf.SetBlacklist, set.Code) ||
				!contains(conf.SetTypes, string(set.SetType)) ||
				(set.ReleasedAt != nil && time.Now().Sub(set.ReleasedAt.Time) > 0) {

				delete(*spoiledCards, set.Code)
				continue
			}
			log.Infof("Checking set %s [release: %s]", set.Name, set.ReleasedAt)

			result, err := client.SearchCards(ctx, fmt.Sprintf("e:%s", set.Code), sco)
			if err != nil {
				log.Error(err)
			}

			for _, card := range result.Cards {
				if contains((*spoiledCards)[set.Code], card.CollectorNumber) {
					continue
				}

				log.Infof("New Card found %s", card.Name)

				for messenger := range conf.Messengers.Get() {
					err = messenger.SendCard(card)
					if err != nil {
						log.Error(err)
					}
				}
				(*spoiledCards)[set.Code] = append((*spoiledCards)[set.Code], card.CollectorNumber)
			}
		}

		select {
		case <-time.After(time.Duration(conf.CheckInterval) * time.Minute):
		case <-sigs:
			log.Info("Terminating")
			SaveSpoiledCards(*SpoiledCardsFile, spoiledCards)
			return
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
