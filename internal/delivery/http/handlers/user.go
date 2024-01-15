package handlers

import (
	"arxivhub/internal/delivery/http/jwt"
	"arxivhub/internal/models"
	"arxivhub/internal/service"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"time"
)

var (
	ErrAlreadyRegistered  = errors.New("user already registered")
	ErrRegisterUser       = errors.New("register error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCreateJWTError     = errors.New("failed to create a token")
)

type UserHandler struct {
	userService service.UserService
	validator   *validator.Validate
	Manager     *jwt.JWTManager
}

func NewUserHandler(userService service.UserService, jwtManager *jwt.JWTManager) (*UserHandler, error) {
	apiValidator := validator.New()

	return &UserHandler{
		validator:   apiValidator,
		Manager:     jwtManager,
		userService: userService,
	}, nil
}

func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
	var requestParams models.RegisterUserRequest
	err := parseBody(c, &requestParams, h.validator)
	if err != nil {
		return err
	}

	user, err := h.userService.RegisterUser(c.UserContext(), requestParams)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return fiber.NewError(fiber.StatusBadRequest, ErrAlreadyRegistered.Error())
			}
		}
		return fiber.NewError(fiber.StatusInternalServerError, ErrRegisterUser.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(models.NewRegisterResponse(user))
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var requestParams models.LoginUserRequest
	err := parseBody(c, &requestParams, h.validator)
	if err != nil {
		return err
	}

	user, err := h.userService.LoginUser(c.UserContext(), requestParams)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, ErrInvalidCredentials.Error())
	}

	token, claims, err := h.Manager.CreateToken(user.ID, time.Hour*20)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, ErrCreateJWTError.Error())
	}

	return c.Status(fiber.StatusOK).JSON(models.LoginUserResponse{
		ID:    claims.ID,
		Token: token,
	})
}
