package handlers

import (
	"arxivhub/internal/delivery/http/jwt"
	"arxivhub/internal/delivery/http/middleware"
	"arxivhub/internal/models"
	"arxivhub/internal/service"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"strconv"
	"time"
)

var (
	ErrAlreadySaved        = errors.New("user already saved this paper")
	ErrSavePaper           = errors.New("error occurred during saving")
	ErrEmptySearchQuery    = errors.New("empty search query")
	ErrSearchError         = errors.New("error occurred during searching")
	ErrInvalidYear         = errors.New("invalid year format")
	ErrInvalidAmountFormat = errors.New("invalid amount parameter format")
	ErrInvalidAmount       = errors.New("invalid amount parameter, amount should between 1 and 120")
	ErrEmptyPaperID        = errors.New("empty paper_id param")
)

type PaperHandler struct {
	paperService service.PaperService
	Manager      *jwt.JWTManager
	validator    *validator.Validate
}

func NewPaperHandler(paperService service.PaperService, manager *jwt.JWTManager) *PaperHandler {
	apiValidator := validator.New()

	return &PaperHandler{
		Manager:      manager,
		validator:    apiValidator,
		paperService: paperService,
	}
}

func (h *PaperHandler) SearchPaper(c *fiber.Ctx) error {
	params, err := GetRequestSearchParams(c)
	if err != nil {
		return err
	}
	searchResults, err := h.paperService.Search(c.Context(), params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, ErrSearchError.Error())
	}

	err = h.UpdateSavedField(c, searchResults)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.SearchPapersResponse{ScoredPapers: searchResults})
}

func GetRequestSearchParams(c *fiber.Ctx) (models.SearchPaperRequest, error) {
	query := c.Query("query")
	if query == "" {
		return models.SearchPaperRequest{}, fiber.NewError(fiber.StatusBadRequest, ErrEmptySearchQuery.Error())
	}

	queryYear := c.Query("year")
	year := 0

	if queryYear != "" {
		convertedYear, err := strconv.Atoi(queryYear)
		if err != nil {
			return models.SearchPaperRequest{}, fiber.NewError(fiber.StatusBadRequest, ErrInvalidYear.Error())
		}

		if !isValidYear(convertedYear) {
			return models.SearchPaperRequest{}, fiber.NewError(fiber.StatusBadRequest, ErrInvalidYear.Error())
		}

		year = convertedYear
	}

	amount := 0
	queryAmount := c.Query("amount")
	if queryAmount != "" {
		convertedAmount, err := strconv.Atoi(queryAmount)
		if err != nil {
			return models.SearchPaperRequest{}, fiber.NewError(fiber.StatusBadRequest, ErrInvalidAmountFormat.Error())
		}
		if convertedAmount > 120 || convertedAmount < 1 {
			return models.SearchPaperRequest{}, fiber.NewError(fiber.StatusBadRequest, ErrInvalidAmount.Error())
		}
		amount = convertedAmount
	}

	category := c.Query("category")

	return models.SearchPaperRequest{
		Query:    query,
		Category: category,
		Year:     year,
		Amount:   amount,
	}, nil
}

func (h *PaperHandler) UpdateSavedField(c *fiber.Ctx, papers []models.ScoredPaper) error {
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		claims, err := middleware.AuthHeaderValidation(authHeader, h.Manager)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		err = h.paperService.UpdateSavedField(c.Context(), claims.ID, papers)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return nil
}

func (h *PaperHandler) SearchScoredPapers(c *fiber.Ctx) ([]models.ScoredPaper, error) {
	query := c.Query("query")
	if query == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, ErrEmptySearchQuery.Error())
	}

	queryYear := c.Query("year")
	year := 0

	if queryYear != "" {
		convertedYear, err := strconv.Atoi(queryYear)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, ErrInvalidYear.Error())
		}

		if !isValidYear(convertedYear) {
			return nil, fiber.NewError(fiber.StatusBadRequest, ErrInvalidYear.Error())
		}

		year = convertedYear
	}

	amount := 0
	queryAmount := c.Query("amount")
	if queryAmount != "" {
		convertedAmount, err := strconv.Atoi(queryAmount)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, ErrInvalidAmountFormat.Error())
		}
		if convertedAmount > 120 || convertedAmount < 1 {
			return nil, fiber.NewError(fiber.StatusBadRequest, ErrInvalidAmount.Error())
		}
		amount = convertedAmount
	}

	category := c.Query("category")
	searchResults, err := h.paperService.Search(c.Context(), models.SearchPaperRequest{
		Query:    query,
		Category: category,
		Year:     year,
		Amount:   amount,
	})
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, ErrSearchError.Error())
	}

	return searchResults, nil
}

func isValidYear(year int) bool {
	currentYear := time.Now().Year()
	if year <= 0 || year > currentYear {
		return false
	}
	return true
}

func (h *PaperHandler) SavePaper(c *fiber.Ctx) error {
	paperID := c.Query("paper_id")
	claims := c.Locals("claims").(*jwt.Claims)
	params := models.SavePaperForUserParams{
		UserID:  claims.ID,
		PaperID: paperID,
	}

	paper, err := h.paperService.SavePaper(c.Context(), params)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return fiber.NewError(fiber.StatusBadRequest, ErrAlreadySaved.Error())
			}
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(paper)
}

func (h *PaperHandler) DeleteSavedPaper(c *fiber.Ctx) error {
	paperID := c.Query("paper_id")
	if paperID == "" {
		return fiber.NewError(fiber.StatusBadRequest, ErrEmptyPaperID.Error())
	}
	claims := c.Locals("claims").(*jwt.Claims)
	params := models.DeleteSavedPaperParams{
		UserID:  claims.ID,
		PaperID: paperID,
	}

	err := h.paperService.DeleteSavedPaper(c.Context(), params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *PaperHandler) GetSavedPapers(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*jwt.Claims)

	resp, err := h.paperService.GetSavedPaper(c.Context(), claims.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
