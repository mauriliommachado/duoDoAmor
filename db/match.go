package db

import (
	"log"
)

type Match struct {
	Id    int `json:"id"`
	MatchId int `json:"matchId"`
	Status bool `json:"status"`
}


type Matchs []Match


func (user *Match) Persist() 	error {
	c := GetDB()
	res, err := c.Exec("INSERT INTO duo.user_match( id, id_match, status) VALUES ($1, $2, $3);",user.Id, user.MatchId, user.Status)
	if err != nil {
		panic(err)
	}
	_ , err = res.RowsAffected()
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
	rows, err := s.Query("SELECT id, \"summonerId\", name, email FROM duo.\"user\" WHERE id <> $1 AND id not in(select id_match from duo.user_match um where um.id = $1);", match.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id,&user.SummonerId,&user.Name,&user.Email)
		if err != nil {
			return nil, err
		}
		array = append(array, user)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return array, nil
}


func (users Matchs) FindAll() (Matchs, error) {
	//defer dbutil.CloseSession(c)
	//err := c.Find(bson.M{}).All(&users)
	//if err != nil {
	//	return users, err
	//}
	return users, nil
}
