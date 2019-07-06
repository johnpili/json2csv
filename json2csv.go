// json2csv
// Author: John Pili
// Website: www.johnpili.com

package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var rules []Rule // The global variable for rules

// parameterPair - Used in URL parameters
type parameterPair struct {
	k string
	v string
}

type ioMappingPair struct {
	Incoming string `json:"incoming"`
	Outgoing string `json:"outgoing"`
}

// Rule - a generic way of defining a rule or endpoint
type Rule struct {
	RuleName       string          `json:"ruleName"`       // Rule name
	IOMapping      []ioMappingPair `json:"ioMapping"`      // This is the input and output result mapping
	TargetURLToken string          `json:"targetUrlToken"` // Authorization key for target api, leave blank if not required
}

// loadRules this method is responsible for loading the rules.json
func loadRules(location string) []Rule {
	raw, err := ioutil.ReadFile(location) // Read the endpoint.json
	if err != nil {
		fmt.Println(err.Error()) // Print the error in reading the rule
		os.Exit(1)               // Exit the program with error flag
	}

	var e []Rule
	json.Unmarshal(raw, &e) // We are expecting an json array here
	return e
}

func main() {
	args := os.Args[1:] // Fetch the program arguments excluding the program name
	if len(args) == 0 {
		fmt.Println("Usage: json2csv rules.json \"rule-name\" \"url\"")
	} else {
		rules = loadRules(args[0]) // Load the rules json location
		ruleName := args[1]        // Set the rule name
		urlPayload := args[2]      // This is the URL where the JSON file will be fetch from

		rule, err := getRule(ruleName) // Fetch the rule from the rule list
		if err != nil {
			fmt.Println(err.Error()) // Print the error
			os.Exit(1)               // Exit program with error flag
		}

		csvWriter := csv.NewWriter(os.Stdout) // Write the CSV to the standard output
		csvWriter.UseCRLF = true
		csvWriter.WriteAll(galacticCSV(restfulClientExchanger(rule, urlPayload), rule))
	}
}

func getRule(ruleName string) (Rule, error) {
	for _, item := range rules { // Iterate through the rule list
		if item.RuleName == ruleName {
			return item, nil // Return the found rule and nil error
		}
	}
	var rule Rule // Return an empty rule instead of nil
	return rule, errors.New("No rule found")
}

// printHeader - This will print the header
func printHeader(ioMapping []string) []string {
	headers := make([]string, 0)
	for _, targetKey := range ioMapping { // Iterate IOMapping
		headers = append(headers, fmt.Sprint(targetKey))
	}
	return headers
}

// galacticColumnizer - This will columnized the value of CSV
func galacticColumnizer(jsonNode map[string]interface{}, ioMapping []string) []string {
	columns := make([]string, 0)
	for _, targetKey := range ioMapping { // Iterate IOMapping
		for nodeKey, nodeValue := range jsonNode {
			if nodeKey == targetKey {
				if nodeValue == nil { // if nodeValue is nil then just give blank string
					columns = append(columns, "")
				} else {
					columns = append(columns, fmt.Sprint(nodeValue))
				}
			}
		}
	}
	return columns
}

// galacticCSV - This will convert JSON string into a CSV
func galacticCSV(body []byte, rule Rule) [][]string {
	var jsonArray []interface{} // A generic way of handling json without struct
	json.Unmarshal([]byte(body), &jsonArray)
	var overall = make([][]string, 0) // Make a zero sized 2 dimensional array
	ioMapping := make([]string, 0)
	headers := make([]string, 0)

	isHeaderPrinted := false
	isIOMappingGenerated := false
	for _, jsonNode := range jsonArray {
		row := make([]string, 0)

		if isIOMappingGenerated == false {
			if len(ioMapping) <= 0 && len(rule.IOMapping) <= 0 {
				for nodeKey, _ := range jsonNode.(map[string]interface{}) {
					ioMapping = append(ioMapping, nodeKey) // Generate a IO Mapping based on the received json file. Take note that golang don't preserve the JSON column order
					headers = append(headers, nodeKey)
				}
			} else if len(rule.IOMapping) > 0 { // If IOMapping has been populated
				for _, nodeValue := range rule.IOMapping {
					ioMapping = append(ioMapping, nodeValue.Incoming) // Copy the IO Mapping from Endpoint configuration
					headers = append(headers, nodeValue.Outgoing)
				}
			}
			isIOMappingGenerated = true
		}

		if isHeaderPrinted == false {
			overall = append(overall, printHeader(headers))
			isHeaderPrinted = true
		}

		row = append(row, galacticColumnizer(jsonNode.(map[string]interface{}), ioMapping)...)
		overall = append(overall, row)
	}
	return overall
}

// galaticURLEncoder - This handles the URL parameters and encoding
func galaticURLEncoder(targetURL string) string {
	var x *url.URL
	x, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal(err)
	}
	return x.String()
}

// restfulClientExchanger - This method handles the HTTP communication to DBFlare, IOREST or other REST services
func restfulClientExchanger(r Rule, u string) []byte {
	urlPayload := galaticURLEncoder(u)
	client := &http.Client{
		Timeout: time.Second * 180,
	}
	req, _ := http.NewRequest("GET", urlPayload, nil)
	if len(r.TargetURLToken) > 0 { // If TargetURLToken is blank or empty then don't create Authorization header
		req.Header.Set("Authorization", r.TargetURLToken) // Set JWT based Authorization token
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}
