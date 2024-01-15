package app

import (
	"arxivhub/internal/delivery/http/handlers"
	"arxivhub/internal/delivery/http/jwt"
	"arxivhub/internal/delivery/http/server"
	"arxivhub/internal/hubpopulator/vectorencoder"
	"arxivhub/internal/pinecone"
	postgresql "arxivhub/internal/repository/postgresql"
	"arxivhub/internal/service"
	"os"
)

type App struct {
	server *server.Server
}

func NewApp() (*App, error) {
	conn, err := postgresql.ConnectToDB(os.Getenv("PSQL_CONN_STRING"))

	if err != nil {
		return nil, err
	}

	userRepository := postgresql.NewUserPSQLRepository(conn)
	paperRepository := postgresql.NewPaperPSQLRepository(conn)

	userService := service.NewUserService(userRepository)

	encoderClient := vectorencoder.NewClient(os.Getenv("ENCODER_URL"))
	pineconeClient := pinecone.NewClient(os.Getenv("PINECONE_KEY"), os.Getenv("PINECONE_HOST"))

	paperService := service.NewPaperService(paperRepository, encoderClient, pineconeClient)

	jwt, err := jwt.NewJWTManager(os.Getenv("JWT_SECRET"))
	if err != nil {
		return nil, err
	}

	userHandler, err := handlers.NewUserHandler(userService, jwt)
	if err != nil {
		return nil, err
	}

	paperHandler := handlers.NewPaperHandler(paperService, jwt)

	serv, err := server.NewServer(userHandler, paperHandler)
	if err != nil {
		return nil, err
	}

	return &App{server: serv}, nil
}

func (app *App) Run(address string) {
	app.server.Listen(address)
}
