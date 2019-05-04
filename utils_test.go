package yuwiki

import (
	"github.com/matoous/go-nanoid"
	"github.com/stretchr/testify/assert"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"
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

func TestFileSha1(t *testing.T) {
	log.Println(FileSha1("/db/yuwiki.db"))
	fileSha1, _ := FileSha1("/db/yuwiki.db")
	backupSha1, _ := FileSha1("/db/yuwiki.20190504.db")
	assert.Equal(t, fileSha1, backupSha1)
}

func TestToday(t *testing.T) {
	today := Today(DateLayout)
	yesterday := Yesterday(DateLayout)
	log.Println(today, yesterday)
	assert.True(t, strings.Compare(yesterday, today) < 0)
}

func TestByteSize(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 10; i++ {
		log.Println(ByteSize(uint64(r.Int63n(10000000))))
	}
}
