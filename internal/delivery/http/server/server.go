package server

import (
	"arxivhub/internal/delivery/http/handlers"
	"arxivhub/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Server struct {
	app *fiber.App
}

func NewServer(user *handlers.UserHandler, paper *handlers.PaperHandler) (*Server, error) {
	app := fiber.New()

	app.Use(cors.New())

	app.Post("/register", user.RegisterUser)
	app.Post("/login", user.Login)
	app.Post("/save", middleware.AuthMiddleware(user.Manager), paper.SavePaper)
	app.Delete("/unsave", middleware.AuthMiddleware(user.Manager), paper.DeleteSavedPaper)
	app.Get("/search", paper.SearchPaper)
	app.Get("/savedpapers", middleware.AuthMiddleware(user.Manager), paper.GetSavedPapers)

	err := app.Listen(":8080")
	if err != nil {
		return nil, err
	}

	return &Server{app: app}, nil
}

func (s *Server) Listen(address string) {
	s.app.Listen(address)
}
