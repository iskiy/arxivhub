package repository

import (
	models "arxivhub/internal/models"
	"context"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg models.RegisterUserRequest) (models.User, error)
	GetUser(ctx context.Context, id int64) (models.User, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	UpdateUserEmail(ctx context.Context, arg models.UpdateUserEmailParams) (models.User, error)
}

type PaperRepository interface {
	CreatePaper(ctx context.Context, arg models.CreatePaperParams) (models.Paper, error)
	GetPaper(ctx context.Context, arxivID string) (models.Paper, error)
	DeletePaper(ctx context.Context, arxivID string) error
	DeleteSavedPaper(ctx context.Context, arg models.DeleteSavedPaperParams) error
	GetPapers(ctx context.Context, arg models.GetPapersLimit) ([]models.Paper, error)
	GetSavedPaperIDs(ctx context.Context, userID int64) ([]string, error)
	SavePaperForUser(ctx context.Context, arg models.SavePaperForUserParams) (models.SavedPaper, error)
}
