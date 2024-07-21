package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type AIModelConnector struct {
	Client *http.Client
}

type Inputs struct {
	Table map[string][]string `json:"table"`
	Query string              `json:"query"`
}

type Response struct {
	Answer      string   `json:"answer"`
	Coordinates [][]int  `json:"coordinates"`
	Cells       []string `json:"cells"`
	Aggregator  string   `json:"aggregator"`
}

func CsvToSlice(data string) (map[string][]string, error) {
	reader := csv.NewReader(strings.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	headers := records[0]
	result := make(map[string][]string)
	for _, header := range headers {
		result[header] = []string{}
	}

	for _, record := range records[1:] {
		for i, value := range record {
			result[headers[i]] = append(result[headers[i]], value)
		}
	}

	return result, nil // TODO: replace this
}

func (c *AIModelConnector) ConnectAIModel(payload interface{}, token string) (Response, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Response{}, err
	}

	req, err := http.NewRequest("POST", "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", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return Response{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return Response{}, err
	}

	return response, nil // TODO: replace this
}

func main() {
	token := ""

	// Membuat instance dari AIModelConnector
	connector := AIModelConnector{
		Client: &http.Client{},
	}

	// Membaca data dari file CSV
	csvData, err := os.ReadFile("data-series.csv")
	if err != nil {
		log.Fatalf("Error reading CSV file: %v", err)
	}

	// Mengkonversi data CSV menjadi slice
	dataSlice, err := CsvToSlice(string(csvData))
	if err != nil {
		log.Fatalf("Error converting CSV to slice: %v", err)
	}

	// Contoh query untuk model AI
	query := "What is the total energy consumption for the refrigerator?"

	// Membuat payload untuk model AI
	payload := Inputs{
		Table: dataSlice,
		Query: query,
	}

	// Menghubungkan ke model AI dan mendapatkan respons
	response, err := connector.ConnectAIModel(payload, token)
	if err != nil {
		log.Fatalf("Error connecting to AI model: %v", err)
	}

	// Mencetak respons dari model AI
	fmt.Printf("AI Model Response: %+v\n", response) // TODO: answer here
}
