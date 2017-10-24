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

func DeleteUser(w http.ResponseWriter, req *http.Request) {
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

func InsertUser(w http.ResponseWriter, req *http.Request) {
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

	user.Admin = false
	user.Token = base64.StdEncoding.EncodeToString([]byte(strings.ToLower(user.Name) + ":" + user.Pwd))
	summoner, err := db.FindByName(user.Name)
	if err != nil {
		badRequest(w, err)
		return
	}
	user.SummonerId = summoner.Id
	err = summoner.Persist()
	if err != nil {
		badRequest(w, err)
		return
	}

	elos, err := db.FindEloById(summoner.Id)
	if err != nil {
		badRequest(w, err)
		return
	}

	err = elos.Persist()
	if err != nil {
		badRequest(w, err)
		return
	}

	err = user.Persist()
	if err != nil {
		badRequest(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", req.URL.Path+"/"+strconv.Itoa(user.Id))
	w.WriteHeader(http.StatusCreated)
}

func UpdateUser(w http.ResponseWriter, req *http.Request) {
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

	user.Admin = false
	user.Token = base64.StdEncoding.EncodeToString([]byte(strings.ToLower(user.Name) + ":" + user.Pwd))
	user.FindById(user.Id)
	if err != nil {
		badRequest(w, err)
		return
	}

	err = user.Merge()
	if err != nil {
		badRequest(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func FindAllUsers(w http.ResponseWriter, req *http.Request) {
	if !validAuthHeader(req) {
		unauthorized(w)
		return
	}
	var users db.Users
	users, err := users.FindAll()
	if err != nil {
		badRequest(w, err)
		return
	}
	for i, _ := range users {
		users[i].Pwd = ""
	}
	resp, _ := json.Marshal(users)
	ResponseWithJSON(w, resp, http.StatusOK)
}

func FindById(w http.ResponseWriter, req *http.Request) {
	var user db.User
	id, err := strconv.Atoi(req.URL.Query().Get(":id"))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = user.FindById(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resp, _ := json.Marshal(user)
	ResponseWithJSON(w, resp, http.StatusOK)
}

func MapEndpointsUsers(m pat.PatternServeMux, properties ServerProperties) {
	m.Post(properties.Address, http.HandlerFunc(InsertUser))
	m.Put(properties.Address, http.HandlerFunc(UpdateUser))
	m.Del(properties.Address+"/:id", http.HandlerFunc(DeleteUser))
	m.Get(properties.Address, http.HandlerFunc(FindAllUsers))
	m.Get(properties.Address+"/:id", http.HandlerFunc(FindById))
	m.Get(properties.Address+"/validate/:hash", http.HandlerFunc(Validate))
}
