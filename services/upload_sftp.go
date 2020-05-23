package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	dt "projects/datanet2/datastruct"
	lib "projects/datanet2/lib"
)

const (
	urlSftp      string = "10.1.35.39"
	usernameSftp string = "sris_tb"
	passwordSftp string = "2020stsSRIS---"
	portSftp     string = "22"
	destSftp     string = "/var/integration/datanet/"
	srcSftp      string = "/home/vagrant/txt/"
)

//UploadSftpServices is an interface
type UploadSftpServices interface {
}

//UploadSftpService is a struct
type UploadSftpService struct{}

//UploadData is a init function
func UploadData(w http.ResponseWriter, r *http.Request) dt.DatanetResponse {
	var resp = dt.DatanetResponse{}
	// respCode

	countErr := 0
	var errMsg []string

	//Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return resp
	}

	// Unmarshall
	var _producer dt.DatanetRequest
	err = json.Unmarshal(b, &_producer)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return resp
	}

	errorCode, errorStr, listFile := GetFileFTP("")

	arrName := ValidationInsert(listFile)

	log.Println(errorCode, errorStr)

	if len(arrName) > 0 {
		for _, rangeFilename := range arrName {

			sftpObj := lib.InitNewSFTP()

			var sftpConfig = map[string]interface{}{
				"url":      urlSftp,
				"port":     portSftp, //fill when port is not 22
				"uname":    usernameSftp,
				"pass":     passwordSftp,
				"destPath": destSftp,
				"srcPath":  srcSftp,
				"filename": rangeFilename,
			}
			err1 := sftpObj.UploadFile(sftpConfig)
			if err1 != nil {
				countErr++
				tmpErr := "There is error with filename " + rangeFilename
				errMsg = append(errMsg, tmpErr)
			}
			// fmt.Println(err1, _producer.Filename)
		}
	} else {
		fmt.Println("Tidak ada data baru")
	}

	var response dt.DatanetResponse
	var responseData = map[string]interface{}{
		"Filename": _producer.Filename,
		"err":      errMsg,
	}

	response.Status = "success"
	response.Data = map[string]interface{}{
		"data": responseData,
	}

	if countErr > 0 {
		response.Status = "There is some errors"
	}
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
	return response

}

// ValidationInsert is a function for validate and insert data to database
func ValidationInsert(flname []string) []string {
	// dbConn.Open()
	dbConn := lib.InitDb()
	var result, finalResult []string
	var fname string

	defer dbConn.Close()

	for _, arrFileName := range flname {
		rows, err := dbConn.Query("select filename from tsd_file where filename = ?", arrFileName)

		if err != nil {
			panic(err.Error())
			return nil
		}

		for rows.Next() {
			err = rows.Scan(&fname)
			if err != nil {
				panic(err.Error())
			}

			result = append(result, fname)
		}

		defer rows.Close()
	}
	count := 0
	for i := 0; i < len(flname); i++ {
		count = 0
		for j := 0; j < len(result); j++ {
			if flname[i] == result[j] {
				count++
				break
			}
		}
		if count == 0 {
			var dtfile = map[string]string{
				"filename": flname[i],
			}

			finalResult = append(finalResult, flname[i])
			dbConn.InsertData("tsd_file", dtfile)
		}
	}
	return finalResult
}

//GetFileFTP is a function, getting a whole file in source folder
func GetFileFTP(filename string) (errCode int, errStr string, listFile []string) {
	var PathString string

	PathString = "/interface/itf/datanet/"
	// destPath := dt.FilesRootFolder + types + "/#tgl/" //required -> it will be change based on devops
	destPath := "/home/vagrant/txt/"

	var ftp = lib.InitNewFtp()
	var FtpConfig = map[string]interface{}{
		"url":        "10.1.31.96",    //required
		"username":   "sap_po_comlog", //required
		"password":   "Smartfren21",   //required
		"sourcePath": PathString,      //required
		"destPath":   destPath,
	}
	// var idxSuccess, idxFailed int
	// idxFailed, idxSuccess = 0, 0
	errFtp, listFile := ftp.DownloadWholeFolder(FtpConfig)
	// fmt.Println(len(listFile), "--------------")
	errCode = 1
	errStr = "Success Download FTP"
	if errFtp != nil {
		errCode = -1
		errStr = "Error when Download FTP " + errFtp.Error()
		// return
	}
	log.Println("Process Finish")
	return
}
