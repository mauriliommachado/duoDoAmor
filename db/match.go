package db

import (
	"log"
)

type Match struct {
	Id      int  `json:"id"`
	MatchId int  `json:"matchId"`
	Status  bool `json:"status"`
}

type Matchs []Match

func (user *Match) Persist() error {
	c := GetDB()
	res, err := c.Exec("INSERT INTO duo.user_match( id, id_match, status) VALUES ($1, $2, $3);", user.Id, user.MatchId, user.Status)
	if err != nil {
		panic(err)
	}
	_, err = res.RowsAffected()
	if err != nil {
		panic(err)
	}
	if user.Status {
		log.Println("Usuário", user.Id, "match com id ", user.MatchId)
	} else {
		log.Println("Usuário", user.Id, "un match com id ", user.MatchId)
	}

	return nil
}

func (user *Match) Merge() error {
	var err error
	//defer dbutil.CloseSession(c)
	//err = c.Update(bson.M{"_id": user.Id}, &user)
	log.Println("Usuário", user.Id, "atualizado")
	if err != nil {
		return err
	}
	return nil
}

func (user *Match) Remove() error {
	var err error
	//defer dbutil.CloseSession(c)
	//err = c.Remove(bson.M{"_id": user.Id})
	log.Println("Usuário", user.Id, "removido")
	if err != nil {
		return err
	}
	return nil
}

func (match *Match) FindById() {
	s := GetDB()
	match.Status = false
	s.QueryRow("SELECT um.status FROM duo.user_match um where um.id = $2 and um.id_match = $1 AND um.status = true;", match.Id, match.MatchId).Scan(&match.Status)
}

func (match *Match) FindNew() (Users, error) {
	s := GetDB()
	var array Users
	rows, err := s.Query("SELECT u.id, u.\"summonerId\", u.name, u.email, r.\"queueType\", r.tier, r.rank, r.\"leaguePoints\", r.wins, r.losses FROM duo.\"user\" u join duo.rank r on r.id = u.\"summonerId\" WHERE u.id <> $1 AND u.id not in(select id_match from duo.user_match um where um.id = $1) order by u.id, r.\"queueType\";", match.Id)
	if err != nil {
		return nil, err
	}
	var user User
	line := 0
	defer rows.Close()
	for rows.Next() {

		err = rows.Scan(&user.Id, &user.SummonerId, &user.Name, &user.Email, &user.Elo[line].QueueType, &user.Elo[line].Tier, &user.Elo[line].Rank, &user.Elo[line].LeaguePoints, &user.Elo[line].Wins, &user.Elo[line].Losses)
		if err != nil {
			return nil, err
		}
		if line == 0 {
			line++
		} else {
			array = append(array, user)
			line = 0
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	for i, item := range array {
		array[i].Champions, err = item.Champions.FindById(item.SummonerId)
		if err != nil {
			return nil, err
		}
	}
	return array, nil
}

func (match *Match) FindAll() (Users, error) {
	s := GetDB()
	var array Users
	rows, err := s.Query("SELECT u.id, u.\"summonerId\", u.name, u.email, r.\"queueType\", r.tier, r.rank, r.\"leaguePoints\", r.wins, r.losses FROM duo.\"user\" u join duo.rank r on r.id = u.\"summonerId\" WHERE u.id <> $1 AND u.id not in(select id_match from duo.user_match um where um.id = $1) order by u.id, r.\"queueType\";", match.Id)
	if err != nil {
		return nil, err
	}
	var user User
	line := 0
	defer rows.Close()
	for rows.Next() {

		err = rows.Scan(&user.Id, &user.SummonerId, &user.Name, &user.Email, &user.Elo[line].QueueType, &user.Elo[line].Tier, &user.Elo[line].Rank, &user.Elo[line].LeaguePoints, &user.Elo[line].Wins, &user.Elo[line].Losses)
		if err != nil {
			return nil, err
		}
		if line == 0 {
			line++
		} else {
			array = append(array, user)
			line = 0
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	for i, item := range array {
		array[i].Champions, err = item.Champions.FindById(item.SummonerId)
		if err != nil {
			return nil, err
		}
	}
	return array, nil
}
