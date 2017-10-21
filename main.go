package main

import (
	"os"
	"github.com/duoDoAmor/server"
	"github.com/duoDoAmor/db"
	"fmt"
	"log"
	"github.com/rs/cors"
	"github.com/bmizerany/pat"
	"net/http"
)

func main() {
	startDb()
	m := pat.New()
	handler := cors.AllowAll().Handler(m)
	server.MapEndpointsUsers(*m, server.ServerProperties{Address: "/api/users", Port: determineListenAddress()})
	server.MapEndpointsMatchs(*m, server.ServerProperties{Address: "/api/match", Port: determineListenAddress()})
	http.Handle("/", handler)
	fmt.Println("servidor iniciado no endereço localhost:" + determineListenAddress())
	err := http.ListenAndServe(":"+determineListenAddress(), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func startDb() {
	db.Start()
}

func determineListenAddress() (string) {
	port := os.Getenv("PORT")
	if port == "" {
		return "8080"
	}
	return port
}
