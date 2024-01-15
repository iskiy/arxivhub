package api

//
//import (
//	"arxiv/internal/service"
//	"arxiv/pkg/token"
//	"context"
//	"encoding/json"
//	"github.com/go-playground/validator/v10"
//	"github.com/gorilla/mux"
//	"log"
//	"net"
//	"net/http"
//	"os"
//	"time"
//)
//
//type Rest struct {
//	address    string
//	mux        *mux.Router
//	listener   net.Listener
//	service    *service.Service
//	server     *http.Server
//	validator  *validator.Validate
//	tokenMaker token.Maker
//}
//
//func New(address string, service *service.Service) (*Rest, error) {
//	rest := &Rest{
//		address: address,
//		service: service,
//	}
//
//	api := mux.NewRouter()
//	api.HandleFunc("/users", rest.register).Methods("POST")
//	api.HandleFunc("/login", rest.login).Methods("POST")
//	api.Handle("/payload", rest.auth(rest.payload)).Methods("GET")
//
//	//api.Handle("/events", rest.BasicAuthMiddleware(http.HandlerFunc(rest.addEvent))).Methods("POST")
//
//	rest.mux = api
//	rest.validator = validator.New()
//	tokenMaker, err := token.NewJWTMaker("12312312311231231231123123123112")
//
//	rest.tokenMaker = tokenMaker
//	if err != nil {
//		return nil, err
//	}
//
//	return rest, nil
//}
//
//func (rest *Rest) BasicAuthMiddleware(handler http.HandlerFunc) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		user, pass, ok := r.BasicAuth()
//		err := rest.service.Users.CheckPassword(user, pass)
//		if !ok || err != nil {
//			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password for this site"`)
//			w.WriteHeader(401)
//			w.Write([]byte("Unauthorised.\n"))
//			log.Println(err.Error())
//			return
//		}
//		handler(w, r)
//	}
//}
//
//func (rest *Rest) Listen() (err error) {
//	rest.listener, err = net.Listen("tcp", rest.address)
//	if err != nil {
//		return err
//	}
//
//	r := http.NewServeMux()
//	r.Handle("/", rest.mux)
//	rest.server = &http.Server{
//		Handler: r,
//	}
//
//	rest.setupMiddleware()
//	return rest.server.Serve(rest.listener)
//}
//
//func (rest *Rest) Stop() {
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	if err := rest.server.Shutdown(ctx); err != nil {
//		log.Printf("Could not shut down server correctly: %v\n", err)
//		os.Exit(1)
//	}
//}
//
//func (rest *Rest) setupMiddleware() {
//	rest.mux.Use(func(handler http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			log.Println("new request", r.RequestURI)
//			handler.ServeHTTP(w, r)
//		})
//	})
//}
//
//type Response struct {
//	Status int
//	Data   interface{}
//}
//
//func (rest *Rest) sendError(w http.ResponseWriter, statusCode int, err error) {
//	w.WriteHeader(statusCode)
//
//	bytes, err := json.Marshal(Response{
//		Status: statusCode,
//		Data:   err.Error(),
//	})
//
//	if err != nil {
//		log.Println(err)
//	}
//
//	_, err = w.Write(bytes)
//	if err != nil {
//		log.Println(err)
//	}
//}
//
//func (rest *Rest) sendData(w http.ResponseWriter, data interface{}) {
//	bytes, err := json.Marshal(Response{
//		Status: 1,
//		Data:   data,
//	})
//
//	if err != nil {
//		log.Println(err)
//	}
//
//	_, err = w.Write(bytes)
//	if err != nil {
//		log.Println(err)
//	}
//}
