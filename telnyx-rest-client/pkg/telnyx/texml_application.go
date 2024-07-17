package telnyx

import (
	"fmt"
	"go.uber.org/zap"
)

// CreateTeXMLApplication creates a new TeXML application.
func (client *TelnyxClient) CreateTeXMLApplication(request TeXMLApplicationRequest) (*TeXMLApplication, error) {
	var result struct {
		Data TeXMLApplication `json:"data"`
	}
	err := client.doRequest("POST", "/texml_applications", request, &result)
	if err != nil {
		client.logger.Error("Error creating TeXML application", zap.Error(err))
		return nil, err
	}
	return &result.Data, nil
}

// UpdateTeXMLApplication updates an existing TeXML application.
func (client *TelnyxClient) UpdateTeXMLApplication(applicationID string, request TeXMLApplicationRequest) (*TeXMLApplication, error) {
	var result struct {
		Data TeXMLApplication `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/texml_applications/%s", applicationID), request, &result)
	if err != nil {
		client.logger.Error("Error updating TeXML application", zap.Error(err), zap.String("applicationID", applicationID))
		return nil, err
	}
	return &result.Data, nil
}

// DeleteTeXMLApplication deletes a TeXML application.
func (client *TelnyxClient) DeleteTeXMLApplication(applicationID string) error {
	err := client.doRequest("DELETE", fmt.Sprintf("/texml_applications/%s", applicationID), nil, nil)
	if err != nil {
		client.logger.Error("Error deleting TeXML application", zap.Error(err), zap.String("applicationID", applicationID))
	}
	return err
}

// GetTeXMLApplication retrieves a TeXML application by ID.
func (client *TelnyxClient) GetTeXMLApplication(applicationID string) (*TeXMLApplication, error) {
	var result struct {
		Data TeXMLApplication `json:"data"`
	}
	err := client.doRequest("GET", fmt.Sprintf("/texml_applications/%s", applicationID), nil, &result)
	if err != nil {
		client.logger.Error("Error fetching TeXML application", zap.Error(err), zap.String("applicationID", applicationID))
		return nil, err
	}
	return &result.Data, nil
}
