package app

//
//import (
//	"arxiv/api"
//	postgresql2 "arxiv/internal/repository/postgresql"
//	"arxiv/internal/service"
//)
//
//type App struct {
//	Service *service.Service
//	Api     *api.Rest
//}
//
//func New() (*App, error) {
//	var err error
//	app := &App{}
//
//	conn, err := postgresql2.ConnectToDB("postgresql://root:password@localhost:5432/arxiv?sslmode=disable")
//
//	if err != nil {
//		return nil, err
//	}
//
//	userRepository := postgresql2.NewUserPSQLRepository(conn)
//	papersRepository := postgresql2.NewPaperPSQLRepository(conn)
//
//	app.Service = service.NewService(&service.Config{
//		UserRepository:  userRepository,
//		PaperRepository: papersRepository,
//	})
//
//	app.Api, err = api.New(":8080", app.Service)
//	if err != nil {
//		return nil, err
//	}
//	return app, nil
//}
//
//func (a *App) Run() error {
//	return a.Api.Listen()
//}
//
//func (a *App) Stop() {
//	a.Api.Stop()
//}
