package telnyx

import (
	"fmt"
	"go.uber.org/zap"
)

func (client *TelnyxClient) CreateFQDNConnection(profile FQDNConnection) (*FQDNConnection, error) {
	var result struct {
		Data FQDNConnection `json:"data"`
	}
	err := client.doRequest("POST", "/fqdn_connections", profile, &result)
	if err != nil {
		client.logger.Error("Error creating FQDN connection", zap.Error(err))
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateFQDNConnection(fqdnConnectionID string, profile FQDNConnection) (*FQDNConnection, error) {
	var result struct {
		Data FQDNConnection `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/fqdn_connections/%s", fqdnConnectionID), profile, &result)
	if err != nil {
		client.logger.Error("Error updating FQDN connection", zap.Error(err), zap.String("fqdnConnectionID", fqdnConnectionID))
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeleteFQDNConnection(connectionID string) error {
	err := client.doRequest("DELETE", fmt.Sprintf("/fqdn_connections/%s", connectionID), nil, nil)
	if err != nil {
		client.logger.Error("Error deleting FQDN connection", zap.Error(err), zap.String("connectionID", connectionID))
	}
	return err
}

func (client *TelnyxClient) GetFQDNConnection(fqdnConnectionID string) (*FQDNConnection, error) {
	var result struct {
		Data FQDNConnection `json:"data"`
	}
	err := client.doRequest("GET", fmt.Sprintf("/fqdn_connections/%s", fqdnConnectionID), nil, &result)
	if err != nil {
		client.logger.Error("Error fetching FQDN connection", zap.Error(err), zap.String("fqdnConnectionID", fqdnConnectionID))
		return nil, err
	}
	return &result.Data, nil
}
