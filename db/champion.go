package db

import (
	"log"
	"strconv"
)

type Champion struct {
	SummonerId     int   `json:"summonerId,omitempty"`
	ChampionLevel  int   `json:"championLevel,omitempty"`
	ChampionPoints int   `json:"championPoints,omitempty"`
	ChampionId     int   `json:"championId,omitempty"`
	LastPlayTime   int64 `json:"lastPlayTime,omitempty"`
}

type Champions []Champion

func (champions Champions) Persist() error {
	c := GetDB()
	for _, champion := range champions {
		err := c.QueryRow("INSERT INTO duo.champion(\"summonerId\", level, points, id, \"ultimoGame\") VALUES (" + strconv.Itoa(champion.SummonerId) + ", " + strconv.Itoa(champion.ChampionLevel) + ", " + strconv.Itoa(champion.ChampionPoints) + ", " + strconv.Itoa(champion.ChampionId) + ", " + strconv.FormatInt(champion.LastPlayTime,10) + ") RETURNING id;").Scan(&champion.SummonerId)
		if err != nil {
			return err
		}
		log.Println("Champion", champion.ChampionId, "inserido para "+strconv.Itoa(champion.SummonerId))
	}
	return nil
}

func (champions Champions) FindById(id int) (Champions, error) {
	s := GetDB()
	rows,err := s.Query("SELECT \"summonerId\", level, points, id, \"ultimoGame\" FROM duo.champion c where c.\"summonerId\" = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c Champion
		err := rows.Scan(&c.SummonerId, &c.ChampionLevel, &c.ChampionPoints, &c.LastPlayTime)
		if err != nil {
			return nil, err
		}
		champions = append(champions, c)
	}
	return champions,nil
}

func (champion *Champion) Merge() error {
	/*c := GetDB()
	err := c.QueryRow("update duo.\"user\" set name = '" + user.Name + "', email = '" + user.Email + "', pwd = '" + user.Pwd + "', token = '" + user.Token + "', discord = '" + user.Discord + "' where id = " + strconv.Itoa(user.Id) + " RETURNING id;").Scan(&user.Id)
	if err != nil {
		return err
	}
	log.Println("Usuário", user.Name, "alterado com id "+strconv.Itoa(user.Id))*/
	return nil
}

func (champion *Champion) Remove() error {
	var err error
	//defer dbutil.CloseSession(c)
	//err = c.Remove(bson.M{"_id": user.Id})
	//log.Println("Usuário", user.Name, "removido")
	if err != nil {
		return err
	}
	return nil
}