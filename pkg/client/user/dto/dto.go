package dto

// TaskData represents one task row for the user panel.
type TaskData struct {
	ID        int
	Title     string
	Status    string
	Message   string
	UpdatedAt string
}

// AccountPageData represents data for rendering the user account page.
type AccountPageData struct {
	ID    int
	Login string
	Name  string
	Tasks []TaskData
}

// AdminReviewData represents one review row for the admin panel.
type AdminReviewData struct {
	ID        int64
	UserLogin string
	TaskTitle string
	Status    string
	Message   string
	FileName  string
	CreatedAt string
	UpdatedAt string
}

// AdminReviewsPageData represents data for rendering the admin reviews page.
type AdminReviewsPageData struct {
	Reviews []AdminReviewData
}
