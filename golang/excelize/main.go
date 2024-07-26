package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	filePath := os.Getenv("EXCEL_FILE_PATH")
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		// Close the spreadsheet
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	sheetName := f.GetSheetName(0)

	// Set week date
	now := time.Now()
	f.SetCellValue(sheetName, "I4", fmt.Sprintf("%s - %s", now.Add(-7*24*time.Hour).Format("2006/01/02"),
		now.Format("2006/01/02")))

	// Report table

	// Save spreadsheet by the given path.
	if err := f.SaveAs("/tmp/slo-report.xlsx"); err != nil {
		log.Println(err)
	}
}
