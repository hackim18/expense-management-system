package utils

import "go-expense-management-system/internal/model"

func NewPageMetadata(page, size int, total int64) model.PageMetadata {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	totalPage := total / int64(size)
	if total%int64(size) != 0 {
		totalPage++
	}

	return model.PageMetadata{
		CurrentPage: page,
		PageSize:    size,
		TotalItem:   total,
		TotalPage:   totalPage,
		HasNext:     int64(page) < totalPage,
		HasPrevious: page > 1,
	}
}
