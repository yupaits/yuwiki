package yuwiki

const (
	defaultPage = 1
	defaultSize = 10
)

type Pager struct {
	Page       uint
	Size       uint
	Total      uint
	Content    interface{}
	HasPrev    bool
	HasNext    bool
	TotalPages uint
}

func Paging(page uint, size uint, total uint, data interface{}) *Pager {
	totalPages := total / size
	if total%size != 0 {
		totalPages += 1
	}
	pager := &Pager{
		Page:       page,
		Size:       size,
		Total:      total,
		Content:    data,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
	return pager
}
