package server

import (
	"github.com/bmizerany/pat"
	"net/http"
	"log"
	"fmt"
	"encoding/json"
	"github.com/duoDoAmor/db"
	"gopkg.in/mgo.v2/bson"
	"encoding/base64"
	"github.com/rs/cors"
)

type ServerProperties struct {
	Port    string
	Address string
}

func DeleteUser(w http.ResponseWriter, req *http.Request) {
	if !validAuthHeader(req) {
		unauthorized(w)
		return
	}
	var user db.User
	id := req.URL.Query().Get(":id")
	err := user.FindById(db.GetCollectionUsers(), bson.ObjectIdHex(id))
	user.Remove(db.GetCollectionUsers())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	ResponseWithJSON(w, nil, http.StatusNoContent)
}


func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
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
	if len(user.Id.Hex()) > 0 {
		badRequest(w, nil)
		return
	}
	user.Admin = false
	user.Token = base64.StdEncoding.EncodeToString([]byte(user.Email + ":" + user.Pwd))
	err = user.Persist(db.GetCollectionUsers())
	if err != nil {
		badRequest(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", req.URL.Path+"/"+user.Id.Hex())
	w.WriteHeader(http.StatusCreated)
}

func UpdateUser(w http.ResponseWriter, req *http.Request) {
	if !validAuthHeader(req) {
		unauthorized(w)
		return
	}
	var user db.User
	var userUp db.User

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&user)
	userUp.FindById(db.GetCollectionUsers(), user.Id)
	if len(userUp.Id.Hex()) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		badRequest(w, err)
		return
	}
	user.Token = base64.StdEncoding.EncodeToString([]byte(user.Email + ":" + user.Pwd))
	err = user.Merge(db.GetCollectionUsers())
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
	users, err := users.FindAll(db.GetCollectionUsers())
	if err != nil {
		badRequest(w, err)
		return
	}
	for i,_ := range users{
		users[i].Pwd=""
	}
	resp, _ := json.Marshal(users)
	ResponseWithJSON(w, resp, http.StatusOK)
}
func validAuthHeader(req *http.Request) bool {
	auth := req.Header.Get("Authorization")
	if len(auth) <= 6 {
		return false
	}
	var user db.User
	user.Token = auth[6:]
	if user.FindHash(db.GetCollectionUsers()){
		return true
	}else{
		return false
	}
}

func FindById(w http.ResponseWriter, req *http.Request) {
	var user db.User
	id := req.URL.Query().Get(":id")
	err := user.FindById(db.GetCollectionUsers(), bson.ObjectIdHex(id))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resp, _ := json.Marshal(user)
	ResponseWithJSON(w, resp, http.StatusOK)
}

func Validate(w http.ResponseWriter, req *http.Request) {
	var user db.User
	hash := req.URL.Query().Get(":hash")
	user.Token = hash
	if user.FindHash(db.GetCollectionUsers()) {
		resp, _ := json.Marshal(user)
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

func StartUsers(properties ServerProperties) {
	m := pat.New()
	handler := cors.AllowAll().Handler(m)
	mapEndpoints(*m, properties)
	http.Handle("/", handler)
	fmt.Println("servidor iniciado no endereço localhost:" + properties.Port + properties.Address)
	err := http.ListenAndServe(":"+properties.Port, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
func mapEndpoints(m pat.PatternServeMux, properties ServerProperties) {
	m.Post(properties.Address, http.HandlerFunc(InsertUser))
	m.Put(properties.Address, http.HandlerFunc(UpdateUser))
	m.Del(properties.Address+"/:id", http.HandlerFunc(DeleteUser))
	m.Get(properties.Address, http.HandlerFunc(FindAllUsers))
	m.Get(properties.Address+"/:id", http.HandlerFunc(FindById))
	m.Get(properties.Address+"/validate/:hash", http.HandlerFunc(Validate))
}


