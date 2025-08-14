package aws

import "errors"

var (
	ErrFailedToUnmarshalGeoJSON  = errors.New("failed to unmarshal GeoJSON payload")
	ErrInvalidGeoJSONCoordinates = errors.New("invalid GeoJSON point: coordinates must have at least 2 elements")
	ErrPositionResolutionIsEmpty = errors.New("position resolution is empty")
)
