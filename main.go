package main

import (
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/7Maliko7/telegram-bot/internal/config"
	"github.com/7Maliko7/telegram-bot/internal/scenario"
	"github.com/7Maliko7/telegram-bot/internal/script"
	"github.com/7Maliko7/telegram-bot/internal/service"
	redisCache "github.com/7Maliko7/telegram-bot/pkg/cache/redis"
	"github.com/7Maliko7/telegram-bot/pkg/db/driver/postgres"
	"github.com/7Maliko7/telegram-bot/pkg/messenger/telegram"
)

const (
	appVersion = "0.1.2.3"
)

var (
	configPath, scenarioPath string
)

func init() {
	flag.StringVar(&configPath, "c", "", "Custom config path")
	flag.StringVar(&scenarioPath, "s", "", "Custom scenario path")
	flag.Parse()
}

func main() {
	log.Info().Msgf("Scenario Bot %v", appVersion)
	log.Info().Msg("app configuration...")
	appConfig, err := config.New(configPath)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg("app configurated")

	log.Info().Msg("tg connection...")
	tgBot, err := telegram.New(appConfig.Telegram.Token)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg("tg connected")

	log.Info().Msg("bot scenario...")
	scr, err := script.New(scenarioPath)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg("bot scenario loaded")

	scnr := scenario.New(scr)

	log.Info().Msg("cache connection...")
	cache, err := redisCache.New(appConfig.Cache.ConnectionString)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg("cache connected")

	log.Info().Msg("db connection...")
	db, err := postgres.New(appConfig.DB.ConnectionString)
	log.Info().Msg("db connected")

	log.Info().Msg("app starting")
	svc := service.New(&log.Logger, tgBot, cache, db, scnr)
	log.Info().Msg("app successfully started")

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2112", nil)

	log.Info().Msg("listening tg update channel...")
	svc.Run()
}
