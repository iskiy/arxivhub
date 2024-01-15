package service

import (
	"arxivhub/internal/hubpopulator/vectorencoder"
	models "arxivhub/internal/models"
	"arxivhub/internal/pinecone"
	"arxivhub/internal/repository"
	"context"
	"fmt"
)

type PaperServiceImpl struct {
	repository     repository.PaperRepository
	encoderClient  *vectorencoder.Client
	pineconeClient *pinecone.Client
}

func NewPaperService(repository repository.PaperRepository, encoderClient *vectorencoder.Client, pineconeClient *pinecone.Client) *PaperServiceImpl {
	return &PaperServiceImpl{
		repository:     repository,
		encoderClient:  encoderClient,
		pineconeClient: pineconeClient,
	}
}

func (s *PaperServiceImpl) DeleteSavedPaper(ctx context.Context, params models.DeleteSavedPaperParams) error {
	return s.repository.DeleteSavedPaper(ctx, params)
}

func (s *PaperServiceImpl) SavePaper(ctx context.Context, params models.SavePaperForUserParams) (models.SavedPaper, error) {
	if params.PaperID == "" {
		return models.SavedPaper{}, fmt.Errorf("empty paper ID")
	}
	papers, err := s.pineconeClient.FetchPapers([]string{params.PaperID})
	if err != nil {
		return models.SavedPaper{}, fmt.Errorf("error during paper ID validation")
	}

	if len(papers.Vectors) == 0 {
		return models.SavedPaper{}, fmt.Errorf("there is no article with id: %s", params.PaperID)
	}

	return s.repository.SavePaperForUser(ctx, params)
}

func (s *PaperServiceImpl) Search(_ context.Context, params models.SearchPaperRequest) ([]models.ScoredPaper, error) {
	encodeResponse, err := s.encoderClient.Encode(params.Query)
	if err != nil {
		return []models.ScoredPaper{}, err
	}

	if params.Amount == 0 {
		params.Amount = 24
	}

	pineRequest := pinecone.QueryRequest{
		IncludeMetadata: true,
		IncludeValues:   false,
		Vector:          encodeResponse.Output,
		TopK:            params.Amount,
		Filter: pinecone.Filter{
			Year:     params.Year,
			Category: params.Category,
		},
	}

	pineconeResponse, err := s.pineconeClient.Query(pineRequest)
	if err != nil {
		return []models.ScoredPaper{}, err
	}

	scoredPapers := make([]models.ScoredPaper, len(pineconeResponse.Matches), len(pineconeResponse.Matches))
	for i, resp := range pineconeResponse.Matches {
		scoredPapers[i] = models.ScoredPaper{
			PineconePaper: models.PineconePaper{
				ArxivID:   resp.ID,
				PaperData: resp.PaperData,
			},
			Score: resp.Score,
		}
	}

	return scoredPapers, nil
}

func (s *PaperServiceImpl) UpdateSavedField(ctx context.Context, userID int64, papers []models.ScoredPaper) error {
	saved, err := s.repository.GetSavedPaperIDs(ctx, userID)
	if err != nil {
		return err
	}
	savedMap := make(map[string]bool)
	for _, p := range saved {
		savedMap[p] = true
	}
	for i := 0; i < len(papers); i++ {
		if savedMap[papers[i].PineconePaper.ArxivID] {
			papers[i].PineconePaper.Saved = true
		}
	}

	return nil
}

func (s *PaperServiceImpl) GetSavedPaper(ctx context.Context, userID int64) ([]models.PineconePaper, error) {
	ids, err := s.repository.GetSavedPaperIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	pineconeResults, err := s.pineconeClient.FetchPapers(ids)
	if err != nil {
		return nil, err
	}

	res := make([]models.PineconePaper, len(pineconeResults.Vectors), len(pineconeResults.Vectors))

	for i := 0; i < len(ids); i++ {
		pineconeEntry := pineconeResults.Vectors[ids[i]]

		res[i] = models.PineconePaper{
			ArxivID:   pineconeEntry.ID,
			PaperData: pineconeEntry.PaperData,
		}
	}

	return res, nil
}
