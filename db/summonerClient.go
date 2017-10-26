package db

import (
	"net/http"
	"log"
	"encoding/json"
	"strconv"
)

type Summoner struct {
	Id            int    `json:"id"`
	AccountId     int    `json:"accountId"`
	Name          string `json:"name"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}

type Elo struct {
	Id           string `json:"playerOrTeamId"`
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
}

type Elos [2]Elo

const uriSummonerApi = "https://br1.api.riotgames.com/lol/summoner/v3/summoners/by-name/"
const uriEloApi = "https://br1.api.riotgames.com/lol/league/v3/positions/by-summoner/"
const TOKEN = "RGAPI-a617157f-06a1-43f2-8f42-c64c5960b3bd"

func FindByName(name string) (Summoner, error) {
	var summoner Summoner

	req, err := http.NewRequest(http.MethodGet, uriSummonerApi+name, nil)
	req.Header.Set("X-Riot-Token", TOKEN)
	if err != nil {
		log.Println(err)
		return summoner, err
	}
	myClient := &http.Client{}
	resp, err := myClient.Do(req)
	if err != nil || resp.Status != "200 OK" {
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

func FindEloById(id int) (Elos, error) {
	var elo Elos
	req, err := http.NewRequest(http.MethodGet, uriEloApi+strconv.Itoa(id), nil)
	req.Header.Set("X-Riot-Token", TOKEN)
	if err != nil {
		log.Println(err)
		return elo, err
	}
	myClient := &http.Client{}
	resp, err := myClient.Do(req)
	if err != nil || resp.Status != "200 OK" {
		log.Println(err)
		return elo, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&elo);
		err != nil {
		log.Println(err)
		return elo, err
	}
	for i := range elo {
		elo[i].Id = strconv.Itoa(id)
	}
	return elo, nil
}

func (user *Summoner) Persist() error {
	c := GetDB()
	err := c.QueryRow("INSERT INTO duo.summoner(id, \"accId\", \"profileIconId\", \"revisionDate\", level) VALUES (" + strconv.Itoa(user.Id) + ", " + strconv.Itoa(user.AccountId) + ", " + strconv.Itoa(user.ProfileIconId) + ", " + strconv.FormatInt(user.RevisionDate, 10) + ", " + strconv.Itoa(user.SummonerLevel) + ") RETURNING id;").Scan(&user.Id)
	if err != nil {
		return err
	}
	log.Println("Summoner", user.Name, "inserido")
	return nil
}

func (elos Elos) Persist() error {
	c := GetDB()
	for i, elo := range elos {
		id, _ := strconv.ParseInt(elo.Id, 10, 64)
		if elo.Rank == "" {
			if i == 0 && elos[i+1].Rank != "" {
				if elos[i+1].Rank == "RANKED_SOLO_5x5" {
					err := c.QueryRow("INSERT INTO duo.rank(id, \"queueType\", tier, rank, \"leaguePoints\", wins, losses) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;", id, "RANKED_FLEX_SR", "provisional", elo.Rank, elo.LeaguePoints, elo.Wins, elo.Losses).Scan(&elo.Id)
					if err != nil {
						return err
					}
				} else {
					err := c.QueryRow("INSERT INTO duo.rank(id, \"queueType\", tier, rank, \"leaguePoints\", wins, losses) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;", id, "RANKED_SOLO_5x5", "provisional", elo.Rank, elo.LeaguePoints, elo.Wins, elo.Losses).Scan(&elo.Id)
					if err != nil {
						return err
					}
				}
				continue
			} else {
				err := c.QueryRow("INSERT INTO duo.rank(id, \"queueType\", tier, rank, \"leaguePoints\", wins, losses) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;", id, "RANKED_FLEX_SR", "provisional", elo.Rank, elo.LeaguePoints, elo.Wins, elo.Losses).Scan(&elo.Id)
				if err != nil {
					return err
				}
				err = c.QueryRow("INSERT INTO duo.rank(id, \"queueType\", tier, rank, \"leaguePoints\", wins, losses) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;", id, "RANKED_SOLO_5x5", "provisional", elo.Rank, elo.LeaguePoints, elo.Wins, elo.Losses).Scan(&elo.Id)
				if err != nil {
					return err
				}
				log.Println(elo.Id, " inserido elo ", elo.Tier)
				return nil
			}
		}
		err := c.QueryRow("INSERT INTO duo.rank(id, \"queueType\", tier, rank, \"leaguePoints\", wins, losses) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;", id, elo.QueueType, elo.Tier, elo.Rank, elo.LeaguePoints, elo.Wins, elo.Losses).Scan(&elo.Id)
		log.Println(elo.Id, " inserido elo ", elo.Tier)
		if err != nil {
			return err
		}
	}
	return nil
}
