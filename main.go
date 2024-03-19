package main

import (
	"fmt"
	"medLLM/bypass"
	"sync"
)

func main() {
	csvfileinput := "top10.csv"
	entries, err := bypass.ParseCSVDfToStructs(csvfileinput)
	if err != nil {
		fmt.Println("entries are not ready:", err)
		return
	}
	fmt.Println("Entries are ready. Starting the process...")

	var wg sync.WaitGroup // WaitGroup to wait for all goroutines to finish
	var fileMutex sync.Mutex

	// Iterate over the entries slice and start a goroutine for each entry
	for _, e := range entries {

		wg.Add(1) // Increment the WaitGroup counter for each goroutine
		go func(entry bypass.Entry) {
			defer wg.Done()
			fileMutex.Lock()
			defer fileMutex.Unlock() // Signal the WaitGroup that this goroutine is done
			Destruct(entry)
		}(e)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All entries processed")
}

func Destruct(e bypass.Entry) {
	data := bypass.Bypass(bypass.CreatePdfFilefromUrl(e.Link, e.ApplNO), bypass.Codeperneed(e.ApplNO))
	finaldata := bypass.Filter(string(data), "PATIENT COUNSELING INFORMATION")
	finalbyte := []byte(finaldata)

	defer bypass.CreateDeletefile(e.ApplNO, finalbyte)
}
