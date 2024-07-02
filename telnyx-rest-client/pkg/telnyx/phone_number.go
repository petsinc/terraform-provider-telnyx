package telnyx

import (
	"fmt"
	"go.uber.org/zap"
)

func (client *TelnyxClient) UpdatePhoneNumber(phoneNumberID, customerReference, connectionID, billingGroupID string, tags []string, hdVoiceEnabled bool) (*PhoneNumber, error) {
	body := map[string]interface{}{
		"customer_reference": customerReference,
		"connection_id":      connectionID,
		"billing_group_id":   billingGroupID,
		"tags":               tags,
		"hd_voice_enabled":   hdVoiceEnabled,
	}
	var result struct {
		Data PhoneNumber `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/phone_numbers/%s", phoneNumberID), body, &result)
	if err != nil {
		client.logger.Error("Error updating phone number", zap.Error(err), zap.String("phone_number_id", phoneNumberID))
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) DeletePhoneNumber(phoneNumberID string) error {
	err := client.doRequest("DELETE", fmt.Sprintf("/phone_numbers/%s", phoneNumberID), nil, nil)
	if err != nil {
		client.logger.Error("Error deleting phone number", zap.Error(err), zap.String("phone_number_id", phoneNumberID))
	}
	return err
}
