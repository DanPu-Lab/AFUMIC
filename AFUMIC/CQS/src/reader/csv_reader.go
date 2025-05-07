package reader

import (
	"encoding/csv"
	"os"
)

type CSVReader struct {
	reader *csv.Reader
}

func NewCSVReader() *CSVReader {
	return &CSVReader{}
}

func (r *CSVReader) Read(filePath string, seq rune) error {
	csvFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	reader := csv.NewReader(csvFile)
	reader.Comma = seq
	r.reader = reader
	return nil
}

func (r *CSVReader) ReadAll() ([][]string, error) {
	return r.reader.ReadAll()
}

func (r *CSVReader) NextLine() ([]string, error) {
	return r.reader.Read()
}
