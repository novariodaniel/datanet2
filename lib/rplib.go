package lib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"

	log "projects/datanet2/logging"
)

/*
	Rplib (Raditya Pratama Library) for Golang
	Library created by Raditya Pratama
	for everything
*/

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func EncryptText(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func DecryptText(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func encryptFile(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(EncryptText(data, passphrase))
}

func decryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return DecryptText(data, passphrase)
}

func IsMapKeyExists(arr interface{}, key string, required ...bool) (code int, isExist bool, value interface{}) {
	isExist = false
	code = 1
	value = nil
	isRequired := true
	if !Isset(required, 0) {
		isRequired = false
	} else {
		isRequired = true
	}

	if newArr, ok := arr.(map[string]string); ok {
		if value, isExist = newArr[key]; isExist {
			return
		}
	} else if newArr, ok := arr.(map[string]int); ok {
		if value, isExist = newArr[key]; isExist {
			return
		}
	} else if newArr, ok := arr.(map[string]int64); ok {
		if value, isExist = newArr[key]; isExist {
			return
		}
	} else if newArr, ok := arr.(map[string]bool); ok {
		if value, isExist = newArr[key]; isExist {
			return
		}
	} else if newArr, ok := arr.(map[string]interface{}); ok {
		if value, isExist = newArr[key]; isExist {
			return
		}
	} else if newArr, ok := arr.(map[string]float32); ok {
		if value, isExist = newArr[key]; isExist {
			return
		}
	}

	if isRequired && !isExist {
		code = -1
	} else if !isExist {
		code = 0
	}
	return
}

/*
	Example of using StringToDate Function
	Created by Raditya Pratama
	ex :
		GetDate()
		return : current timestamp with format YYYY-mm-dd HH:ii:ss

		GetDate("yyyy-mm-dd") || StringToDate("y-m-d")
		return : current timestamp with format yyyy-mm-dd

		GetDate("ymd", "20181231")
		return : 20181231
*/

func GetDate(optionParam ...string) (string, time.Time) {
	// var err error
	location, _ := time.LoadLocation("Asia/Jakarta")

	dateTimeToConvert, formatDateTime := time.Now().In(location), "2006-01-02 15:04:05"
	if Isset(optionParam, 0) {
		sTime := optionParam[0]
		formatDateTime = ReplaceFormat(sTime)
	}
	if Isset(optionParam, 1) {
		sTime := optionParam[1]
		layoutParse := getParseLayout(sTime)
		if layoutParse == "" {
			log.Logf("Invalid String Date Format %s", sTime)
			return "", time.Time{}
		}
		dateTimeToConvert, _ = time.ParseInLocation(layoutParse, sTime, location)
	}
	return dateTimeToConvert.Format(formatDateTime), dateTimeToConvert
}

/*
	Example of using StringToDate Function
	Created by Raditya Pratama
	ex :
		GetDate()
		return : current timestamp with format YYYY-mm-dd HH:ii:ss

		GetDate("yyyy-mm-dd") || StringToDate("y-m-d")
		return : current timestamp with format yyyy-mm-dd

		GetDate("ymd", "20181231")
		return : 20181231
*/

func DateAdd(dateParam time.Time, additionalValue ...int) string {
	// var err error
	var newDate string
	addDay, addMonth, addYear := 0, 0, 0
	if Isset(additionalValue, 0) {
		addDay = additionalValue[0]
	}
	if Isset(additionalValue, 1) {
		addMonth = additionalValue[1]
	}
	if Isset(additionalValue, 2) {
		addYear = additionalValue[2]
	}
	addDate := dateParam.AddDate(addYear, addMonth, addDay)
	nDate := addDate.Format("2006-01-02 15:04:05")
	newDate, _ = GetDate("y-m-d h:i:s", nDate)
	return newDate
}

// DateDiff is used to get diff between two string date
// created by Raditya Pratama
func DateDiff(dateParam ...string) map[string]float64 {
	totalParam := len(dateParam)
	if totalParam < 1 || totalParam > 2 {
		log.Logf("Unsufficient Param DateDiff")
		return nil
	}

	defaultFormat := "y-m-d h:i:s"

	strTime1, time1 := GetDate(defaultFormat, dateParam[0])
	strTime2, time2 := "", time.Time{}
	if !Isset(dateParam, 1) {
		strTime2, time2 = GetDate(defaultFormat)
	} else {
		strTime2, time2 = GetDate(defaultFormat, dateParam[1])
	}

	if strTime1 == "" || strTime2 == "" {
		return nil
	}
	duration := time2.Sub(time1)

	var diffResult = map[string]float64{
		"seconds": duration.Seconds(),
		"minutes": duration.Minutes(),
		"hours":   duration.Hours(),
	}

	diffResult["days"] = diffResult["hours"] / 24
	diffResult["mounths"] = diffResult["days"] / 30
	diffResult["years"] = diffResult["mounths"] / 12
	/*var diff float64
	if reqDuration == "all" {

	}else if reqDuration == "s" {
		diff = duration.Seconds()
	} else if reqDuration == "i" {
		diff = duration.Minutes()
	} else if reqDuration == "h" {
		diff = duration.Hours()
	} else if reqDuration == "d" {
		diff = (duration.Hours()/24)
	} else if reqDuration == "m" {
		diff = (duration.Hours()/24/30)
	}
	else if reqDuration == "y" {
		diff = (duration.Hours()/24/30/12)
	}*/

	return diffResult
}

func getParseLayout(dateString string) (layoutParse string) {
	// if time.Parse("2016-01-02", dateString)
	var err error
	var formatLists = []string{
		"2006-01-02 15:04:05",
		"2006-02-01 15:04:05",
		"02-01-2006 15:04:05",
		"01-02-2006 15:04:05",
		"2006-01-02 15:04",
		"2006-02-01 15:04",
		"02-01-2006 15:04",
		"01-02-2006 15:04",

		"2006/01/02 15:04:05",
		"2006/02/01 15:04:05",
		"02/01/2006 15:04:05",
		"01/02/2006 15:04:05",
		"2006/01/02 15:04",
		"2006/02/01 15:04",
		"02/01/2006 15:04",
		"01/02/2006 15:04",

		"02012006",
		"01022006",
		"20060102",
		"20060201",

		"02-01-2006",
		"01-02-2006",
		"2006-01-02",
		"2006-02-01",

		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
		"2006/02/01",

		"02012006 150405",
		"01022006 150405",
		"20060102 150405",
		"20060201 150405",

		"02/01/2006 150405",
		"01/02/2006 150405",
		"2006/01/02 150405",
		"2006/02/01 150405",

		"02012006 1504",
		"01022006 1504",
		"20060102 1504",
		"20060201 1504",

		"02/01/2006 1504",
		"01/02/2006 1504",
		"2006/01/02 1504",
		"2006/02/01 1504",
	}

	for _, format := range formatLists {
		_, err = time.Parse(format, dateString)
		if err == nil {
			layoutParse = format
			break
		}
	}
	return
}

func ReplaceByArr(old []string, new string, source string) string {
	for _, e := range old {
		source = strings.Replace(source, e, new, -1)
	}
	return source
}

func ReplaceFormat(dateFormat string) (formatDateTime string) {
	sFormat := dateFormat
	sFormat = ReplaceByArr([]string{"YYYY", "yyyy", "Y", "y"}, "2006", sFormat)
	sFormat = ReplaceByArr([]string{"MM", "mm", "M", "m"}, "01", sFormat)
	sFormat = ReplaceByArr([]string{"DD", "dd", "D", "d"}, "02", sFormat)
	sFormat = ReplaceByArr([]string{"HH", "hh", "H", "h"}, "15", sFormat)
	sFormat = ReplaceByArr([]string{"II", "ii", "I", "i"}, "04", sFormat)
	sFormat = ReplaceByArr([]string{"SS", "ss", "S", "s"}, "05", sFormat)
	sFormat = ReplaceByArr([]string{"W"}, "Monday", sFormat)
	sFormat = ReplaceByArr([]string{"w"}, "Mon", sFormat)
	sFormat = ReplaceByArr([]string{"F"}, "January", sFormat)
	sFormat = ReplaceByArr([]string{"f"}, "Jan", sFormat)

	formatDateTime = sFormat
	return
}

func DeFormatDate(dateFormat string) (formatDateTime string, reformatDate string) {
	// log.Logf("%s", dateFormat)
	reformatDate = getParseLayout(dateFormat)
	// dateTimeToConvert, _ := time.Parse(reformatDate, dateFormat)
	sFormat := reformatDate

	sFormat = strings.Replace(sFormat, "2006", "yyyy", -1)
	sFormat = strings.Replace(sFormat, "01", "mm", -1)
	sFormat = strings.Replace(sFormat, "02", "dd", -1)
	sFormat = strings.Replace(sFormat, "15", "hh", -1)
	sFormat = strings.Replace(sFormat, "04", "ii", -1)
	sFormat = strings.Replace(sFormat, "05", "ss", -1)

	formatDateTime = sFormat
	// log.Logf("%s", dateFormat, reformatDate, sFormat, dateTimeToConvert)
	return
}

/*
	example of using in_array
	names := []string{"Mary", "Anna", "Beth", "Johnny", "Beth"}
    fmt.Println(in_array("Anna", names)) // results true, 1
    fmt.Println(in_array("Jon", names)) // results false, -1
*/
func in_array(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func Issets(arr interface{}, index int) bool {
	if _, ok := arr.([]string); ok {
		return (len(arr.([]string)) > index)
	} else if _, ok := arr.([]int); ok {
		return (len(arr.([]int)) > index)
	} else if _, ok := arr.([]interface{}); ok {
		return (len(arr.([]interface{})) > index)
	} else if _, ok := arr.([]bool); ok {
		return (len(arr.([]bool)) > index)
	} else if _, ok := arr.([]map[string]interface{}); ok {
		return (len(arr.([]map[string]interface{})) > index)
	} else if _, ok := arr.([]map[string]bool); ok {
		return (len(arr.([]map[string]bool)) > index)
	} else if _, ok := arr.([]map[string]string); ok {
		return (len(arr.([]map[string]string)) > index)
	} else if _, ok := arr.([]map[string]int); ok {
		return (len(arr.([]map[string]int)) > index)
	} else if _, ok := arr.([]map[int]interface{}); ok {
		return (len(arr.([]map[int]interface{})) > index)
	} else if _, ok := arr.([]map[int]bool); ok {
		return (len(arr.([]map[int]bool)) > index)
	} else if _, ok := arr.([]map[int]string); ok {
		return (len(arr.([]map[int]string)) > index)
	} else if _, ok := arr.([]map[int]int); ok {
		return (len(arr.([]map[int]int)) > index)
	}
	return false
}
