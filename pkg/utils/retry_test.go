package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	_, err := Retry(3, time.Second, func() error {
		val := 1
		if val == 1 {
			return nil
		} else {
			return errors.New("usual error")
		}
	})

	require.NoError(t, err)

	var tests = []struct {
		name  string
		input int
		want  string
	}{
		{name: "usual", input: 2, want: "usual error"},
		{name: "stop", input: 1, want: "stop error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Retry(4, time.Second, func() error {
				switch {
				case tt.input >= 2:
					return errors.New("usual error")
				case tt.input >= 1:
					return Stop{errors.New("stop error")}
				default:
					return nil
				}
			})
			require.Error(t, err)
			require.EqualError(t, err, tt.want)
		})
	}
}
