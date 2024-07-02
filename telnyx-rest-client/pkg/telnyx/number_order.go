package telnyx

import (
	"fmt"
)

func (client *TelnyxClient) CreateNumberOrder(request PhoneNumberOrderRequest) (*PhoneNumberOrderResponse, error) {
	phoneNumbers := []map[string]interface{}{}
	for _, number := range request.PhoneNumbers {
		phoneNumbers = append(phoneNumbers, map[string]interface{}{
			"phone_number": number,
		})
	}

	body := map[string]interface{}{
		"phone_numbers":        phoneNumbers,
		"connection_id":        request.ConnectionID,
		"messaging_profile_id": request.MessagingProfileID,
		"billing_group_id":     request.BillingGroupID,
		"customer_reference":   request.CustomerReference,
		"sub_number_order_ids": request.SubNumberOrderIDs,
	}

	var result struct {
		Data PhoneNumberOrderResponse `json:"data"`
	}
	err := client.doRequest("POST", "/number_orders", body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateNumberOrder(numberOrderID string, request UpdateNumberOrderRequest) (*PhoneNumberOrderResponse, error) {
	regulatoryRequirements := []map[string]interface{}{}
	for _, req := range request.RegulatoryRequirements {
		regulatoryRequirements = append(regulatoryRequirements, map[string]interface{}{
			"requirement_id": req.RequirementID,
			"field_value":    req.FieldValue,
		})
	}

	body := map[string]interface{}{
		"customer_reference":      request.CustomerReference,
		"regulatory_requirements": regulatoryRequirements,
	}

	var result struct {
		Data PhoneNumberOrderResponse `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/number_orders/%s", numberOrderID), body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}
