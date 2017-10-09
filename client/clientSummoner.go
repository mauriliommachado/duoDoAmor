package client

import (
	"net/http"
	"log"
	"encoding/json"
)

const URI_SUMMONER_API = "https://br1.api.riotgames.com/lol/summoner/v3/summoners/by-name/"
const TOKEN = "RGAPI-7bd306d6-5fdf-40c1-8c89-f9f38ddfc412"

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

	req, err := http.NewRequest(http.MethodGet, URI_SUMMONER_API + name, nil)
	req.Header.Set("X-Riot-Token", "Basic "+TOKEN)
	if err != nil {
		log.Println(err)
		return summoner, err
	}
	myClient := &http.Client{}
	resp, err := myClient.Do(req)
	if err != nil {
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
