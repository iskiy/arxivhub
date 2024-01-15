package arxiv

import (
	"arxivhub/internal/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/orijtech/arxiv/v1"
)

func FindPapers(query string, maxPageNumber, maxResultPerPage int) []models.Paper {
	var papers []models.Paper

	resChan, cancel, err := arxiv.Search(context.Background(), &arxiv.Query{
		Terms:             query,
		MaxPageNumber:     5,
		MaxResultsPerPage: 100,
	})
	if err != nil {
		log.Fatal(err)
	}

	for resPage := range resChan {
		if err := resPage.Err; err != nil {
			fmt.Printf("#%d err: %v", resPage.PageNumber, err)
			continue
		}

		fmt.Printf("#Page number: %d\n", resPage.PageNumber)
		feed := resPage.Feed

		for _, entry := range feed.Entry {
			date, err := time.Parse("2006-01-02T15:04:05Z", string(entry.Published))
			if err != nil {
				fmt.Println(err)
			}

			paper := models.Paper{
				ArxivID:  entry.ID,
				Title:    entry.Title,
				Abstract: entry.Summary.Body,
				Authors:  entry.Author.Name,
				Date:     date,
			}

			papers = append(papers, paper)
		}

		if resPage.PageNumber >= int64(maxPageNumber+1) {
			cancel()
		}
	}
	return papers
}
