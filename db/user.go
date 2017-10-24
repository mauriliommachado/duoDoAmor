package db

import (
	"log"
	"strconv"
	"database/sql"
)

type User struct {
	Id         int    `json:"id,omitempty"`
	SummonerId int    `json:"summonerId,omitempty"`
	Name       string `json:"name,omitempty"`
	Email      string `json:"email,omitempty"`
	Pwd        string `json:"pwd,omitempty"`
	Token      string `json:"token,omitempty"`
	Admin      bool   `json:"admin,omitempty"`
	Elo        Elos   `json:"elo,omitempty"`
}

type Users []User

func (user *User) Persist() error {
	c := GetDB()
	err := c.QueryRow("INSERT INTO duo.\"user\"( \"summonerId\", name, email, pwd, token, admin) VALUES (" + strconv.Itoa(user.SummonerId) + ",'" + user.Name + "', '" + user.Email + "', '" + user.Pwd + "', '" + user.Token + "', false) RETURNING id;").Scan(&user.Id)
	if err != nil {
		return err
	}
	log.Println("Usuário", user.Name, "inserido com id "+strconv.Itoa(user.Id))
	return nil
}

func (user *User) Merge() error {
	c := GetDB()
	err := c.QueryRow("update duo.\"user\" set name = '" + user.Name + "', email = '" + user.Email + "', pwd = '" + user.Pwd + "', token = '" + user.Token + "' where id = " + strconv.Itoa(user.Id) + " RETURNING id;").Scan(&user.Id)
	if err != nil {
		return err
	}
	log.Println("Usuário", user.Name, "inserido com id "+strconv.Itoa(user.Id))
	return nil
}

func (user *User) Remove() error {
	var err error
	//defer dbutil.CloseSession(c)
	//err = c.Remove(bson.M{"_id": user.Id})
	log.Println("Usuário", user.Name, "removido")
	if err != nil {
		return err
	}
	return nil
}

func (user *User) FindById(id int) error {
	s := GetDB()
	log.Println(user)
	row := s.QueryRow("SELECT id, token, name, admin, email FROM duo.\"user\" WHERE id = $1", user.Id)
	err := row.Scan(&user.Id, &user.Token, &user.Name, &user.Admin, &user.Email)
	if err == sql.ErrNoRows {
		log.Println(err)
		log.Println(user.Token)
		return err
	} else if err != nil {
		log.Println(err)
		log.Println(user.Token)
		return err
	}
	return nil
}

func (user *User) FindLogin() bool {
	//defer dbutil.CloseSession(c)
	//err := c.Find(bson.M{"email": user.Email,"pwd":user.Pwd}).One(&user)
	//if err != nil {
	//	return false
	//}
	return true
}

func (user *User) FindHash() bool {
	s := GetDB()
	row := s.QueryRow("SELECT id, token, name, admin, email FROM duo.\"user\" WHERE token = $1", user.Token)
	err := row.Scan(&user.Id, &user.Token, &user.Name, &user.Admin, &user.Email)
	if err == sql.ErrNoRows {
		log.Println(err)
		log.Println(user.Token)
		return false
	} else if err != nil {
		log.Println(err)
		log.Println(user.Token)
		return false
	}
	return true
}

func (users Users) FindAll() (Users, error) {
	//defer dbutil.CloseSession(c)
	//err := c.Find(bson.M{}).All(&users)
	//if err != nil {
	//	return users, err
	//}
	return users, nil
}
