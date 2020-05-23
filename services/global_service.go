package services

/*
	Created by Daniel N Pangaribuan
	For Global usage of function
	21 September 2018
*/

import (
	"encoding/json"
	"fmt"
	"net/http"
	lib "projects/datanet2/lib"
	"reflect"
)

func TestGetType(variable ...interface{}) {
	var reflectValue = reflect.ValueOf(variable[0])

	switch reflectValue.Kind() {
	case reflect.Map:
		fmt.Println("variable is Map")
	case reflect.Interface:
		fmt.Println("interface")
	case reflect.Slice:
		fmt.Println("slice")
		/* reflectType := vInterface.Type()
		for i := 0; i < vInterface.NumField(); i++ {
			fmt.Println("nama      :", reflectType.Field(i).Name)
			fmt.Println("tipe data :", reflectType.Field(i).Type)
			fmt.Println("nilai     :", vInterface.Field(i).Interface())
			fmt.Println("")
		} */

	}
}

func DateAddDay(strDate string, additionalDay int) string {
	_, dateInstance := lib.GetDate("y-m-d h:i:s", strDate)
	return lib.DateAdd(dateInstance, 0, 0, additionalDay)
}

func DateAddMonth(strDate string, additionalMonth int) string {
	_, dateInstance := lib.GetDate("y-m-d h:i:s", strDate)
	return lib.DateAdd(dateInstance, 0, additionalMonth, 0)
}

func DateAddYear(strDate string, additionalYear int) string {
	_, dateInstance := lib.GetDate("y-m-d h:i:s", strDate)
	return lib.DateAdd(dateInstance, additionalYear, 0, 0)
}

//summary or detail
// func GetFileFTP(types string, db lib.DbConnection) (errCode int, errStr string) {
// 	var PathString string
// 	var idxDate, idxTime int
// 	idxDate = 1
// 	PathString = "/interface/itf/"
// 	destPath := dt.FilesRootFolder + types + "/#tgl/" //required

// 	if types == "summary" {
// 		PathString += "comlog/sodogi_copy"
// 		idxDate = 1
// 	} else if types == "dso" {
// 		PathString += "zsmart/USR004_copy"
// 		idxDate = 2
// 	} else if types == "gto" {
// 		PathString += "zsmart/USR005_copy"
// 		idxDate = 2
// 	} else if types == "cso" {
// 		PathString += "zsmart/USR006_copy"
// 		idxDate = 2
// 	} else if types == "rso" {
// 		PathString += "zsmart/USR007_copy"
// 		idxDate = 2
// 	} else if types == "rto" {
// 		PathString += "zsmart/USR008_copy"
// 		idxDate = 2
// 	} else if types == "responsesrallocation" {
// 		PathString += "zsmart/testcase_upload/sr/response_allocation"
// 		idxDate = 0
// 		destPath = dt.FilesRootFolder + types + "/"
// 	} else if types == "buyback" {
// 		PathString = dt.ConstPathBuybackRequest
// 		idxDate = 2
// 		destPath = dt.FilesRootFolder + "tmp/" + types + "/request/"
// 	}

// 	idxTime = idxDate + 1
// 	var ftp = lib.InitNewFtp()
// 	var FtpConfig = map[string]interface{}{
// 		"url":        "10.1.31.96",    //required
// 		"username":   "sap_po_comlog", //required
// 		"password":   "Smartfren21",   //required
// 		"sourcePath": PathString,      //required
// 		"destPath":   destPath,
// 		"options": map[string]interface{}{
// 			// default separator is '-'
// 			// "separator":    "_",
// 			"#tgl":         strconv.Itoa(idxDate),
// 			"saveToDb":     db,
// 			"GenerateDate": strconv.Itoa(idxDate) + "|" + strconv.Itoa(idxTime),
// 		},
// 	}
// 	// var idxSuccess, idxFailed int
// 	// idxFailed, idxSuccess = 0, 0
// 	errFtp := ftp.DownloadWholeFolder(FtpConfig)
// 	errCode = 1
// 	errStr = "Success Download FTP"
// 	if errFtp != nil {
// 		errCode = -1
// 		errStr = "Error when Download FTP " + errFtp.Error()
// 		// return
// 	}
// 	log.Logf("Process Finish")
// 	return
// }

// func DeleteIndexSlice(s []string, index int) []string {
// 	return append(s[:index], s[index+1:]...)
// }

