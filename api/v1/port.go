package v1

import (
	"errors"
)

// Port describes port's properties.
type Port struct {
	// City is a city of a port.
	City string `json:"city,omitempty"`
	// Coordinates is a coordinates of a port.
	Coordinates []float64 `json:"coordinates,omitempty"`
	// Country is a country of a port.
	Country string `json:"country,omitempty"`
	// Name is a name of a port.
	Name string `json:"name,omitempty"`
	// Province is a province of a port.
	Province string `json:"province,omitempty"`
}

// Validate validates port input data.
func (p Port) Validate() error {
	if len(p.Name) == 0 {
		return errors.New("port's name can not be empty")
	}

	if len(p.Coordinates) == 0 {
		return errors.New("port's coordinates can not be empty")
	} else if len(p.Coordinates) != 2 {
		return errors.New("port's coordinates should have only 2 values")
	}

	if len(p.Country) == 0 {
		return errors.New("port's country can not be empty")
	}

	// Let's assume that city and province can be empty.

	return nil
}
