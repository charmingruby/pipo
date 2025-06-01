package csv

import (
	"encoding/csv"
	"log"
	"os"
)

// ReadFile reads a CSV file and returns the records.
//
// filePath is the path to the file to be read.
// amountOfRecords is the number of records to be read.
//
// Returns the records and an error if one occurs.
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
