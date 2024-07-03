package telnyx

import (
	"fmt"
	"go.uber.org/zap"
)

func (client *TelnyxClient) CreateFQDN(fqdn FQDN) (*FQDN, error) {
	var result struct {
		Data FQDN `json:"data"`
	}
	err := client.doRequest("POST", "/fqdns", fqdn, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateFQDN(fqdnID string, fqdn FQDN) (*FQDN, error) {
	var result struct {
		Data FQDN `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/fqdns/%s", fqdnID), fqdn, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeleteFQDN(fqdnID string) error {
	err := client.doRequest("DELETE", fmt.Sprintf("/fqdns/%s", fqdnID), nil, nil)
	if err != nil {
		client.logger.Error("Error deleting FQDN", zap.Error(err), zap.String("fqdnID", fqdnID))
	}
	return err
}
