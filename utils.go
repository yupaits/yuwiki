package yuwiki

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/matoous/go-nanoid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	passwordLen = 12
	saltLen     = 6
	filenameLen = 8
)

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
)

const (
	DateLayout         = "2006-01-02"
	DateTimeLayout     = "2006-01-02 15:04:05"
	DateTimeLogLayout  = "2006-01-02 15:04:05.000"
	DateTimeFileLayout = "20060102150405"
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
		log.WithField("error", err).Error("生成密码密文失败")
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

//生成文件名
func GenFilename() string {
	filename, _ := gonanoid.Nanoid(filenameLen)
	return filename
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

//获取今天的日期
func Today(layout string) string {
	if layout == "" {
		layout = DateLayout
	}
	return time.Now().Format(layout)
}

//获取昨天的日期
func Yesterday(layout string) string {
	if layout == "" {
		layout = DateLayout
	}
	return time.Now().AddDate(0, 0, -1).Format(layout)
}

//获取当前时间用于生成文件名
func NowFilename(layout string) string {
	if layout == "" {
		layout = DateTimeFileLayout
	}
	return time.Now().Format(layout)
}

//计算文件SHA1值
func FileSha1(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error(err)
		}
	}()
	hash := sha1.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

//复制文件
func CopyFile(dstFile, srcFile string) (written int64, err error) {
	src, err := os.Open(srcFile)
	if err != nil {
		return
	}
	defer func() {
		if err := src.Close(); err != nil {
			log.Error(err)
		}
	}()
	dst, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer func() {
		if err := dst.Close(); err != nil {
			log.Error(err)
		}
	}()
	return io.Copy(dst, src)
}

/** copy from cloudfoundry/bytefmt **/
var invalidByteQuantityError = errors.New("byte quantity must be a positive integer with a unit of measurement like M, MB, MiB, G, GiB, or GB")

// ByteSize returns a human-readable byte string of the form 10M, 12.5K, and so forth.  The following units are available:
//	T: Terabyte
//	G: Gigabyte
//	M: Megabyte
//	K: Kilobyte
//	B: Byte
// The unit that results in the smallest number greater than or equal to 1 is always chosen.
func ByteSize(bytes uint64) string {
	unit := ""
	value := float64(bytes)

	switch {
	case bytes >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0"
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}

// ToMegabytes parses a string formatted by ByteSize as megabytes.
func ToMegabytes(s string) (uint64, error) {
	bytes, err := ToBytes(s)
	if err != nil {
		return 0, err
	}

	return bytes / MEGABYTE, nil
}

// ToBytes parses a string formatted by ByteSize as bytes. Note binary-prefixed and SI prefixed units both mean a base-2 units
// KB = K = KiB	= 1024
// MB = M = MiB = 1024 * K
// GB = G = GiB = 1024 * M
// TB = T = TiB = 1024 * G
func ToBytes(s string) (uint64, error) {
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)

	i := strings.IndexFunc(s, unicode.IsLetter)

	if i == -1 {
		return 0, invalidByteQuantityError
	}

	bytesString, multiple := s[:i], s[i:]
	bytes, err := strconv.ParseFloat(bytesString, 64)
	if err != nil || bytes <= 0 {
		return 0, invalidByteQuantityError
	}

	switch multiple {
	case "T", "TB", "TIB":
		return uint64(bytes * TERABYTE), nil
	case "G", "GB", "GIB":
		return uint64(bytes * GIGABYTE), nil
	case "M", "MB", "MIB":
		return uint64(bytes * MEGABYTE), nil
	case "K", "KB", "KIB":
		return uint64(bytes * KILOBYTE), nil
	case "B":
		return uint64(bytes), nil
	default:
		return 0, invalidByteQuantityError
	}
}

/** copy from cloudfoundry/bytefmt **/
