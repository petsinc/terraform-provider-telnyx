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

func PrettyPrintRequestBody(body map[string]interface{}) error {
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