// func CheckMapKeyString(keyQuery string, arr interface{}) bool {
// 	if rec, ok := arr.(map[string]int); ok {
// 		for key, _ := range rec {
// 			if key == keyQuery {
// 				return true
// 			}
// 		}
// 	}
// 	if rec, ok := arr.(map[string]string); ok {
// 		for key, _ := range rec {
// 			if key == keyQuery {
// 				return true
// 			}
// 		}
// 	}
// 	if rec, ok := arr.(map[string]bool); ok {
// 		for key, _ := range rec {
// 			if key == keyQuery {
// 				return true
// 			}
// 		}
// 	}

// 	if rec, ok := arr.(map[string]interface{}); ok {
// 		for key, _ := range rec {
// 			if key == keyQuery {
// 				return true
// 			}
// 		}
// 	}

// 	if rec, ok := arr.(map[string]float64); ok {
// 		for key, _ := range rec {
// 			if key == keyQuery {
// 				return true
// 			}
// 		}
// 	}

// 	return false
// }

// func EncryptText(plainText string) string {
// 	chipperText := lib.EncryptText([]byte(plainText), dt.ConstEncyprtPass)
// 	return string(chipperText[:])
// }

// func DecryptText(chipperText string) string {
// 	plainText := lib.DecryptText([]byte(chipperText), dt.ConstEncyprtPass)
// 	return string(plainText[:])
// }

// //GetDate to get Date
// func GetDate(optParam ...string) string {
// 	stringDate, _ := lib.GetDate(optParam...)
// 	return stringDate
// }

// func GetDateFormat(dateParam string) string {
// 	formatDate, _ := lib.DeFormatDate(dateParam)
// 	return formatDate
// }

// func DateDiff(dateParam ...string) map[string]float64 {
// 	diffResult := lib.DateDiff(dateParam...)
// 	return diffResult
// }

// func GetTimeInstance(optParam ...string) time.Time {
// 	_, timeInstance := lib.GetDate(optParam...)
// 	return timeInstance
// }
// func ChangeDateFormat(date string) string {
// 	var getDate, getTime, dateNew string
// 	if date != "" {
// 		splitDate := strings.Split(date, " ")
// 		if len(splitDate) > 0 {
// 			getDate = splitDate[0]
// 			getTime = splitDate[1]
// 		} else {
// 			getDate = splitDate[0]
// 		}
// 	}

// 	scan := regexp.MustCompile("(0?[1-9]|[12][0-9]|3[01])/(0?[1-9]|1[012])/((19|20)\\d\\d)")
// 	if scan.MatchString(getDate) {
// 		splitChar := strings.Split(getDate, "/")
// 		dateNew = splitChar[2] + "-" + splitChar[1] + "-" + splitChar[0] + " " + getTime
// 	}
// 	return dateNew
// }
func SendResponses(rw http.ResponseWriter, dt interface{}) {
	js, err := json.Marshal(dt)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(js)
}

// func StoreNotificationReminder(label string, emailData map[string]string) (int64, error) {
// 	var execResult int64
// 	var execErr error
// 	loyaltyConn := lib.InitDb("loyalty")

// 	// typeIdInterface := param[0].(int64)
// 	// typeId := strconv.FormatInt(typeIdInterface, 10)

// 	emailLabel := "ECL-" + label
// 	tableName := "sf_lty_mtr_email_config"
// 	detailTableName := "sf_lty_trx_email"
// 	mainCondition := "EmailLabel = '" + emailLabel + "'"
// 	_, getReminder, _ := loyaltyConn.GetDetailData(tableName, mainCondition)
// 	if getReminder == nil {
// 		var newData = map[string]string{
// 			"EmailLabel": emailLabel,
// 			"Sender":     "inhouse.app@smartfren.com",
// 			"Recipient":  "-",
// 			"Subject":    emailData["Subject"],
// 			"Body":       emailData["Body"],
// 		}
// 		execResult, execErr = loyaltyConn.InsertData(tableName, newData)
// 	} else {
// 		if emailData["Content"] != "" {
// 			_, checkDetailEmail, _ := loyaltyConn.GetDetailData(detailTableName, "Label = '"+emailLabel+"' and Sent = '0'", "Content")
// 			if checkDetailEmail == nil {
// 				var contentData = map[string]string{
// 					"Content":   emailData["Content"],
// 					"Label":     emailLabel,
// 					"CreatedBy": "ecl_" + dt.ConstSys,
// 				}
// 				execResult, execErr = loyaltyConn.InsertData(detailTableName, contentData)
// 			} else {
// 				_, execResult, execErr = loyaltyConn.UpdateData(tableName, map[string]string{"Content": emailData["Content"]}, map[string]string{"Label": emailLabel, "Sent": "0"})
// 			}

