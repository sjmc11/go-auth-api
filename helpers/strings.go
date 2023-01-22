package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"net/mail"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func SplitStringIdsToIntSlice(stringids string) []int {
	var intids []int
	splitIds := strings.Split(stringids, ",")

	for _, v := range splitIds {
		intid, err := strconv.Atoi(strings.TrimSpace(v))
		if err == nil {
			intids = append(intids, intid)
		}
	}

	sort.Slice(intids, func(i, j int) bool {
		return intids[i] < intids[j]
	})

	return intids
}

func Capitalize(text string) string {
	trimText := strings.TrimSpace(text)
	lowerText := strings.ToLower(trimText)
	firstRune := []rune(lowerText)[0]
	upper := unicode.ToUpper(firstRune)
	upperStr := string(upper)
	result := upperStr + lowerText[1:]
	return result
}

func TokenizeString(text string) string {
	randToken := md5.Sum([]byte(text + strconv.Itoa(int(time.Now().Unix()))))
	return hex.EncodeToString(randToken[:])
}

func StringToHash(pString string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pString), bcrypt.MinCost)

	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numberBytes = letterBytes + "1234567890"

func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RandStringMixed(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = numberBytes[rand.Intn(len(numberBytes))]
	}
	return string(b)
}

func EmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func RoleNiceName(role int) string {
	switch true {
	case role == 1:
		return "Default"
	case role == 2:
		return "Default"
	default:
		return "Unknown role"
	}
}

func SanitizeString(text string) string {
	return strings.TrimSpace(strings.ToLower(text))
}
