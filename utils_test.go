package yuwiki

import (
	"github.com/dchest/captcha"
	"github.com/matoous/go-nanoid"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestRandomPassword(t *testing.T) {
	random, _ := gonanoid.Nanoid(passwordLen)
	log.Info(random)
	assert.Len(t, random, passwordLen)
}

func TestMatch(t *testing.T) {
	password := GenPassword()
	salt := GenSalt()
	encPassword, _ := EncPassword(password, salt)
	log.Info(password)
	log.Info(salt)
	log.Info(encPassword)
	assert.True(t, Match(password, salt, encPassword))
}

func TestFileSha1(t *testing.T) {
	log.Info(FileSha1("/db/yuwiki.db"))
	fileSha1, _ := FileSha1("/db/yuwiki.db")
	backupSha1, _ := FileSha1("/db/yuwiki.20190504.db")
	assert.Equal(t, fileSha1, backupSha1)
}

func TestToday(t *testing.T) {
	today := Today(DateLayout)
	yesterday := Yesterday(DateLayout)
	log.Info(today, yesterday)
	assert.True(t, strings.Compare(yesterday, today) < 0)
}

func TestByteSize(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < 10; i++ {
		log.Info(ByteSize(uint64(r.Int63n(10000000))))
	}
}

func TestNewCaptcha(t *testing.T) {
	captchaId := captcha.New()
	assert.True(t, captchaId != "")
}
