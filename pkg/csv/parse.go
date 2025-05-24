package csv

import (
	"encoding/csv"
	"os"
)

func ReadFile(filePath string, amountOfRecords int) ([][]string, error) {
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
		records = records[1 : amountOfRecords+1] // skip header
	}

	return records, nil
}
