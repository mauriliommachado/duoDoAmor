package server

import (
	"github.com/bmizerany/pat"
	"net/http"
	"encoding/json"
	"github.com/duoDoAmor/db"
	"encoding/base64"
	"strconv"
	"strings"
)

func DeleteMatch(w http.ResponseWriter, req *http.Request) {
	if !validAuthHeader(req) {
		unauthorized(w)
		return
	}
	var user db.User
	id, err := strconv.Atoi(req.URL.Query().Get(":id"))
	if err != nil {
		badRequest(w, err)
		return
	}
	err = user.FindById(id)
	if err != nil {
		badRequest(w, err)
		return
	}
	user.Remove()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	ResponseWithJSON(w, nil, http.StatusNoContent)
}

func InsertMatch(w http.ResponseWriter, req *http.Request) {
	if !validAuthHeader(req) {
		unauthorized(w)
		return
	}
	var match db.Match
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&match)
	if err != nil {
		badRequest(w, err)
		return
	}
	err = match.Persist()
	if err != nil {
		badRequest(w, err)
		return
	}
	if match.Status {
		match.FindById()
		if err != nil {
			badRequest(w, err)
			return
		}
	}
	resp, _ := json.Marshal(match)
	ResponseWithJSON(w, resp, http.StatusCreated)
}

func UpdateMatch(w http.ResponseWriter, req *http.Request) {
	if !validAuthHeader(req) {
		unauthorized(w)
		return
	}
	var user db.User
	var userUp db.User

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&user)
	userUp.FindById(user.Id)
	if len(string(userUp.Id)) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		badRequest(w, err)
		return
	}
	user.Token = base64.StdEncoding.EncodeToString([]byte(strings.ToLower(user.Name) + ":" + user.Pwd))
	err = user.Merge()
	if err != nil {
		badRequest(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func FindAllMatchs(w http.ResponseWriter, req *http.Request) {
	if !validAuthHeader(req) {
		unauthorized(w)
		return
	}

	var user db.User
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&user)
	if err != nil {
		badRequest(w, err)
		return
	}
	var users db.Matchs
	users, err = users.FindAll()
	if err != nil {
		badRequest(w, err)
		return
	}

	resp, _ := json.Marshal(users)
	ResponseWithJSON(w, resp, http.StatusOK)
}

func FindNewMatchs(w http.ResponseWriter, req *http.Request) {
	if !validAuthHeader(req) {
		unauthorized(w)
		return
	}
	id, err := strconv.Atoi(req.URL.Query().Get(":id"))
	if err != nil {
		badRequest(w, err)
		return
	}
	var match db.Match
	match.Id = id
	users, err := match.FindNew()
	if err != nil {
		badRequest(w, err)
		return
	}
	resp, _ := json.Marshal(users)
	ResponseWithJSON(w, resp, http.StatusOK)
}

func MapEndpointsMatchs(m pat.PatternServeMux, properties ServerProperties) {
	m.Post(properties.Address, http.HandlerFunc(InsertMatch))
	m.Put(properties.Address, http.HandlerFunc(UpdateMatch))
	m.Del(properties.Address+"/:id", http.HandlerFunc(DeleteMatch))
	m.Get(properties.Address, http.HandlerFunc(FindAllMatchs))
	m.Get(properties.Address+"/:id/new", http.HandlerFunc(FindNewMatchs))
}
