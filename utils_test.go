package yuwiki

import (
	"github.com/matoous/go-nanoid"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestRandomPassword(t *testing.T) {
	random, _ := gonanoid.Nanoid(passwordLen)
	log.Print(random)
	assert.Len(t, random, passwordLen)
}

func TestMatch(t *testing.T) {
	password := GenPassword()
	salt := GenSalt()
	encPassword, _ := EncPassword(password, salt)
	log.Println(password)
	log.Println(salt)
	log.Println(encPassword)
	assert.True(t, Match(password, salt, encPassword))
}
