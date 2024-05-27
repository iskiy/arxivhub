package postgresql

import (
	"arxivhub/internal/models"
	"context"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func TestPaperRepository_Create_Get(t *testing.T) {
	type args struct {
		ctx         context.Context
		args        []models.CreatePaperParams
		getArxivIDs []string
	}
	tests := []struct {
		name          string
		args          args
		wantErrCreate bool
		wantErrGet    bool
	}{
		{
			name: "Create and Get Paper",
			args: args{
				ctx: context.Background(),
				args: []models.CreatePaperParams{
					{
						ArxivID:  "1234.5678",
						Title:    "Test Paper",
						Abstract: "This is a sample abstract.",
						Authors:  "Author One",
						Date:     time.Now(),
					},
					{
						ArxivID:  "1234.5679",
						Title:    "Test Paper 2",
						Abstract: "This is a sample abstract.",
						Authors:  "Author One",
						Date:     time.Now(),
					},
				},

				getArxivIDs: []string{"1234.5678", "1234.5679"},
			},
			wantErrCreate: false,
			wantErrGet:    false,
		},
		{
			name: "Create duplicate",
			args: args{
				ctx: context.Background(),
				args: []models.CreatePaperParams{
					{
						ArxivID:  "1234.5678",
						Title:    "Test Paper",
						Abstract: "This is a sample abstract.",
						Authors:  "Author One",
						Date:     time.Now(),
					},
					{
						ArxivID:  "1234.5678",
						Title:    "Test Paper 2",
						Abstract: "This is a sample abstract.",
						Authors:  "Author One",
						Date:     time.Now(),
					},
				},

				getArxivIDs: []string{"1234.5678"},
			},
			wantErrCreate: true,
			wantErrGet:    false,
		},
		{
			name: "Get not created",
			args: args{
				ctx: context.Background(),
				args: []models.CreatePaperParams{
					{
						ArxivID:  "1234.5678",
						Title:    "Test Paper",
						Abstract: "This is a sample abstract.",
						Authors:  "Author One",
						Date:     time.Now(),
					},
				},

				getArxivIDs: []string{"1234.5679"},
			},
			wantErrCreate: false,
			wantErrGet:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer clearFunc()
			r := NewPaperPSQLRepository(dbConn)

			for _, arg := range tt.args.args {
				_, err := r.CreatePaper(tt.args.ctx, arg)
				if err != nil && tt.wantErrCreate == false {
					t.Errorf("PaperPSQLRepository.CreatePaper() error = %v, wantErr %v", err, tt.wantErrCreate)
					return
				}
			}

			for i, arxivID := range tt.args.getArxivIDs {
				if arxivID == "" {
					continue
				}

				got, err := r.GetPaper(tt.args.ctx, arxivID)
				if tt.wantErrGet == true {
					continue
				}

				if err != nil {
					t.Errorf("PaperPSQLRepository.GetPaper() error = %v, wantErr %v", err, tt.wantErrGet)
					return
				}

				expected := tt.args.args[i]
				if got.ArxivID != expected.ArxivID || got.Title != expected.Title ||
					got.Abstract != expected.Abstract || got.Authors != expected.Authors {
					t.Errorf("PaperPSQLRepository.GetPaper() = %v, want %v", got, expected)
				}
			}
		})
	}
}
