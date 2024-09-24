package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeTag(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  Man United  ", "manunited"},
		{"\tMan\tUnited\t", "manunited"},
		{"Man\nUnited", "manunited"},
		{"Man UNITED", "manunited"},
		{"Man   United", "manunited"},
	}

	for _, test := range tests {
		result := NormalizeTag(test.input)
		assert.Equal(t, test.expected, result, "they should be equal")
	}
}
