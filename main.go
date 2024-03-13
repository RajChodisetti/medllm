package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"

	"github.com/go-gota/gota/dataframe"
)

type entry struct {
	drug_name string
	ai        string
	route     string
	applNO    string
	company   string
	date      string
	link      string
}

func main() {
	// s, err := parseDataFrameToStructs(newDf("Medication Guides.csv"))
	// if err != nil {
	// 	fmt.Println("Error in main injection", err)
	// }
	text, _ := extractcheckPDF("https://www.accessdata.fda.gov/drugsatfda_docs/label/2023/204311s000lbl.pdf#page=34")
	fmt.Println(text)

}

func newDf(filename string) dataframe.DataFrame {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		panic(err)
	}
	defer file.Close()

	df := dataframe.ReadCSV(file)

	return df
}
func parseDataFrameToStructs(df dataframe.DataFrame) ([]entry, error) {
	var result []entry

	// Get the records from the DataFrame
	records := df.Records()

	// Skip the first row (header row)
	if len(records) > 0 {
		records = records[1:]
	}

	// Iterate over each row in the DataFrame
	for _, record := range records {
		if len(record) != 3 {
			return nil, fmt.Errorf("unexpected number of columns in DataFrame record")
		}

		// Parse each record into MyStruct
		myStruct := entry{
			drug_name: record[0],
			ai:        record[1],
			route:     record[2],
			applNO:    record[3],
			company:   record[4],
			date:      record[5],
			link:      record[6],
		}

		// Append the parsed struct to the result slice
		result = append(result, myStruct)
	}

	return result, nil
}
func createMetaData(s []entry) error {
	for _, e := range s {
		file, err := os.Create("/metadata" + e.applNO + ".txt")
		if err != nil {
			fmt.Println("Error creating file:")
			return err
		}
		defer file.Close()
	}
	return nil

}
func extractcheckPDF(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error downloading PDF file: %v", err)
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading PDF contents: %v", err)
	}

	pdfReader, err := model.NewPdfReader(bytes.NewReader(buf))
	if err != nil {
		return "", fmt.Errorf("error creating PDF reader: %v", err)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", fmt.Errorf("error getting number of pages: %v", err)
	}

	var content strings.Builder
	var isInSection bool

	// Create a new text extractor instance.
	extractor, err := extractor.New(pdfReader)
	if err != nil {
		return "", fmt.Errorf("error creating text extractor: %v", err)
	}

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		// Extract text from the current page.
		pageText, _, _, err := extractor.ExtractPageText(pageNum)
		if err != nil {
			return "", fmt.Errorf("error extracting text from page %d: %v", pageNum, err)
		}

		// Check if the current page contains the section name.
		if strings.Contains(pageText, "PATIENT COUNSELING INFORMATION") {
			// Start recording content if the section is found.
			isInSection = true
		}

		// If in the section, append the content of the current page.
		if isInSection {
			content.WriteString(pageText)
			content.WriteString("\n")
		}
	}

	return content.String(), nil
}
