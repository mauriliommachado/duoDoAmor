package client

import (
	"net/http"
	"log"
	"encoding/json"
	"github.com/duoDoAmor/db"
	"strconv"
)

const uriSummonerApi = "https://br1.api.riotgames.com/lol/summoner/v3/summoners/by-name/"
const TOKEN = "RGAPI-13fcb017-0c03-4739-acba-435f5ec5f733"

type Summoner struct {
	Id int	`json:"id"`
	AccountId int `json:"accountId"`
	Name string `json:"name"`
	ProfileIconId int `json:"profileIconId"`
	RevisionDate int64 `json:"revisionDate"`
	SummonerLevel int `json:"summonerLevel"`
}


func FindByName(name string)(Summoner, error) {
	var summoner Summoner

	req, err := http.NewRequest(http.MethodGet, uriSummonerApi + name, nil)
	req.Header.Set("X-Riot-Token", TOKEN)
	if err != nil{
		log.Println(err)
		return summoner, err
	}
	myClient := &http.Client{}
	resp, err := myClient.Do(req)
	if err != nil || resp.Status != "200 OK"{
		log.Println(err)
		return summoner, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&summoner); err != nil {
		log.Println(err)
		return summoner, err
	}
	return summoner, nil
}

func (user *Summoner) Persist() error {
	c := db.GetDB()
	err := c.QueryRow("INSERT INTO duo.summoner(id, \"accId\", \"profileIconId\", \"revisionDate\", level) VALUES (" + strconv.Itoa(user.Id) + ", " + strconv.Itoa(user.AccountId) + ", " + strconv.Itoa(user.ProfileIconId) + ", " + strconv.FormatInt(user.RevisionDate,10)  + ", " + strconv.Itoa(user.SummonerLevel) + ") RETURNING id;").Scan(&user.Id)
	if err != nil {
		return err
	}
	log.Println("Summoner", user.Name, "inserido")
	return nil
}