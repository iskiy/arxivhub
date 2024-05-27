package handlers

import (
	"arxivhub/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestGetRequestSearchParams(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedResult models.SearchPaperRequest
		expectedError  error
	}{
		{
			name: "Valid",
			queryParams: map[string]string{
				"query":    "test",
				"year":     "2023",
				"amount":   "10",
				"category": "science",
			},
			expectedResult: models.SearchPaperRequest{
				Query:    "test",
				Category: "science",
				Year:     2023,
				Amount:   10,
			},
			expectedError: nil,
		},
		{
			name: "Empty",
			queryParams: map[string]string{
				"query": "",
			},
			expectedResult: models.SearchPaperRequest{},
			expectedError:  fiber.NewError(fiber.StatusBadRequest, ErrEmptySearchQuery.Error()),
		},
		{
			name: "Invalid year format",
			queryParams: map[string]string{
				"query": "test",
				"year":  "invalid",
			},
			expectedResult: models.SearchPaperRequest{},
			expectedError:  fiber.NewError(fiber.StatusBadRequest, ErrInvalidYear.Error()),
		},
		{
			name: "Invalid amount format",
			queryParams: map[string]string{
				"query":  "test",
				"amount": "invalid",
			},
			expectedResult: models.SearchPaperRequest{},
			expectedError:  fiber.NewError(fiber.StatusBadRequest, ErrInvalidAmountFormat.Error()),
		},
		{
			name: "Amount out of range",
			queryParams: map[string]string{
				"query":  "test",
				"amount": "121",
			},
			expectedResult: models.SearchPaperRequest{},
			expectedError:  fiber.NewError(fiber.StatusBadRequest, ErrInvalidAmount.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Get("/", func(c *fiber.Ctx) error {
				result, err := GetRequestSearchParams(c)

				assert.Equal(t, tt.expectedResult, result)
				if tt.expectedError != nil {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				} else {
					assert.Nil(t, err)
				}
				return nil
			})

			req := httptest.NewRequest("GET", "/?"+buildQuery(tt.queryParams), nil)
			resp, err := app.Test(req)

			assert.Nil(t, err)
			assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		})
	}
}

func buildQuery(params map[string]string) string {
	query := ""
	for key, value := range params {
		if query != "" {
			query += "&"
		}
		query += key + "=" + value
	}
	return query
}
