package csv

import (
	"encoding/csv"
	"log"
	"os"
)

func ReadFile(filePath string, amountOfRecords int) ([][]string, error) {
	f, err := os.Open(filePath) // #nosec G304
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Print("error closing file", "err", err)
		}
	}()

	csvReader := csv.NewReader(f)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if amountOfRecords > 0 {
		records = records[1 : amountOfRecords+1] // skip header
	}

	return records, nil
}
