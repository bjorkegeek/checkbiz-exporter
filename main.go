package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var checkBizToken string

// readSingleLineFile reads the first line from a file
func readSingleLineFile(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// getCheckbizToken retrieves the API token from an environment variable or a file
func getCheckbizToken() (string, error) {
	tokenFile, fileExists := os.LookupEnv("CHECKBIZ_TOKEN_FILE")
	if fileExists {
		token, err := readSingleLineFile(tokenFile)
		if err != nil {
			return "", err
		}
		return token, nil
	}

	token, tokenExists := os.LookupEnv("CHECKBIZ_TOKEN")
	if tokenExists {
		return token, nil
	}

	return "", fmt.Errorf("missing CHECKBIZ_TOKEN_FILE or CHECKBIZ_TOKEN environment variable")
}

func fetchAPIData() (map[string]interface{}, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.checkbiz.se/api/v1/packagecalls", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Basic "+checkBizToken)
	req.Header.Add("User-Agent", "CheckBiz Prometheus Exporter")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Function to convert API response to Prometheus format
func convertToPrometheusMetrics(apiResponse string) string {
	// Implement conversion logic here
	// For now, returning dummy metrics
	return "api_calls_total 123\n"
}

// printMetrics prints the Prometheus metrics
func printMetrics(w http.ResponseWriter, data map[string]interface{}) {
	wroteMeta := false

	products, ok := data["products"].([]interface{})
	if !ok {
		fmt.Fprintln(w, "Error: Invalid data format for products")
		return
	}

	totals := make(map[string]int)

	for _, p := range products {
		product, ok := p.(map[string]interface{})
		if !ok {
			continue
		}

		packages, ok := product["packages"].([]interface{})
		if !ok {
			continue
		}

		for _, pk := range packages {
			pkg, ok := pk.(map[string]interface{})
			if !ok {
				continue
			}

			for key, value := range pkg {
				if strings.HasPrefix(key, "numberOfCalls") {
					if !wroteMeta {
						fmt.Fprintln(w, `HELP checkbiz_call_count The number of API calls made
TYPE checkbiz_call_count counter`)
						wroteMeta = true
					}

					labels := map[string]string{
						"product": fmt.Sprintf("%v", product["productName"]),
						"package": fmt.Sprintf("%v", pkg["packageName"]),
						"period":  key[13:],
					}
					var labelParts []string
					for k, v := range labels {
						labelParts = append(labelParts, fmt.Sprintf("%s=%q", k, v))
					}
					labelString := "{" + strings.Join(labelParts, ",") + "}"
					fmt.Fprintf(w, "checkbiz_call_count%s %v\n", labelString, value)
					if valueFloat, ok := value.(float64); ok {
						totals[labels["period"]] += int(valueFloat)
					}
				}
			}
		}
	}
	if len(totals) > 0 {
		fmt.Fprintln(w, `
HELP checkbiz_call_count The total number of API calls made for each period
TYPE checkbiz_call_count counter
`)
		for period, total := range totals {
			fmt.Fprintf(w, "checkbiz_call_count_total{period=%q} %d\n", period, total)
		}
	}
}

// Handler for your web server
func metricsHandler(w http.ResponseWriter, r *http.Request) {
	apiResponse, err := fetchAPIData()
	if err != nil {
		http.Error(w, "Error fetching API data", http.StatusInternalServerError)
		return
	}
	printMetrics(w, apiResponse)
}
func init() {
	var err error
	checkBizToken, err = getCheckbizToken()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error retrieving API token:", err)
		os.Exit(2)
	}
}
func main() {
	http.HandleFunc("/metrics", metricsHandler)
	fmt.Println("Server is starting on port 8080...")
	http.ListenAndServe(":8080", nil)
}
