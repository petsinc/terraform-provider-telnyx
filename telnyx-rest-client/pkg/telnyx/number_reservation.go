package telnyx

import (
	"fmt"
)

func (client *TelnyxClient) CreateNumberReservation(phoneNumbers []string, customerReference string) (*PhoneNumberReservation, error) {
	var phoneNumbersMap []map[string]string
	for _, phoneNumber := range phoneNumbers {
		phoneNumbersMap = append(phoneNumbersMap, map[string]string{"phone_number": phoneNumber})
	}
	body := map[string]interface{}{
		"phone_numbers":      phoneNumbersMap,
		"customer_reference": customerReference,
	}
	var result struct {
		Data PhoneNumberReservation `json:"data"`
	}
	err := client.doRequest("POST", "/number_reservations", body, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}

func (client *TelnyxClient) ExtendPhoneNumberReservation(reservationID string) (*PhoneNumberReservation, error) {
	var result struct {
		Data PhoneNumberReservation `json:"data"`
	}
	err := client.doRequest("POST", fmt.Sprintf("/number_reservations/%s/actions/extend", reservationID), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result.Data, nil
}
