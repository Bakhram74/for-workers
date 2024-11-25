package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatPhone(t *testing.T) {

	var tests = []struct {
		name  string
		input string
		want  string
	}{
		{name: "valid", input: "8-919-(327)30-90", want: "79193273090"},
		{name: "valid-2", input: "+79193273090", want: "79193273090"},
		{name: "valid-3", input: "7(919)32730-90", want: "79193273090"},
		{name: "valid-4", input: "79193273090", want: "79193273090"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidatePhone(tt.input)
			require.NoError(t, err)
			require.Equal(t, got, tt.want)
		})
	}

}

func TestInValidPhone(t *testing.T) {

	var tests = []struct {
		name  string
		input string
		want  string
	}{
		{name: "invalid", input: "9-919-327-3090", want: "invalid phone number"},
		{name: "invalid-1", input: "8-919-327-90", want: "invalid phone number"},
		{name: "invalid-2", input: "791932730901", want: "invalid phone number"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidatePhone(tt.input)
			require.Error(t, err)
			require.Equal(t, err.Error(), tt.want)
			require.Equal(t, got, "")
		})
	}
}

func TestRecoverUUID(t *testing.T) {

	id1 := "4877866c-2a06-4db1-bcd7-e93b1bc0da82"
	id2 := "79c3364e-a58e-4b2f-895d-b24e5c621b6b"
	id3 := "a7889f03-0e5b-4143-bf47-19610ec37f70"

	var tests = []struct {
		name  string
		input string
		want  string
	}{
		{name: "valid", input: strings.ReplaceAll(id1, "-", ""), want: id1},
		{name: "valid2", input: strings.ReplaceAll(id2, "-", ""), want: id2},
		{name: "valid3", input: strings.ReplaceAll(id3, "-", ""), want: id3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RecoverUUID(tt.input)

			require.Equal(t, got, tt.want)
		})
	}
}
