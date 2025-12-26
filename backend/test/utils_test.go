package test

import (
	"testing"

	"go-expense-management-system/internal/utils"

	"github.com/stretchr/testify/require"
)

func TestFormatIDR(t *testing.T) {
	tests := []struct {
		name   string
		amount int64
		want   string
	}{
		{name: "zero", amount: 0, want: "Rp 0"},
		{name: "thousand", amount: 1500, want: "Rp 1.500"},
		{name: "million", amount: 1500000, want: "Rp 1.500.000"},
	}

	for _, tt := range tests {
		got := utils.FormatIDR(tt.amount)
		require.Equal(t, tt.want, got, tt.name)
	}
}

func TestNewPageMetadata(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		size        int
		total       int64
		wantPage    int
		wantSize    int
		wantTotal   int64
		wantPages   int64
		wantHasNext bool
		wantHasPrev bool
	}{
		{
			name:        "defaults-with-empty-total",
			page:        0,
			size:        0,
			total:       0,
			wantPage:    1,
			wantSize:    10,
			wantTotal:   0,
			wantPages:   0,
			wantHasNext: false,
			wantHasPrev: false,
		},
		{
			name:        "middle-page",
			page:        2,
			size:        10,
			total:       25,
			wantPage:    2,
			wantSize:    10,
			wantTotal:   25,
			wantPages:   3,
			wantHasNext: true,
			wantHasPrev: true,
		},
		{
			name:        "last-page",
			page:        3,
			size:        10,
			total:       25,
			wantPage:    3,
			wantSize:    10,
			wantTotal:   25,
			wantPages:   3,
			wantHasNext: false,
			wantHasPrev: true,
		},
	}

	for _, tt := range tests {
		meta := utils.NewPageMetadata(tt.page, tt.size, tt.total)
		require.Equal(t, tt.wantPage, meta.CurrentPage, tt.name)
		require.Equal(t, tt.wantSize, meta.PageSize, tt.name)
		require.Equal(t, tt.wantTotal, meta.TotalItem, tt.name)
		require.Equal(t, tt.wantPages, meta.TotalPage, tt.name)
		require.Equal(t, tt.wantHasNext, meta.HasNext, tt.name)
		require.Equal(t, tt.wantHasPrev, meta.HasPrevious, tt.name)
	}
}
