package Dao

import "testing"

func TestParsePriceString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{name: "euro suffix", input: "499€", want: 499},
		{name: "european decimal comma", input: "12,54€", want: 12.54},
		{name: "currency prefix", input: "€ 12.54", want: 12.54},
		{name: "thousands and decimal separators", input: "¥ 1,234.56", want: 1234.56},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePriceString(tt.input)
			if err != nil {
				t.Fatalf("parsePriceString(%q): %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("parsePriceString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
