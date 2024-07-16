package telnyx

import (
	"fmt"
)

func (client *TelnyxClient) CreateNumberOrder(request CreateNumberOrderRequest) (*PhoneNumberOrderResponse, error) {
	var result struct {
		Data PhoneNumberOrderResponse `json:"data"`
	}
	err := client.doRequest("POST", "/number_orders", request, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdateNumberOrder(numberOrderID string, request UpdateNumberOrderRequest) (*PhoneNumberOrderResponse, error) {
	var result struct {
		Data PhoneNumberOrderResponse `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/number_orders/%s", numberOrderID), request, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) GetNumberOrder(numberOrderID string) (*PhoneNumberOrderResponse, error) {
	var result struct {
		Data PhoneNumberOrderResponse `json:"data"`
	}
	err := client.doRequest("GET", fmt.Sprintf("/number_orders/%s", numberOrderID), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}
