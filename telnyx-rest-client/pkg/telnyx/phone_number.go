package telnyx

import (
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"strconv"
	"strings"
)

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
	params := url.Values{}
	if filters.StartsWith != "" {
		params.Set("filter[phone_number][starts_with]", filters.StartsWith)
	}
	if filters.EndsWith != "" {
		params.Set("filter[phone_number][ends_with]", filters.EndsWith)
	}
	if filters.Contains != "" {
		params.Set("filter[phone_number][contains]", filters.Contains)
	}
	if filters.Locality != "" {
		params.Set("filter[locality]", filters.Locality)
	}
	if filters.AdministrativeArea != "" {
		params.Set("filter[administrative_area]", filters.AdministrativeArea)
	}
	if filters.CountryCode != "" {
		params.Set("filter[country_code]", filters.CountryCode)
	}
	if filters.NationalDestinationCode != "" {
		params.Set("filter[national_destination_code]", filters.NationalDestinationCode)
	}
	if filters.RateCenter != "" {
		params.Set("filter[rate_center]", filters.RateCenter)
	}
	if filters.PhoneNumberType != "" {
		params.Set("filter[phone_number_type]", filters.PhoneNumberType)
	}
	if len(filters.Features) > 0 {
		params.Set("filter[features]", strings.Join(filters.Features, ","))
	}
	if filters.Limit > 0 {
		params.Set("filter[limit]", strconv.Itoa(filters.Limit))
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

	var result AvailablePhoneNumbersResponse
	err := client.doRequest("GET", fmt.Sprintf("/available_phone_numbers?%s", params.Encode()), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
