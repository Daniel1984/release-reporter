package main

import (
	"time"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/release-reporter/api/slack"
	"github.com/release-reporter/gitscan"
	"github.com/release-reporter/gracefulshutdown"
	"github.com/release-reporter/logger"
	"github.com/release-reporter/models"
	"github.com/spf13/viper"
)

func main() {
	// configure viper to consume configuration
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	// create new instance of logger
	log := logger.New()

	// consumer configuration
	if err := viper.ReadInConfig(); err != nil {
		log.Err.Fatalf("Error reading config file, %s", err)
	}

	cfg := &models.Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Err.Fatalf("unable to decode into struct, %v", err)
	}

	// initiate db for persistence
	db, err := badger.Open(badger.DefaultOptions(cfg.DBPath))
	if err != nil {
		log.Err.Fatal(err)
	}

	// start check process in separate go routine so that rest of process could
	// still be handeled by gracefulshutdown service
	go func() {
		for {
			gitReport := gitscan.ScanReleases(cfg, db, log)
			if len(gitReport) > 0 {
				if err := slack.SendMessage(gitReport, cfg.Webhook); err != nil {
					log.Err.Printf("%s\n", err)
				}
			}

			time.Sleep(cfg.CheckIntervalMs * time.Millisecond)
		}
	}()

	gracefulshutdown.Init(log.Err, func() {
		log.Info.Println("performing shutdown")
		if err := db.Close(); err != nil {
			log.Err.Printf("failed closing DB connection:%s", err)
		}
	})
}
