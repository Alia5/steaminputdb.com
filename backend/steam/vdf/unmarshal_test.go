package vdf_test

import (
	"testing"

	"github.com/Alia5/steaminputdb.com/steam/vdf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalScalarIntoMap(t *testing.T) {
	type target struct {
		Details map[string]string `json:"details"`
	}

	type testCase struct {
		name        string
		data        string
		expected    map[string]string
		expectedErr bool
	}

	testCases := []testCase{
		{
			name:     "BLOCK",
			data:     `"root" { "details" { "id" "value" } }`,
			expected: map[string]string{"id": "value"},
		},
		{
			// steam sends "" instead of an empty block
			name:     "EMPTY_STRING_IS_NO_ENTRIES",
			data:     `"root" { "details" "" }`,
			expected: nil,
		},
		{
			name:        "NON_EMPTY_STRING_IS_AN_ERROR",
			data:        `"root" { "details" "nope" }`,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res struct {
				Root target `json:"root"`
			}
			err := vdf.Unmarshal(tc.data, &res)
			if tc.expectedErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.expected, res.Root.Details)
		})
	}
}
