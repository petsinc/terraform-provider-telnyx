package telnyx

import (
	"fmt"

	"go.uber.org/zap"
)

// Create Credential Connection
func (client *TelnyxClient) CreateCredentialConnection(profile FQDNConnection) (*FQDNConnection, error) {
	var result struct {
		Data FQDNConnection `json:"data"`
	}
	err := client.doRequest("POST", "/credential_connections", profile, &result)
	if err != nil {
		client.logger.Error("Error creating credential connection", zap.Error(err))
		return nil, err
	}
	return &result.Data, nil
}

// Update Credential Connection
func (client *TelnyxClient) UpdateCredentialConnection(credentialConnectionID string, profile FQDNConnection) (*FQDNConnection, error) {
	var result struct {
		Data FQDNConnection `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/credential_connections/%s", credentialConnectionID), profile, &result)
	if err != nil {
		client.logger.Error("Error updating credential connection", zap.Error(err), zap.String("credentialConnectionID", credentialConnectionID))
		return nil, err
	}
	return &result.Data, nil
}

// Delete Credential Connection
func (client *TelnyxClient) DeleteCredentialConnection(credentialConnectionID string) error {
	err := client.doRequest("DELETE", fmt.Sprintf("/credential_connections/%s", credentialConnectionID), nil, nil)
	if err != nil {
		client.logger.Error("Error deleting credential connection", zap.Error(err), zap.String("credentialConnectionID", credentialConnectionID))
	}
	return err
}
