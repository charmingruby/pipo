package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/charmingruby/pipo/config"
	"github.com/charmingruby/pipo/pkg/logger"
	"github.com/charmingruby/pipo/pkg/redis"
	"github.com/joho/godotenv"
)

const (
	MAX_RECORDS = 241145
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

	_, err = redis.New(cfg.RedisURL)
	if err != nil {
		logger.Error("failed to connect to redis", "error", err)

		os.Exit(1)
	}

	logger.Info("redis connected")

	records, err := readCsvFile("./dataset/sentiment_data.csv", args.Records)
	if err != nil {
		logger.Error("failed to read csv file", "error", err)

		os.Exit(1)
	}

	fmt.Println(records)
}

type Args struct {
	FilePath string
	Records  int
}

func parseArgs() (Args, error) {
	filePath := flag.String("file", "./dataset/sentiment_data.csv", "path to the csv file, default is ./dataset/sentiment_data.csv")
	records := flag.Int("records", MAX_RECORDS, "number of records to read, max value is "+strconv.Itoa(MAX_RECORDS))

	flag.Parse()

	if *records > MAX_RECORDS {
		return Args{}, errors.New("records must be less than or equal to " + strconv.Itoa(MAX_RECORDS))
	}

	return Args{
		FilePath: *filePath,
		Records:  *records,
	}, nil
}

func readCsvFile(filePath string, amountOfRecords int) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if amountOfRecords > 0 {
		records = records[1 : amountOfRecords+1]
	}

	return records, nil
}
