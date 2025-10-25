package helpers

import (
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "Positive length",
			length:  6,
			wantErr: false,
		},
		{
			name:    "Zero length",
			length:  0,
			wantErr: true,
		},
		{
			name:    "Negative length",
			length:  -5,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateShortCode(tt.length)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateShortCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(got) != tt.length {
					t.Errorf("GenerateShortCode() got length = %v, want %v", len(got), tt.length)
				}
				// Check if the generated code contains only characters from the charset
				for _, r := range got {
					if !contains(charset, r) {
						t.Errorf("GenerateShortCode() got character %c not in charset", r)
					}
				}
			}
		})
	}
}

// Helper function to check if a rune is present in a string
func contains(s string, r rune) bool {
	for _, char := range s {
		if char == r {
			return true
		}
	}
	return false
}