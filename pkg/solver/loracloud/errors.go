package loracloud

import "errors"

var (
	ErrContextPortNotFound       = errors.New("context port not found")
	ErrContextDevEuiNotFound     = errors.New("context DevEUI not found")
	ErrContextFCountNotFound     = errors.New("context frame counter not found")
	ErrContextInvalidPort        = errors.New("context port is invalid, must be a number between 0 and 255")
	ErrContextInvalidDevEui      = errors.New("context DevEUI is invalid, must be a valid hex string of length 16")
	ErrContextInvalidFCount      = errors.New("context frame counter is invalid, must be a positive integer")
	ErrSemtechLoRaCloudShutdown  = errors.New("LoRa Cloud is no longer available after 31.07.2025, see https://www.semtech.com/loracloud-shutdown")
	ErrSendingRequest            = errors.New("error sending request")
	ErrUnexpectedStatusCode      = errors.New("unexpected status code returned")
	ErrDecodingResponse          = errors.New("error decoding response")
	ErrMultipleDevicesInResponse = errors.New("multiple devices found in response")
	ErrDeviceEuiNotInResponse    = errors.New("device EUI not found in response")
)
