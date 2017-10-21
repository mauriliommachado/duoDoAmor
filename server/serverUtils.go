package server

import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/duoDoAmor/db"
)

type ServerProperties struct {
	Port    string
	Address string
}


func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

func validAuthHeader(req *http.Request) bool {
	auth := req.Header.Get("Authorization")
	if len(auth) <= 6 {
		return false
	}
	var user db.User
	user.Token = auth[6:]
	if user.FindHash(){
		return true
	}else{
		return false
	}
}

func Validate(w http.ResponseWriter, req *http.Request) {
	var user db.User
	hash := req.URL.Query().Get(":hash")
	user.Token = hash
	if user.FindHash() {
		resp, _ := json.Marshal(user)
		log.Println(user.Name + " Autenticado")
		ResponseWithJSON(w, resp, http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
}

func badRequest(w http.ResponseWriter, err error) {
	log.Println(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
}
