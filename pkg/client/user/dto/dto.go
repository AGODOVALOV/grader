package dto

type TaskData struct {
	ID        int
	Title     string
	Status    string
	Message   string
	UpdatedAt string
}

type AccountPageData struct {
	ID    int
	Login string
	Name  string
	Tasks []TaskData
}
