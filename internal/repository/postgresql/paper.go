package postgresql

import (
	"arxivhub/internal/models"
	sqlc "arxivhub/internal/repository/postgresql/sqlc"
	"context"
	"database/sql"
)

type PaperPSQLRepository struct {
	queries *sqlc.Queries
}

func NewPaperPSQLRepository(db *sql.DB) *PaperPSQLRepository {
	return &PaperPSQLRepository{queries: sqlc.New(db)}
}

func (r *PaperPSQLRepository) CreatePaper(ctx context.Context, arg models.CreatePaperParams) (models.Paper, error) {
	paper, err := r.queries.CreatePaper(ctx, sqlc.CreatePaperParams(arg))

	return models.Paper(paper), err
}

func (r *PaperPSQLRepository) GetPaper(ctx context.Context, arxivID string) (models.Paper, error) {
	paper, err := r.queries.GetPaper(ctx, arxivID)

	return models.Paper(paper), err
}

func (r *PaperPSQLRepository) DeletePaper(ctx context.Context, arxivID string) error {
	return r.queries.DeletePaper(ctx, arxivID)
}

func (r *PaperPSQLRepository) DeleteSavedPaper(ctx context.Context, arg models.DeleteSavedPaperParams) error {
	return r.queries.DeleteSavedPaper(ctx, sqlc.DeleteSavedPaperParams(arg))
}

func (r *PaperPSQLRepository) GetPapers(ctx context.Context, arg models.GetPapersLimit) ([]models.Paper, error) {
	papers, err := r.queries.GetPapers(ctx, sqlc.GetPapersParams(arg))

	return sqlcPapersToPapers(papers), err
}

func (r *PaperPSQLRepository) SavePaperForUser(ctx context.Context, arg models.SavePaperForUserParams) (models.SavedPaper, error) {
	paper, err := r.queries.SavePaperForUser(ctx, sqlc.SavePaperForUserParams(arg))

	return models.SavedPaper(paper), err
}

func (r *PaperPSQLRepository) GetSavedPaperIDs(ctx context.Context, userID int64) ([]string, error) {
	return r.queries.GetSavedPaperIDs(ctx, userID)
}

func sqlcPapersToPapers(papers []sqlc.Paper) []models.Paper {
	length := len(papers)

	res := make([]models.Paper, length, length)

	for i := 0; i < length; i++ {
		res[i] = models.Paper(papers[i])
	}

	return res
}
