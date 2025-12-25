package utils

import "testing"

func TestFormatIDR(t *testing.T) {
	tests := []struct {
		amount int64
		want   string
	}{
		{amount: 0, want: "Rp 0"},
		{amount: 1500, want: "Rp 1.500"},
		{amount: 1500000, want: "Rp 1.500.000"},
	}

	for _, tt := range tests {
		got := FormatIDR(tt.amount)
		if got != tt.want {
			t.Fatalf("expected %q, got %q", tt.want, got)
		}
	}
}
