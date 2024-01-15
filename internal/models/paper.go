package models

import "time"

type Paper struct {
	ArxivID  string    `json:"arxiv_id"`
	Title    string    `json:"title"`
	Abstract string    `json:"abstract"`
	Authors  string    `json:"authors"`
	Date     time.Time `json:"date"`
}

type SavedPaper struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	PaperID   string    `json:"paper_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatePaperParams struct {
	ArxivID  string    `json:"arxiv_id"`
	Title    string    `json:"title"`
	Abstract string    `json:"abstract"`
	Authors  string    `json:"authors"`
	Date     time.Time `json:"date"`
}

type DeleteSavedPaperRequest struct {
	PaperID string `json:"paper_id" validate:"required"`
}

type DeleteSavedPaperParams struct {
	UserID  int64  `json:"user_id"`
	PaperID string `json:"paper_id"`
}

type GetPapersLimit struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type GetSavedPapersForUserParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type SavePaperRequest struct {
	PaperID string `json:"paper_id" validate:"required"`
}

type SavePaperForUserParams struct {
	UserID  int64  `json:"user_id"`
	PaperID string `json:"paper_id"`
}

// PaperData contains information about the article without arxivID
type PaperData struct {
	Title           string    `json:"title"`
	Abstract        string    `json:"abstract"`
	Authors         []string  `json:"authors"`
	Categories      []string  `json:"categories"`
	PublicationDate time.Time `json:"publication_date"`
	Saved           bool      `json:"saved"`
	Year            int       `json:"year"`
}

type PineconePaper struct {
	ArxivID string `json:"arxiv_id""`
	PaperData
}

type ScoredPaper struct {
	PineconePaper PineconePaper `json:"paper"`
	Score         float64       `json:"score"`
}

type SearchPaperRequest struct {
	Query    string
	Category string
	Year     int
	Amount   int
}

type SearchPapersResponse struct {
	ScoredPapers []ScoredPaper `json:"scored_papers"`
}
