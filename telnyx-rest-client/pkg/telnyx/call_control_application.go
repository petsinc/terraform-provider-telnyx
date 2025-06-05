package telnyx

import (
	"fmt"

	"go.uber.org/zap"
)

// CreateCallControlApplication creates a new Call Control Application.
func (client *TelnyxClient) CreateCallControlApplication(request CallControlApplicationRequest) (*CallControlApplication, error) {
	var result struct {
		Data CallControlApplication `json:"data"`
	}
	err := client.doRequest("POST", "/call_control_applications", request, &result)
	if err != nil {
		client.logger.Error("Error creating Call Control Application", zap.Error(err))
		return nil, err
	}
	return &result.Data, nil
}

// GetCallControlApplication retrieves a Call Control Application by ID.
func (client *TelnyxClient) GetCallControlApplication(applicationID string) (*CallControlApplication, error) {
	var result struct {
		Data CallControlApplication `json:"data"`
	}
	err := client.doRequest("GET", fmt.Sprintf("/call_control_applications/%s", applicationID), nil, &result)
	if err != nil {
		client.logger.Error("Error fetching Call Control Application", zap.Error(err), zap.String("applicationID", applicationID))
		return nil, err
	}
	return &result.Data, nil
}

// UpdateCallControlApplication updates an existing Call Control Application.
func (client *TelnyxClient) UpdateCallControlApplication(applicationID string, request CallControlApplicationRequest) (*CallControlApplication, error) {
	var result struct {
		Data CallControlApplication `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/call_control_applications/%s", applicationID), request, &result)
	if err != nil {
		client.logger.Error("Error updating Call Control Application", zap.Error(err), zap.String("applicationID", applicationID))
		return nil, err
	}
	return &result.Data, nil
}

// DeleteCallControlApplication deletes a Call Control Application.
func (client *TelnyxClient) DeleteCallControlApplication(applicationID string) error {
	err := client.doRequest("DELETE", fmt.Sprintf("/call_control_applications/%s", applicationID), nil, nil)
	if err != nil {
		client.logger.Error("Error deleting Call Control Application", zap.Error(err), zap.String("applicationID", applicationID))
	}
	return err
}

// ListCallControlApplications lists all Call Control Applications.
func (client *TelnyxClient) ListCallControlApplications() ([]CallControlApplication, error) {
	var result struct {
		Data []CallControlApplication `json:"data"`
	}
	err := client.doRequest("GET", "/call_control_applications", nil, &result)
	if err != nil {
		client.logger.Error("Error listing Call Control Applications", zap.Error(err))
		return nil, err
	}
	return result.Data, nil
}
