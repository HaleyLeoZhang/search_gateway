package vo

type BlogFrontRequest struct {
	Id       *int   `form:"id"`
	Title    string `form:"title"`
	Describe string `form:"describe"`
	Category string `form:"category"`

	Page     *int `form:"page"`
	PageSize *int `form:"page_size"`
}

func (s *BlogFrontRequest) SetPage() {
	if s.Page == nil || *s.Page < 1 {
		page := 1
		s.Page = &page
	}
}

func (s *BlogFrontRequest) SetPageSize() {
	if s.PageSize == nil || *s.PageSize > 50 {
		pageSize := 10
		s.PageSize = &pageSize
	}
}

func (s *BlogFrontRequest) GetOffset() int {
	s.SetPageSize()
	s.SetPage()
	return (*s.Page - 1) * (*s.PageSize)
}

type BlogFrontResponse struct {
	Ids   []int `json:"ids"`
	Total int64 `json:"total"`
}
