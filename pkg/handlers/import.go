package handlers

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/tsagita/go-vue-postgre/pkg/database"
)

var files = map[string]string{
	"order_items":        "Test task - Postgres - order_items.csv",
	"customer_companies": "Test task - Postgres - customer_companies.csv",
	"customers":          "Test task - Postgres - customers.csv",
	"deliveries":         "Test task - Postgres - deliveries.csv",
	"orders":             "Test task - Postgres - orders.csv",
}

var testDataDir = "./test_data"

func Import(c *gin.Context) {
	var successRecords []map[string]interface{}
	for tableName, file := range files {
		filePath := filepath.Join(testDataDir, file)
		result, err := readCSV(filePath)
		if err != nil {
			fmt.Printf("Error reading %s: %s\n", file, err)
			continue
		}

		err = insert(tableName, result)
		if err != nil {
			fmt.Printf("Error inserting into PostgreSQL for %s: %s\n", tableName, err)
			continue
		}

		successRecord := make(map[string]interface{})
		successRecord["table"] = tableName
		successRecord["records"] = generateSuccessRecords(result)
		successRecords = append(successRecords, successRecord)

	}
	c.JSON(200, successRecords)
}

func generateSuccessRecords(data []map[string]string) string {
	var successRecords []string
	for _, row := range data {
		successRecords = append(successRecords, generateRecordString(row))
	}
	return "Success Insert: " + strings.Join(successRecords, " - ")
}

func generateRecordString(row map[string]string) string {
	var values []string
	for _, value := range row {
		values = append(values, value)
	}
	return strings.Join(values, " - ")
}

func readCSV(filePath string) ([]map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error opening file %s: %s", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // Set the delimiter to tab

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error reading CSV %s: %s", filePath, err)
	}

	headers := records[0]
	fmt.Println("Headers:", headers)

	var result []map[string]string

	for _, row := range records[1:] {
		data := make(map[string]string)
		for i := 0; i < len(row); i++ {
			data[headers[i]] = row[i]
		}
		result = append(result, data)
	}

	return result, nil
}

func insert(tableName string, data []map[string]string) error {
	var err error
	db, err = database.Connection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the table if it doesn't exist
	_, err = db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s
		)
	`, tableName, generateTableColumns(data[0])))
	if err != nil {
		return err
	}

	// Insert data into the table
	for _, row := range data {
		columns, values := generateInsertQuery(row)
		fmt.Println("columns:", columns)
		fmt.Println("values:", values)

		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columns, values)
		fmt.Println("query:", query)
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateInsertQuery(row map[string]string) (string, string) {
	var columns, values string
	for key, value := range row {
		columns += key + ","
		if value == "" {
			values += "NULL,"
		} else {
			values += fmt.Sprintf("'%s',", value)
		}
	}
	return columns[:len(columns)-1], values[:len(values)-1]
}

func generateTableColumns(row map[string]string) string {
	var columns string
	for key := range row {
		columns += fmt.Sprintf("%s VARCHAR(255),", key)
	}
	return columns[:len(columns)-1]
}
