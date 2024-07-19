package telnyx

import (
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"strings"
)

func (client *TelnyxClient) GetPhoneNumber(phoneNumberID string) (*PhoneNumberResponse, error) {
	var result struct {
		Data PhoneNumberResponse `json:"data"`
	}
	err := client.doRequest("GET", fmt.Sprintf("/phone_numbers/%s", phoneNumberID), nil, &result)
	if err != nil {
		client.logger.Error("Error retrieving phone number", zap.Error(err), zap.String("phone_number_id", phoneNumberID))
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) UpdatePhoneNumber(phoneNumberID string, request UpdatePhoneNumberRequest) (*UpdatePhoneNumberResponse, error) {
	var result struct {
		Data UpdatePhoneNumberResponse `json:"data"`
	}
	err := client.doRequest("PATCH", fmt.Sprintf("/phone_numbers/%s", phoneNumberID), request, &result)
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

// ListAvailablePhoneNumbers retrieves available phone numbers based on the provided filters.
func (client *TelnyxClient) ListAvailablePhoneNumbers(filters AvailablePhoneNumbersRequest) (*AvailablePhoneNumbersResponse, error) {
	queryParams := filters.toQueryParams()
	var result AvailablePhoneNumbersResponse
	err := client.doRequest("GET", fmt.Sprintf("/available_phone_numbers?%s", queryParams), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (filters AvailablePhoneNumbersRequest) toQueryParams() string {
	params := url.Values{}
	addParam := func(key, value string) {
		if value != "" {
			params.Set(key, value)
		}
	}
	addParam("filter[phone_number][starts_with]", filters.StartsWith)
	addParam("filter[phone_number][ends_with]", filters.EndsWith)
	addParam("filter[phone_number][contains]", filters.Contains)
	addParam("filter[locality]", filters.Locality)
	addParam("filter[administrative_area]", filters.AdministrativeArea)
	addParam("filter[country_code]", filters.CountryCode)
	addParam("filter[national_destination_code]", filters.NationalDestinationCode)
	addParam("filter[rate_center]", filters.RateCenter)
	addParam("filter[phone_number_type]", filters.PhoneNumberType)
	if len(filters.Features) > 0 {
		params.Set("filter[features]", strings.Join(filters.Features, ","))
	}
	if filters.Limit > 0 {
		params.Set("filter[limit]", fmt.Sprintf("%d", filters.Limit))
	}
	if filters.BestEffort {
		params.Set("filter[best_effort]", "true")
	}
	if filters.Quickship {
		params.Set("filter[quickship]", "true")
	}
	if filters.Reservable {
		params.Set("filter[reservable]", "true")
	}
	if filters.ExcludeHeldNumbers {
		params.Set("filter[exclude_held_numbers]", "true")
	}
	return params.Encode()
}
