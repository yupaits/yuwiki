package yuwiki

import (
	"errors"
	"github.com/matoous/go-nanoid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strings"
)

const (
	passwordLen = 12
	saltLen     = 6
)

//生成随机 Salt
func GenSalt() string {
	salt, _ := gonanoid.Nanoid(saltLen)
	return salt
}

//生成随机密码
func GenPassword() string {
	password, _ := gonanoid.Nanoid(passwordLen)
	return password
}

//生成密码密文
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

//校验密码是否匹配
func Match(raw string, salt string, enc string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(enc), []byte(raw+salt))
	return err == nil
}

//创建全路径目录，自动忽略文件
func Mkdirs(path string) bool {
	dir := string(path[:strings.LastIndex(path, "/")])
	if _, err := os.Open(dir); err == nil {
		return true
	}
	err := os.MkdirAll(dir, os.ModePerm)
	return err == nil
}
