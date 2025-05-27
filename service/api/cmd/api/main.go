package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"strconv"

	"github.com/charmingruby/pipo/lib/broker/redis"
	"github.com/charmingruby/pipo/lib/logger"
	"github.com/charmingruby/pipo/service/api/config"
	"github.com/charmingruby/pipo/service/api/internal/core/service"
	"github.com/joho/godotenv"
)

func main() {
	logger := logger.New()

	args, err := parseArgs()
	if err != nil {
		logger.Error("failed to parse args", "error", err)

		os.Exit(1)
	}

	logger.Info("args parsed", "args", args)

	if err := godotenv.Load(); err != nil {
		logger.Warn("failed to find .env file", "error", err)
	}

	cfg, err := config.New()
	if err != nil {
		logger.Error("failed to load config", "error", err)

		os.Exit(1)
	}

	logger.Info("config loaded")

	redisClient, err := redis.NewClient(cfg.RedisURL)
	if err != nil {
		logger.Error("failed to connect to redis", "error", err)

		os.Exit(1)
	}

	logger.Info("redis connected")

	redisBroker := redis.NewStream(redisClient)

	sentimentService := service.New(logger, redisBroker, cfg.SentimentIngestedTopic)

	if _, err := sentimentService.IngestRawData(
		context.Background(),
		service.IngestRawDataInput{
			FilePath: args.FilePath,
			Records:  args.Records,
		},
	); err != nil {
		logger.Error("failed to process raw sentiment data", "error", err)

		os.Exit(1)
	}
}

type Args struct {
	FilePath string
	Records  int
}

func parseArgs() (Args, error) {
	maxRecords := 241145

	filePath := flag.String("file", "../../data/sentiment_data.csv", "path to the csv file, default is ./data/sentiment_data.csv")
	records := flag.Int("records", maxRecords, "number of records to read, max value is "+strconv.Itoa(maxRecords))

	flag.Parse()

	if *records > maxRecords {
		return Args{}, errors.New("records must be less than or equal to " + strconv.Itoa(maxRecords))
	}

	return Args{
		FilePath: *filePath,
		Records:  *records,
	}, nil
}
