package telnyx

import (
	"fmt"
)

func (client *TelnyxClient) CreateFQDN(connectionID, fqdn, dnsRecordType string, port int) (*FQDN, error) {
	body := map[string]interface{}{
		"connection_id":   connectionID,
		"fqdn":            fqdn,
		"dns_record_type": dnsRecordType,
		"port":            port,
	}
	var result struct {
		Data FQDN `json:"data"`
	}
	err := client.doRequest("POST", "/fqdns", body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateFQDN(fqdnID string, fqdn, dnsRecordType string, port int) (*FQDN, error) {
	body := map[string]interface{}{
		"fqdn":            fqdn,
		"dns_record_type": dnsRecordType,
		"port":            port,
	}
	var result struct {
		Data FQDN `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/fqdns/%s", fqdnID), body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeleteFQDN(fqdnID string) error {
	return client.doRequest("DELETE", fmt.Sprintf("/fqdns/%s", fqdnID), nil, nil)
}
