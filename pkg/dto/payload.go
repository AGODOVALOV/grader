package dto

type GraderPayload struct {
	UserID        string `json:"userId"`
	TaskID        string `json:"taskId"`
	ReviewID      string `json:"reviewId"`
	FileIDs       []File `json:"files"`
	EventID       string `json:"eventId"`
	ContainerName string `json:"containerName"`
}

type File struct {
	Label    string `json:"label"`
	FileName string `json:"fileName"`
}
