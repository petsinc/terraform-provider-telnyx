package telnyx

import (
	"fmt"
)

func (client *TelnyxClient) CreateBillingGroup(name string) (*BillingGroup, error) {
	body := map[string]interface{}{
		"name": name,
	}
	var result struct {
		Data BillingGroup `json:"data"`
	}
	err := client.doRequest("POST", "/billing_groups", body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateBillingGroup(billingGroupID, name string) (*BillingGroup, error) {
	body := map[string]interface{}{
		"name": name,
	}
	var result struct {
		Data BillingGroup `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/billing_groups/%s", billingGroupID), body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeleteBillingGroup(billingGroupID string) error {
	return client.doRequest("DELETE", fmt.Sprintf("/billing_groups/%s", billingGroupID), nil, nil)
}
