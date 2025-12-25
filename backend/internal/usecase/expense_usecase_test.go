package usecase

import "testing"

func TestValidateExpenseAmount(t *testing.T) {
	tests := []struct {
		name   string
		amount int64
		wantOK bool
	}{
		{name: "below-min", amount: 5000, wantOK: false},
		{name: "min", amount: 10000, wantOK: true},
		{name: "above-max", amount: 50000001, wantOK: false},
		{name: "max", amount: 50000000, wantOK: true},
	}

	for _, tt := range tests {
		err := validateExpenseAmount(tt.amount)
		if tt.wantOK && err != nil {
			t.Fatalf("%s: expected no error, got %v", tt.name, err)
		}
		if !tt.wantOK && err == nil {
			t.Fatalf("%s: expected error, got nil", tt.name)
		}
	}
}

func TestNormalizeStatusFilter(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "pending", want: "awaiting_approval"},
		{input: "auto-approved", want: "auto_approved"},
		{input: "approved", want: "approved"},
	}

	for _, tt := range tests {
		got := normalizeStatusFilter(tt.input)
		if got != tt.want {
			t.Fatalf("expected %q, got %q", tt.want, got)
		}
	}
}
