package handlers

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrBadDataFormat = errors.New("invalid data format")
	ErrBadData       = errors.New("invalid data")
)

func parseBody[T any](c *fiber.Ctx, requestParams *T, validator *validator.Validate) error {
	err := c.BodyParser(&requestParams)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, ErrBadDataFormat.Error())
	}

	err = validator.Struct(requestParams)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, ErrBadData.Error())
	}

	return nil
}
