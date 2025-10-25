package helpers

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{
			name: "Empty string",
			data: "",
			want: true,
		},
		{
			name: "Non-empty string",
			data: "hello",
			want: false,
		},
		{
			name: "String with spaces",
			data: " ",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.data); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}