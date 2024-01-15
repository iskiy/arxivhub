package api

//
//import (
//	"log"
//	"net/http"
//)
//
//func (rest *Rest) payload(w http.ResponseWriter, r *http.Request) {
//	payload := r.Context().Value("payload")
//	log.Println(payload)
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusCreated)
//	rest.sendData(w, payload)
//}
