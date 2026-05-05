package dto

type GraderPayload struct {
	UserID        string `json:"userId"`
	ReviewID      string `json:"reviewId"`
	FileIDs       []File `json:"files"`
	ContainerName string `json:"containerName"`
}

type File struct {
	Label    string `json:"label"`
	FileName string `json:"fileName"`
}
