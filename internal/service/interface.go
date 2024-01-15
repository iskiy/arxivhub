package service

import (
	"arxivhub/internal/models"
	"context"
)

type PaperService interface {
	SavePaper(ctx context.Context, params models.SavePaperForUserParams) (models.SavedPaper, error)
	DeleteSavedPaper(ctx context.Context, params models.DeleteSavedPaperParams) error
	Search(ctx context.Context, request models.SearchPaperRequest) ([]models.ScoredPaper, error)
	GetSavedPaper(ctx context.Context, userID int64) ([]models.PineconePaper, error)
	UpdateSavedField(ctx context.Context, userID int64, papers []models.ScoredPaper) error
}

type UserService interface {
	RegisterUser(ctx context.Context, params models.RegisterUserRequest) (models.User, error)
	LoginUser(ctx context.Context, params models.LoginUserRequest) (models.User, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
}