// 		}
// 	}
// 	return execResult, execErr
// }

// func handleError(err error) {
// 	if err != nil {
// 		log.Logf("Error: ", err)
// 	}
// }

// func StructToArr(structName interface{}) map[string]string {
// 	// refl := reflect.TypeOf(dt.SodogiColumns).Name()
// 	var reflectValue = reflect.ValueOf(structName)

// 	if reflectValue.Kind() == reflect.Ptr {
// 		reflectValue = reflectValue.Elem()
// 	}

// 	var reflectType = reflectValue.Type()
// 	var ArrStruct = make(map[string]string)
// 	for i := 0; i < reflectValue.NumField(); i++ {
// 		ArrStruct[reflectType.Field(i).Name] = reflectValue.Field(i).Interface().(string)
// 	}
// 	return ArrStruct
// }

// func ReplaceAllStringFor(oldString string, stringTypes string) (newString string) {

// 	regexStr := "[^a-zA-Z0-9]+"
// 	if stringTypes == "alpha" {
// 		regexStr = "[^a-zA-Z]+"
// 	} else if stringTypes == "numeric" {
// 		regexStr = "[^0-9]+"
// 	}
// 	reg, err := regexp.Compile(regexStr)
// 	if err != nil {
// 		log.Errorf("Error when check string is %s for %s", stringTypes, oldString)
// 		newString = ""
// 		return
// 	}
// 	newString = reg.ReplaceAllString(oldString, "")
// 	return
// }

// func ReadDirectory(path string) ([]string, error) {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		log.Errorf("Error when Open Folder %s : %#v", path, err)
// 		return nil, err
// 	}
// 	fileList, err := f.Readdir(-1)
// 	f.Close()
// 	if err != nil {
// 		log.Errorf("Error when Get List of Folder %s : %#v", path, err)
// 		return nil, err
// 	}

// 	var files []string
// 	for _, file := range fileList {
// 		if !file.IsDir() {
// 			files = append(files, file.Name())
// 		}
// 	}

// 	return files, nil
// }

// func visit(files *[]string) filepath.WalkFunc {
// 	return func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			log.Errorf("%#v", err)
// 		}
// 		*files = append(*files, path)
// 		return nil
// 	}
// }

// func ValidationError(w http.ResponseWriter, r *http.Request, request interface{}) int {
// 	var body []byte
// 	var ResponseCode int
// 	var ResponseStatus string
// 	var ResponseMessage string
// 	var ResponseDataResult interface{}
// 	var statusCode int
// 	body, err := ioutil.ReadAll(r.Body)

// 	err = json.Unmarshal(body, &request)

// 	if err != nil {
// 		statusCode = 401
// 		ResponseCode = 401
// 		ResponseStatus = "Error"
// 		ResponseMessage = "Error type data in value"
// 		ResponseDataResult = err.Error()
// 	} else {
// 		statusCode = 200
// 		ResponseCode = 200
// 		ResponseStatus = "Success"
// 		ResponseMessage = "Type data in value match"
// 		ResponseDataResult = ""
// 	}

// 	var Response = dt.GlobalResponse{
// 		ResponseCode:       ResponseCode,
// 		ResponseStatus:     ResponseStatus,
// 		ResponseMessage:    ResponseMessage,
// 		ResponseDataResult: ResponseDataResult,
// 	}
// 	if statusCode != 200 {
// 		SendResponses(w, Response)
// 	}

// 	return statusCode
// }
// func UploadFile(r *http.Request, destination string, prefix_name string) (string, error) {
// 	var err error
// 	var filename, extension string
// 	r.ParseMultipartForm(32 << 20)
// 	uploadedFile, handler, errors := r.FormFile("file")

// 	if errors != nil {
// 		err = errors
// 	} else {
// 		filename = handler.Filename
// 		splitFileName := strings.Split(filename, ".")
// 		extension = splitFileName[len(splitFileName)-1]
// 		filename = prefix_name + "_" + GetDate("Ymd_his") + "." + extension
// 	}
// 	defer uploadedFile.Close()

