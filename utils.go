package yuwiki

import (
	"errors"
	"github.com/matoous/go-nanoid"
	"golang.org/x/crypto/bcrypt"
	"log"
)

const (
	passwordLen = 12
	saltLen     = 6
)

func GenSalt() string {
	salt, _ := gonanoid.Nanoid(saltLen)
	return salt
}

func GenPassword() string {
	password, _ := gonanoid.Nanoid(passwordLen)
	return password
}

func EncPassword(raw string, salt string) (string, error) {
	if len(raw) < 6 {
		return "", errors.New("密码长度必须不小于6位")
	}
	if encHash, err := bcrypt.GenerateFromPassword([]byte(raw+salt), bcrypt.DefaultCost); err != nil {
		log.Fatal("生成密码密文失败", err)
		return "", err
	} else {
		return string(encHash), nil
	}
}

func Match(raw string, salt string, enc string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(enc), []byte(raw+salt))
	return err == nil
}
