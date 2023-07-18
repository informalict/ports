package v1

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPort_Validate tests port validation.
func TestPort_Validate(t *testing.T) {
	type fields struct {
		City        string
		Coordinates []float64
		Country     string
		Name        string
		Province    string
	}
	tests := map[string]struct {
		fields  fields
		wantErr error
	}{
		"empty name": {
			wantErr: errors.New("port's name can not be empty"),
		},
		"empty coordinates": {
			fields: fields{
				Name: "test",
			},
			wantErr: errors.New("port's coordinates can not be empty"),
		},
		"invalid coordinates": {
			fields: fields{
				Name:        "test",
				Coordinates: []float64{1, 1, 1},
			},
			wantErr: errors.New("port's coordinates should have only 2 values"),
		},
		"empty city": {
			fields: fields{
				Name:        "test",
				Coordinates: []float64{1, 1},
			},
			wantErr: errors.New("port's country can not be empty"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			p := Port{
				City:        test.fields.City,
				Coordinates: test.fields.Coordinates,
				Country:     test.fields.Country,
				Name:        test.fields.Name,
				Province:    test.fields.Province,
			}

			err := p.Validate()
			if test.wantErr != nil {
				require.EqualError(t, err, test.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
