package csv

import (
	"encoding/csv"
	"io"
)

// ParseFile reads a CSV file and returns the records.
//
// reader is the io.Reader containing the CSV data.
// amountOfRecords is the number of records to be read.
//
// Returns the records and an error if one occurs.
func ParseFile(reader io.Reader, amountOfRecords int) ([][]string, error) {
	csvReader := csv.NewReader(reader)

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if amountOfRecords > 0 {
		records = records[1 : amountOfRecords+1] // skip header
	}

	return records, nil
}
