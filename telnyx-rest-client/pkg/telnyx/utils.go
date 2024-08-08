package telnyx

import (
	"encoding/json"
	"fmt"
)

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func PrettyPrintRequestBody(body interface{}) error {
	// Pretty print the body as JSON
	prettyJSON, err := json.MarshalIndent(body, "", "  ")
	if err != nil {
		fmt.Println("Failed to generate JSON:", err)
		return err
	}
	fmt.Println("Request Body JSON:")
	fmt.Println(string(prettyJSON))
	return nil
}

// PrettyPrintResponseBody attempts to pretty print the response body if it's JSON.
func PrettyPrintResponseBody(body []byte) error {
	var jsonBody interface{}
	var prettyJSON string

	// Attempt to unmarshal the body to check if it's valid JSON
	if err := json.Unmarshal(body, &jsonBody); err == nil {
		// If valid JSON, pretty print it
		prettyBytes, err := json.MarshalIndent(jsonBody, "", "  ")
		if err != nil {
			fmt.Println("Failed to generate JSON:", err)
			return err
		}
		prettyJSON = string(prettyBytes)
	} else {
		// If not valid JSON, just print the raw body
		prettyJSON = string(body)
	}

	fmt.Println("Response Body JSON:")
	fmt.Println(prettyJSON)
	return nil
}