// 	if strings.Contains(strings.ToLower(dt.ExtensionPermission), strings.ToLower(extension)) {
// 		dir, errors := os.Getwd()
// 		if errors != nil {
// 			err = errors
// 		}
// 		splitDir := strings.Split(destination, "/")
// 		directory := ""
// 		for k, v := range splitDir {
// 			if k > 0 {
// 				directory += "/" + v
// 			} else {
// 				directory += v
// 			}
// 			if _, err := os.Stat(directory); os.IsNotExist(err) {
// 				os.Mkdir(directory, 0777)
// 			}
// 		}
// 		fileLocation := filepath.Join(dir, destination, filename)
// 		targetFile, errors := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0755)
// 		if errors != nil {
// 			err = errors
// 		}
// 		defer targetFile.Close()
// 		if _, errors := io.Copy(targetFile, uploadedFile); err != nil {
// 			err = errors
// 		}
// 	}
// 	return filename, err
// }
// func ReadFileXls(pathFile string) ([]map[string]string, error) {
// 	xlsx, err := excelize.OpenFile(pathFile)
// 	rows, err := xlsx.GetRows("Sheet1")
// 	field := make([]string, 0)

// 	dataRows := make([]map[string]string, 0)
// 	if err == nil {
// 		for k, v := range rows {
// 			if k == 1 {
// 				for _, v1 := range v {
// 					field = append(field, strings.Replace(v1, " ", "", 1))
// 				}
// 			} else if k > 1 {
// 				values := make(map[string]string, 0)
// 				i := 0
// 				for _, v1 := range v {
// 					values[field[i]] = v1
// 					i++
// 				}
// 				dataRows = append(dataRows, values)
// 			}

// 		}
// 	}
// 	return dataRows, err

// }
// func DetectMethod(w http.ResponseWriter, r *http.Request) string {
// 	return r.Method
// }
// func CompareHeaderFile(fieldInDb string, fieldInFile []map[string]string) bool {
// 	var status bool = false
// 	splitField := strings.Split(fieldInDb, ",")

// 	fieldName := make([]string, 0)
// 	if len(splitField) > 0 {
// 		for _, field := range splitField {
// 			splitFieldName := strings.Split(field, "_")
// 			name := ""
// 			if len(splitFieldName) > 0 {
// 				for _, v1 := range splitFieldName {
// 					name += strings.Title(strings.ToLower(v1))
// 				}

// 			}
// 			fieldName = append(fieldName, name)
// 		}
// 	}
// 	JmlFieldInDb := len(fieldName)
// 	JmlFiledTrue := 0
// 	if len(fieldInFile) > 0 {
// 		for k, _ := range fieldInFile[0] {
// 			for _, v1 := range fieldName {
// 				if v1 == k {
// 					JmlFiledTrue++
// 				}
// 			}
// 		}

// 	}
// 	if JmlFieldInDb == JmlFiledTrue {
// 		status = true
// 	}

// 	return status
// }
// func RecordLog(dataLog map[string]string, db lib.DbConnection) error {
// 	var err error
// 	_, erors := db.InsertData("log_activity", dataLog)
// 	if erors != nil {
// 		err = erors
// 	}
// 	return err
// }
// func SortInt(value []int, sorts string) []int {
// 	if strings.ToLower(sorts) == "asc" {
// 		sort.Ints(value)
// 	} else if sorts == "desc" {
// 		sort.Sort(sort.Reverse(sort.IntSlice(value)))
// 	}
// 	return value
// }
// func SortString(value []string, sorts string) []string {
// 	if strings.ToLower(sorts) == "asc" {
// 		sort.Sort(sort.StringSlice(value))
// 	} else if sorts == "desc" {
// 		sort.Sort(sort.Reverse(sort.StringSlice(value)))
// 	}
// 	return value
// }
// func SortByValue(value map[string]int, sorts string) []map[string]int {
// 	crossArray := make(map[int]string)
// 	getInt := make([]int, len(value))
// 	// result := make(map[string]int)
// 	result := []map[string]int{}
// 	for k, v := range value {
// 		crossArray[v] = k
// 	}
// 	i := 0
// 	for k1, _ := range crossArray {
// 		getInt[i] = k1
// 		i++
// 	}
// 	if strings.ToLower(sorts) == "asc" {
// 		sort.Ints(getInt)
// 	} else if sorts == "desc" {
// 		sort.Sort(sort.Reverse(sort.IntSlice(getInt)))
// 	}
// 	for _, v2 := range getInt {
// 		result = append(result, map[string]int{crossArray[v2]: v2})
// 	}
// 	return result

// }
