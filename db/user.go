package db

import (
	"log"
	"strconv"
)

type User struct {
	Id    int `json:"id" bson:"_id,omitempty"`
	SummonerId int `json:"summonerId"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
	Token string `json:"token"`
	Admin bool
}


type Users []User


func (user *User) Persist() error {
	c := GetDB()
	err := c.QueryRow("INSERT INTO duo.\"user\"( \"summonerId\", name, email, pwd, token, admin) VALUES (" + strconv.Itoa(user.SummonerId) + ",'" + user.Name + "', '"+ user.Email +"', '"+ user.Pwd +"', '"+ user.Token +"', false) RETURNING id;").Scan(&user.Id)
	if err != nil {
		return err
	}
	log.Println("Usuário", user.Name, "inserido")
	return nil
}

func (user *User) Merge() error {
	var err error
	//defer dbutil.CloseSession(c)
	//err = c.Update(bson.M{"_id": user.Id}, &user)
	log.Println("Usuário", user.Name, "atualizado")
	if err != nil {
		return err
	}
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
	//defer dbutil.CloseSession(c)
	//err := c.Find(bson.M{"_id": id}).One(&user)
	//if err != nil {
	//	return err
	//}
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
	//defer dbutil.CloseSession(c)
	//err := c.Find(bson.M{"token": user.Token}).One(&user)
	//if err != nil {
	//	log.Println(err,user.Token)
	//	return false
	//}
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
