package telnyx

import (
	"fmt"
)

func (client *TelnyxClient) CreateNumberOrder(phoneNumbers []string, connectionID, messagingProfileID, billingGroupID, customerReference string) (*PhoneNumberOrder, error) {
	var phoneNumbersMap []map[string]string
	for _, phoneNumber := range phoneNumbers {
		phoneNumbersMap = append(phoneNumbersMap, map[string]string{"phone_number": phoneNumber})
	}
	body := map[string]interface{}{
		"phone_numbers":        phoneNumbersMap,
		"connection_id":        connectionID,
		"messaging_profile_id": messagingProfileID,
		"billing_group_id":     billingGroupID,
		"customer_reference":   customerReference,
	}
	var result struct {
		Data PhoneNumberOrder `json:"data"`
	}
	err := client.doRequest("POST", "/number_orders", body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateNumberOrder(numberOrderID, customerReference string, regulatoryRequirements []map[string]string) (*PhoneNumberOrder, error) {
	body := map[string]interface{}{
		"customer_reference":      customerReference,
		"regulatory_requirements": regulatoryRequirements,
	}
	var result struct {
		Data PhoneNumberOrder `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/number_orders/%s", numberOrderID), body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}
