package telnyx

import (
	"fmt"

	"go.uber.org/zap"
)

// Create Credential Connection
func (client *TelnyxClient) CreateCredentialConnection(profile CredentialConnection) (*CredentialConnection, error) {
	var result struct {
		Data CredentialConnection `json:"data"`
	}
	err := client.doRequest("POST", "/credential_connections", profile, &result)
	if err != nil {
		client.logger.Error("Error creating credential connection", zap.Error(err))
		return nil, err
	}
	return &result.Data, nil
}

// Update Credential Connection
func (client *TelnyxClient) UpdateCredentialConnection(credentialConnectionID string, profile CredentialConnection) (*CredentialConnection, error) {
	var result struct {
		Data CredentialConnection `json:"data"`
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

// Get Credential Connection
func (client *TelnyxClient) GetCredentialConnection(credentialConnectionID string) (*CredentialConnection, error) {
	var result struct {
		Data CredentialConnection `json:"data"`
	}
	err := client.doRequest("GET", fmt.Sprintf("/credential_connections/%s", credentialConnectionID), nil, &result)
	if err != nil {
		client.logger.Error("Error getting credential connection", zap.Error(err), zap.String("credentialConnectionID", credentialConnectionID))
		return nil, err
	}
	return &result.Data, nil
}
