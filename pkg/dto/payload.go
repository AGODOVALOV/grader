package dto

// GraderPayload is a struct that represents the payload sent to the grader
type GraderPayload struct {
	UserID        string `json:"user_id"`
	TaskID        string `json:"task_id"`
	ReviewID      string `json:"review_id"`
	FileIDs       []File `json:"files"`
	EventID       string `json:"event_id"`
	ContainerName string `json:"container_name"`
}

// File is a struct that represents a file that is sent to the grader
type File struct {
	Label    string `json:"label"`
	FileName string `json:"file_name"`
}

// GraderPayloadCallback is a struct that represents the payload sent to the callback
type GraderPayloadCallback struct {
	UserID        string `json:"user_id"`
	TaskID        string `json:"task_id"`
	ReviewID      string `json:"review_id"`
	EventID       string `json:"event_id"`
	Passed        bool   `json:"passed"`
	OutputMessage string `json:"error_message"`
	ErrorText     string `json:"error_text"`
}
