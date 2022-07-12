package controllers

import (
	"fmt"
	"mas-kusa-api/db"
)

func AddUser(acct string, instance string, token string) error {
	var u db.User
	if IsAlreadyExistUser(u.Instance, u.Token) {
		return fmt.Errorf("this user is already exist")
	}
	u = db.User{
		Instance: instance,
		Name:     acct,
		Token:    token,
	}
	res := db.Psql.Create(&u)
	return res.Error
}

func IsAlreadyExistUser(instance string, token string) bool {
	err := db.Psql.Where("instance = ?", instance).Where("token = ?", token).First(&db.User{}).Error
	return err == nil
}

func GetUserInfo(token string) (*db.User, error) {
	u := db.User{}
	res := db.Psql.Model(&db.User{}).First(&u)
	if res.Error != nil {
		return nil, res.Error
	}
	return &u, nil
}
