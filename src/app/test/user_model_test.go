package main

import (
	"app/models"
	"testing"
)

func Test_ユーザー登録成功(t *testing.T) {
	save := models.User{
		Name:     "hoge",
		Email:    models.NewNullString("hogeahogecom"),
		Password: models.NewNullString("hogehoge"),
	}
	data, err := models.CreateUser(save)
	if err != nil {
		t.Log(err)
		t.Fatalf("err")
	}
	if data.Token != "" {
		t.Log(data)
	}
}
