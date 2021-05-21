package simpleforce

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/pkg/errors"
)

var (
	// ErrAuthentication is returned when authentication failed.
	ErrAuthentication = errors.New("authentication failure")
)

type jsonError []struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type xmlError struct {
	Message   string `xml:"Body>Fault>faultstring"`
	ErrorCode string `xml:"Body>Fault>faultcode"`
}

// Need to get information out of this package.
func ParseSalesforceError(statusCode int, responseBody []byte) error {
	jsonError := jsonError{}
	jErr := json.Unmarshal(responseBody, &jsonError)
	if jErr != nil {
		xmlError := xmlError{}
		xErr := xml.Unmarshal(responseBody, &xmlError)
		if xErr != nil {
			return errors.Errorf("error while unmarshaling SFDC errors. JSON: %s, XML: %s", jErr, xErr)
		}

		message := fmt.Sprintf("Error. http code: %v Error Message:  %v Error Code: %v", statusCode, xmlError.Message, xmlError.ErrorCode)
		return errors.New(message)
	}

	message := fmt.Sprintf("Error. http code: %v Error Message:  %v Error Code: %v", statusCode, jsonError[0].Message, jsonError[0].ErrorCode)
	return errors.New(message)
}
