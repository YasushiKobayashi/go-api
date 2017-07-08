package main

import (
	"app/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ユーザー登録成功(t *testing.T) {
	assert := assert.New(t)
	save := models.User{
		Name:     "hoge",
		Email:    models.NewNullString("hoge@hogecom"),
		Password: models.NewNullString("hogehoge"),
	}
	data, err := models.CreateUser(save)
	if err != nil {
		t.Log(err)
		t.Fatalf("err")
	}
	assert.NotEqual("", data.Token, "トークンが取得できている")
}

func Test_ユーザーログイン成功(t *testing.T) {
	assert := assert.New(t)
	save := models.User{
		Email:    "hoge@hogecom",
		Password: "hogehoge",
	}
	data, err := models.Login(save)
	if err != nil {
		t.Log(err)
		t.Fatalf("err")
	}
	assert.NotEqual("", data.Token, "トークンが取得できている")
}
