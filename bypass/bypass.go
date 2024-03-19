package bypass

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

// Payload struct represents the structure of the payload JSON
type Payload struct {
	Code  string `json:"Code"`
	Files []File `json:"Files"`
}

// File struct represents the structure of the file details in the payload JSON
type File struct {
	Filename string `json:"Filename"`
	B64str   string `json:"B64str"`
	MimeType string `json:"MimeType"`
	Size     int    `json:"Size"`
}

// Response struct represents the structure of the JSON response
type Response struct {
	Code     string   `json:"Code"`
	Stdout   string   `json:"Stdout"`
	Stderr   string   `json:"Stderr"`
	Modified []string `json:"Modified"`
}
type Entry struct {
	Drug_name string
	Ai        string
	Route     string
	ApplNO    string
	Company   string
	Date      string
	Link      string
}

func Bypass(pdfFileName string, codeContent []byte) []byte {

	// pdfFileName := "check.pdf"

	// Read the content of the PDF file
	pdfContent, err := ioutil.ReadFile(pdfFileName)
	if err != nil {
		fmt.Println("Error reading PDF file:", err)
		return nil
	}

	// Convert the PDF content to base64
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfContent)

	// Create a new payload
	payload := Payload{
		Code: string(codeContent),
		Files: []File{
			{
				Filename: strings.Replace(pdfFileName, "metadata/", "", -1),
				B64str:   pdfBase64,
				MimeType: "application/pdf",
				Size:     len(pdfContent),
			},
		},
	}

	file, err := os.Create("payload.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	val := reflect.ValueOf(payload)
	typ := reflect.TypeOf(payload)

	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fieldValue := fmt.Sprintf("%v", val.Field(i).Interface())

		line := fmt.Sprintf("%s: %s\n", fieldName, fieldValue)
		_, err := file.WriteString(line)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Convert the payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding payload to JSON:", err)
		return nil
	}

	// Define the URL of unidoc playground
	url := "https://play.unidoc.io/api/run"

	// Create a new HTTP request with the payload data

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadJSON))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return nil
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP client
	client := &http.Client{}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[-] Error sending HTTP request:", err)
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("[-] Error reading response body:", err)
		return nil
	}
	file2, err := os.Create("response.json")
	if err != nil {
		log.Fatal(err)
	}
	_, err = file2.Write([]byte(body))
	if err != nil {
		log.Fatal(err)
	}
	file2.Close()

	// Parse the JSON response
	var response Response
	err = json.Unmarshal(body, &response)

	if err != nil {
		fmt.Println("[-] Error parsing JSON response:", err)

		return nil
	}

	return []byte(response.Stdout)

}

// func NewCSVDf(filename string) dataframe.DataFrame {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		fmt.Println("Error opening file:", err)
// 		panic(err)
// 	}
// 	defer file.Close()

// 	df := dataframe.ReadCSV(file)

//		return df
//	}
func ParseCSVDfToStructs(filename string) ([]Entry, error) {
	var result []Entry
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Skip header row
	_, _ = reader.Read()
	count := 0

	// Iterate over each row in the DataFrame
	for {
		record, err := reader.Read()
		if err != nil {
			fmt.Println("Stopped loop after", count)
			break
		}
		if len(record) != 7 {
			return nil, fmt.Errorf("unexpected number of columns in DataFrame record")
		}

		// Parse each record into MyStruct
		myStruct := Entry{
			Drug_name: record[0],
			Ai:        record[1],
			Route:     record[2],
			ApplNO:    record[3],
			Company:   record[4],
			Date:      record[5],
			Link:      record[6],
		}

		// Append the parsed struct to the result slice
		result = append(result, myStruct)
		count++
	}

	return result, nil
}
func CreatePdfFilefromUrl(url string, pdflocal string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("error downloading PDF file: %v", err)
	}
	defer resp.Body.Close()

	// Read the contents of the PDF file
	pdfContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading PDF contents: %v", err)
	}

	// Specify the path where you want to save the PDF file
	savePath := "metadata/" + pdflocal + ".pdf"

	// Write the PDF content to a file
	err = ioutil.WriteFile(savePath, pdfContent, 0644)
	if err != nil {
		log.Fatalf("error writing PDF file: %v", err)
	}
	return savePath
}

func CreateDeletefile(textlocal string, extracted []byte) error {
	// Specify the directory and filename properly
	txtFilePath := fmt.Sprintf("metadata/%s.txt", textlocal)

	// Create the text file
	txtFile, err := os.Create(txtFilePath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %v", txtFilePath, err)
	}
	defer txtFile.Close()

	// Write the extracted data to the text file
	_, err = txtFile.Write(extracted)
	if err != nil {
		return fmt.Errorf("error writing to file %s: %v", txtFilePath, err)
	}

	// Remove the corresponding PDF file
	pdfPath := fmt.Sprintf("metadata/%s.pdf", textlocal)
	if _, err := os.Stat(pdfPath); err == nil {
		if err := os.Remove(pdfPath); err != nil {
			return fmt.Errorf("error removing PDF file %s: %v", pdfPath, err)
		}
	}

	return nil
}

func Codeperneed(pdffilename string) []byte {
	codeFileName := "code.old.txt"
	// Read the content of the local file
	codeContent, err := ioutil.ReadFile(codeFileName)
	if err != nil {
		fmt.Println("Error reading code file:", err)
		return nil
	}
	pdffile := pdffilename + ".pdf"
	text := string(codeContent)
	// Find and replace text
	text = strings.ReplaceAll(text, "place_holder.pdf", pdffile)
	return []byte(text)
}
func Filter(input string, selector string) string {
	lastIndex := strings.LastIndex(input, selector)
	if lastIndex == -1 {
		// Custom string not found, return empty string or handle error as needed
		return ""
	}

	// Extract text starting from the last instance of the custom string
	extractedText := input[lastIndex+len(selector):]
	return extractedText

}
