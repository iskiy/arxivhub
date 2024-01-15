package main

import (
	"arxivhub/internal/hubpopulator/vectorencoder"
	"arxivhub/internal/models"
	"arxivhub/internal/pinecone"
	"arxivhub/internal/repository/postgresql"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/orijtech/arxiv/v1"
	"os"

	"log"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	resChan, cancel, err := arxiv.Search(context.Background(), &arxiv.Query{
		Terms:             "cat:cs.AI",
		MaxPageNumber:     20,
		MaxResultsPerPage: 100,
	})
	if err != nil {
		log.Fatal(err)
	}

	for resPage := range resChan {
		var papers []models.Paper
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

		if resPage.PageNumber >= int64(20) {
			cancel()
		}
		AddPapers(papers)
	}

}

func AddPapers(papers []models.Paper) {
	encodeClient := vectorencoder.NewClient(os.Getenv("ENCODER_URL"))
	vectors := make([]pinecone.Vector, len(papers), len(papers))

	for i, p := range papers {
		encoded, err := encodeClient.Encode(stringForPaperEncoding(p))
		if err != nil {
			log.Fatalln(err)
		}

		vectors[i] = pinecone.Vector{
			ID:     p.ArxivID,
			Values: encoded.Output,
		}
	}

	fmt.Println(len(vectors))

	insertRequest := pinecone.InsertRequest{Vector: vectors}
	data, err := json.Marshal(&insertRequest)
	if err != nil {
		log.Fatalln(err)
	}

	conn, err := postgresql.ConnectToDB(os.Getenv("PSQL_CONN_STRING"))

	if err != nil {
		log.Fatalln(err)
	}

	paperRepository := postgresql.NewPaperPSQLRepository(conn)

	pineconeClient := pinecone.NewClient(os.Getenv("PINECONE_KEY"), os.Getenv("PINECONE_HOST"))
	err = pineconeClient.InsertVectors(data)
	if err != nil {
		log.Fatalln(err)
	}
	for _, p := range papers {
		paperRepository.CreatePaper(context.Background(), models.CreatePaperParams{
			ArxivID:  p.ArxivID,
			Title:    p.Title,
			Abstract: p.Abstract,
			Authors:  p.Authors,
			Date:     p.Date,
		})
	}
}

func stringForPaperEncoding(paper models.Paper) string {
	return "Title: " + paper.Title + " Abstract: " + paper.Abstract
}
