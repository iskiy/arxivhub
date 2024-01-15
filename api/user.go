package api

//
//import (
//	"arxiv/internal/models"
//	"context"
//	"database/sql"
//	"encoding/json"
//	"errors"
//	"io"
//	"net/http"
//	"time"
//)
//
//func (rest *Rest) register(w http.ResponseWriter, r *http.Request) {
//	body := io.LimitReader(r.Body, 1024)
//	defer r.Body.Close()
//
//	var params models.RegisterUserRequest
//
//	err := json.NewDecoder(body).Decode(&params)
//
//	if err != nil {
//		rest.sendError(w, http.StatusBadRequest, errors.New("invalid Data Format"))
//		return
//	}
//
//	err = rest.validator.Struct(params)
//	if err != nil {
//		rest.sendError(w, http.StatusBadRequest, errors.New("invalid Data Format"))
//		return
//	}
//
//	user, err := rest.service.Users.RegisterUser(r.Context(), params)
//	if err != nil {
//		rest.sendError(w, http.StatusInternalServerError, errors.New("internal error"))
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusCreated)
//	rest.sendData(w, models.NewRegisterResponse(user))
//}
//
//func (rest *Rest) login(w http.ResponseWriter, r *http.Request) {
//	body := io.LimitReader(r.Body, 1024)
//	defer r.Body.Close()
//
//	var params models.LoginUserRequest
//
//	err := json.NewDecoder(body).Decode(&params)
//
//	if err != nil {
//		rest.sendError(w, http.StatusBadRequest, errors.New("json decode error"))
//		return
//	}
//
//	err = rest.validator.Struct(params)
//	if err != nil {
//		rest.sendError(w, http.StatusBadRequest, errors.New("validator error"))
//		return
//	}
//
//	user, err := rest.service.Users.GetUserByUsername(context.Background(), params.Username)
//	if err != nil {
//		if err == sql.ErrNoRows {
//			rest.sendError(w, http.StatusNotFound, errors.New("no user error"))
//			return
//		}
//		rest.sendError(w, http.StatusNotFound, errors.New("server error"))
//		return
//	}
//
//	err = rest.service.Users.CheckPassword(user.HashedPassword, params.Password)
//	if err != nil {
//		rest.sendError(w, http.StatusUnauthorized, errors.New("check password error"))
//	}
//
//	token, _, err := rest.tokenMaker.CreateToken(user.ID, user.Username, time.Minute*15)
//	if err != nil {
//		rest.sendError(w, http.StatusNotFound, errors.New("create token error"))
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusCreated)
//	rest.sendData(w, models.LoginUserResponse{ID: user.ID, Username: user.Username, Token: token})
//}
