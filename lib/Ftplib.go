package lib

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	log "projects/datanet/sftp/logging"
	"strconv"
	"strings"

	goftp "gopkg.in/dutchcoders/goftp.v1"
)

/*
	FTP Library
	for connected to FTP server
	created by Raditya Pratama
	1 April 2019

	======================== BEGIN ==================================
*/

type FtpConnection interface {
	/* Private modifier access */
	openFtp(url string) error
	auth(uname, passwd string) error
	connectFtp(url, uname, passwd string) error
	changeDir(path string) error
	closeFtp() error
	storeFileToServer(source, destinationPath string) error
	storeMultipleFiles(source []string, destinationPath string) error
	copyFileToLocal(source, destinationPath, filename string) error
	copyMultipleFilesToLocal(source []string, destinationPath string, options ...map[string]string) error

	/* Public Modifier Access */
	GetFtpFiles(url, uname, pswd, sourcePath string) ([]string, []string, error)
	DownloadSingleFile(config map[string]interface{}) error
	DownloadMultipleFiles(config map[string]interface{}) error
	DownloadWholeFolder(config map[string]interface{}) error
	UploadSingleFile(config map[string]interface{}) error
	UploadMultipleFiles(config map[string]interface{}) error
	UploadWholeFolder(config map[string]interface{}) error
}

type FtpData struct {
	err error
	Dt  *goftp.FTP
}

// initiate new FTP Instance
func InitNewFtp() FtpData {
	return FtpData{}
}

func (ftp *FtpData) openFtp(url string) error {
	if ftp.Dt, ftp.err = goftp.Connect(url + ":21"); ftp.err != nil {
		log.Errorf("FTP not connected, check your connection")
		return ftp.err
		// panic(err)
	}

	return nil
}
func (ftp *FtpData) auth(uname, passwd string) error {
	if ftp.err = ftp.Dt.Login(uname, passwd); ftp.err != nil {
		// panic(err)
		log.Errorf("Account not Authorized")
		return ftp.err
	}
	return nil
}

// GetFtpFiles is for get lists of Files in SourcePath
func (ftp *FtpData) GetFtpFiles(url, uname, pswd, sourcePath string) ([]string, []string, error) {
	sourcePath = CheckLastStringPath(sourcePath)
	if ftp.Dt == nil {
		log.Logf("FTP instance not ready")
		if ftp.err = ftp.connectFTP(url, uname, pswd); ftp.err != nil {
			return nil, nil, ftp.err
		}
	}
	if ftp.err = ftp.changeDir(sourcePath); ftp.err != nil {
		return nil, nil, ftp.err
	}

	var files []string
	if files, ftp.err = ftp.Dt.List(""); ftp.err != nil {
		// panic(err)
		log.Logf("Cannot get list of directory, please check the permission")
		return nil, nil, ftp.err
	}
	var fullPathFileLists []string
	var fileNameLists []string
	for _, file := range files {
		file = strings.Replace(file, "	", " ", -1)
		file = ReplaceByArr([]string{"\t", "\r", "\n"}, "", file)
		file = ReplaceByArr([]string{"  ", "   ", "    ", "     "}, " ", file)
		splStr := strings.Split(file, " ")
		// NamaFile := ""
		// NamaFile := strings.TrimSpace(splStr[len(splStr)-1])
		fmt.Println(splStr, "pancen")
		var arrTmp []string
		for i, rangeSpl := range splStr {
			j := 0
			// k := 0
			// fmt.Println(i, rangeSpl)
			if strings.Contains(rangeSpl, ":") {
				fmt.Println(i, "contain :")
				j = i + 1
				if j == len(splStr)-1 {
					arrTmp = splStr[j:]
				} else {
					arrTmp = splStr[j:len(splStr)]
				}
			}

			// if strings.Contains(rangeSpl, ".pdf") {
			// 	fmt.Println(i, "contain .pdf")
			// 	k = i
			// }
		}

		// fmt.Println(arrTmp)
		NamaFile := strings.Join(arrTmp, " ")
		// fmt.Println(newStr)
		// if len(NamaFile) <= 4 {
		// 	if len(splStr) == 10 {
		// 		NamaFile = strings.TrimSpace(splStr[len(splStr)-1])
		// 	} else {
		// 		for x, str := range splStr {
		// 			if strings.Contains(str, "buyback") && (strings.Contains(str, ".txt") || strings.Contains(str, ".TXT")) {
		// 				NamaFile += " " + str
		// 			} else if x > 10 {
		// 				NamaFile += " " + str
		// 			}
		// 		}
		// 		NamaFile = strings.TrimLeft(NamaFile, " ")
		// 	}
		// }
		if NamaFile == "." || NamaFile == ".." {
			continue
		}
		fullPathFileLists = append(fullPathFileLists, sourcePath+NamaFile)
		fileNameLists = append(fileNameLists, NamaFile)
	}
	return fullPathFileLists, fileNameLists, nil
}

/*
	UploadSingleFile is used to Upload Single File into destination Folder
	Created By Raditya Pratama
	Required Param:
	1. url string
	2. username string
	3. password string
	4. sourcePath string
	5. destPath string
	6. fileName string
*/
func (ftp *FtpData) UploadSingleFile(config map[string]interface{}) error {
	sourcePath, destinationPath, initErr := ftp.initConfig(config)
	if initErr != nil {
		return initErr
	}
	// automatic close ftp if operation is done
	defer ftp.closeFtp()
	var fileName string
	if resCode, exist, val := IsMapKeyExists(config, "fileName", true); exist {
		fileName = val.(string)
	} else if resCode == -1 {
		return errors.New("FileName Parameter is Required Parameter")
	}
	destinationPath = CheckLastStringPath(destinationPath)
	sourcePath = CheckLastStringPath(sourcePath)
	return ftp.storeFileToServer(sourcePath+fileName, destinationPath+fileName)
}

/*
	UploadWholeFolder is used to Upload a Whole Folder in the folder path into 1 destination Folder
	Created By Raditya Pratama
	Required Param:
	1. url string
	2. username string
	3. password string
	4. sourcePath string
	5. destPath string
*/
func (ftp *FtpData) UploadWholeFolder(config map[string]interface{}) error {
	sourcePath, destinationPath, initErr := ftp.initConfig(config)
	if initErr != nil {
		return initErr
	}
	// automatic close ftp if operation is done
	defer ftp.closeFtp()
	var fl = InitNewFileLib()
	files, fileErr := fl.GetListsOfFileInsideFolder(sourcePath)
	if fileErr != nil {
		log.Logf("Error di get all files")
		return fileErr
	}
	destinationPath = CheckLastStringPath(destinationPath)
	var options map[string]string
	if _, exist, val := IsMapKeyExists(config, "options"); exist {
		options = val.(map[string]string)
	}
	return ftp.storeMultipleFilesToServer(files, destinationPath, options)
}

/*
	UploadMultipleFiles is used to Upload several files into 1 destination path
	Created By Raditya Pratama
	Required Config Param:
	1. url string
	2. username string
	3. password string
	4. sourcePath string //not used in this function, just filled ""
	   use : sourcePath: ""
	5. destPath string
	6. fileList []string
*/
func (ftp *FtpData) UploadMultipleFiles(config map[string]interface{}) error {
	_, destinationPath, initErr := ftp.initConfig(config)
	if initErr != nil {
		return initErr
	}

	// automatic close ftp if operation is done
	defer ftp.closeFtp()
	var fileList []string
	if resCode, exist, val := IsMapKeyExists(config, "fileList", true); exist {
		fileList = val.([]string)
	} else if resCode == -1 {
		return errors.New("FTP Url Parameter is Required Parameter")
	}
	destinationPath = CheckLastStringPath(destinationPath)

	return ftp.storeMultipleFilesToServer(fileList, destinationPath)
}

func (ftp *FtpData) getConfigurationVar(config map[string]interface{}) (result map[string]string, err error) {
	result = map[string]string{
		"url":        "",
		"uname":      "",
		"pswd":       "",
		"sourcePath": "",
		"destPath":   "",
		"fileName":   "",
	}
	if resCode, exist, val := IsMapKeyExists(config, "url", true); exist {
		result["url"] = val.(string)
	} else if resCode == -1 {
		err = errors.New("FTP Url Parameter is Required Parameter")
	}
	if resCode, exist, val := IsMapKeyExists(config, "username", true); exist {
		result["uname"] = val.(string)
	} else if resCode == -1 {
		err = errors.New("FTP Username Parameter is Required Parameter")
	}
	if resCode, exist, val := IsMapKeyExists(config, "password", true); exist {
		result["pswd"] = val.(string)
	} else if resCode == -1 {
		err = errors.New("FTP Password Parameter is Required Parameter")
	}
	if resCode, exist, val := IsMapKeyExists(config, "sourcePath", true); exist {
		result["sourcePath"] = val.(string)
	} else if resCode == -1 {
		err = errors.New("FTP Source Path Parameter is Required Parameter")
	}
	if resCode, exist, val := IsMapKeyExists(config, "destPath", true); exist {
		result["destPath"] = val.(string)
	} else if resCode == -1 {
		err = errors.New("FTP Destination Path Parameter is Required Parameter")
	}

	return
}

func (ftp *FtpData) storeMultipleFilesToServer(source []string, destPath string, options ...map[string]string) error {
	var option map[string]string
	var getHashtagFolder []string
	if Isset(options, 0) {
		option = options[0]
	}
	for _, filePath := range source {
		var subLists []string
		if strings.Contains(filePath, "/") {
			subLists = strings.Split(filePath, "/")
		} else if strings.Contains(filePath, "\\") {
			subLists = strings.Split(filePath, "\\")
		}
		fileName := subLists[len(subLists)-1]
		tmpPath := destPath
		if option != nil {
			NamaFileWithoutTxt := ReplaceByArr([]string{".txt", ".TXT", ".xls", ".XLS", ".xlsx", ".doc", ".docx", ".pdf"}, "", fileName)

			// log.Logf("option tidak kosong, dengan localPath " + tmpPath)
			separator := "-"
			// Check if separator is not set ?
			if _, exist, vm := IsMapKeyExists(option, "separator"); exist {
				separator = vm.(string)
			}
			splitFileName := strings.Split(NamaFileWithoutTxt, separator)
			if strings.Contains(destPath, "#") {
				getHashtagFolder = strings.Split(destPath, "#")
				for k, val := range getHashtagFolder {
					if k == 0 {
						continue
					}
					if _, exist, vm := IsMapKeyExists(option, "#"+strings.Replace(val, "/", "", -1)); exist {
						var newVM = vm.(string)
						idx, _ := strconv.Atoi(newVM)

						if !Isset(splitFileName, idx) {
							continue
						}
						tmpPath = strings.Replace(tmpPath, "#"+val, splitFileName[idx]+"/", -1)
					} else {
						return errors.New("Must Have Path Replace Pattern in your localPath")
					}
				}
			}
			// localPath = tmpPath
		}
		if strings.Contains(tmpPath, "#") {
			return errors.New("Please remove or replace your '#' on your path pattern, replace with 'options' parameter")
		}
		if errStore := ftp.storeFileToServer(filePath, tmpPath+fileName); errStore != nil {
			return errStore
		}
	}
	return nil
}

func (ftp *FtpData) storeFileToServer(source, destPath string) error {
	fileObj := InitNewFileLib()
	file, err := fileObj.OpenExistingFile(source)
	if err != nil {
		return err
	}
	var pathList []string
	if strings.Contains(destPath, "/") {
		pathList = strings.Split(destPath, "/")
	} else if strings.Contains(destPath, "\\") {
		pathList = strings.Split(destPath, "\\")
	}
	// log.Logf("%s", destPath, pathList, file, source)
	fileName := pathList[len(pathList)-1]
	pathList = pathList[:len(pathList)-1]
	destinationPath := strings.Join(pathList, "/")
	destPath = strings.Replace(destPath, "!", ":", -1)
	if err := ftp.Dt.Stor(destPath, file); err != nil {
		return err
	}
	log.Logf("File ", fileName, " uploaded successfully to ", destinationPath)
	return nil
}

/*
	DownloadSingleFile is used to Download Single File from server into local destination Folder
	Created By Raditya Pratama
	Required Config Param:
	1. url string
	2. username string
	3. password string
	4. sourcePath string
	5. destPath string
	6. fileName string
*/
func (ftp *FtpData) DownloadSingleFile(config map[string]interface{}) error {
	sourcePath, destinationPath, initErr := ftp.initConfig(config)
	if initErr != nil {
		return initErr
	}

	// automatic close ftp if operation is done
	defer ftp.closeFtp()
	var fileName string
	if resCode, exist, val := IsMapKeyExists(config, "fileName", true); exist {
		fileName = val.(string)
	} else if resCode == -1 {
		return errors.New("FileName Parameter is Required Parameter")
	}
	destinationPath = CheckLastStringPath(destinationPath)
	sourcePath = CheckLastStringPath(sourcePath)

	return ftp.copyFileToLocal(sourcePath, destinationPath, fileName)
}

/*
	DownloadMultipleFiles is used to Download several files into 1 destination path
	Created By Raditya Pratama
	Required Config Param:
	1. url string
	2. username string
	3. password string
	4. sourcePath string //not used in this function, just filled ""
	   use : sourcePath: ""
	5. destPath string
	6. fileList []string
*/
func (ftp *FtpData) DownloadMultipleFiles(config map[string]interface{}) error {
	_, destinationPath, initErr := ftp.initConfig(config)
	if initErr != nil {
		return ftp.err
	}

	// automatic close ftp if operation is done
	defer ftp.closeFtp()
	var fileList []string
	if resCode, exist, val := IsMapKeyExists(config, "fileList", true); exist {
		fileList = val.([]string)
	} else if resCode == -1 {
		return errors.New("FTP Url Parameter is Required Parameter")
	}
	destinationPath = CheckLastStringPath(destinationPath)
	return ftp.storeMultipleFilesToServer(fileList, destinationPath)
}

func (ftp *FtpData) initConfig(config map[string]interface{}) (string, string, error) {
	var url, uname, pswd, sourcePath, destinationPath string

	/*
		==================================== For FTP Authentication here ========================================
	*/
	var configVar = make(map[string]string)
	if configVar, ftp.err = ftp.getConfigurationVar(config); ftp.err != nil {
		return "", "", ftp.err
	}
	url = configVar["url"]
	uname = configVar["uname"]
	pswd = configVar["pswd"]
	destinationPath = configVar["destPath"]
	sourcePath = configVar["sourcePath"]
	// connect to FTP + authentication
	if ftp.err = ftp.connectFTP(url, uname, pswd); ftp.err != nil {
		log.Logf("Error when Connect to FTP")
		return "", "", ftp.err
	}

	/*
		=========================	FTP Authentication End Here ==================================
	*/
	return sourcePath, destinationPath, nil
}

func (ftp *FtpData) copyFileToLocal(sourcePath, localPath, fileName string) error {
	var fl = InitNewFileLib()
	sourcePath = CheckLastStringPath(sourcePath) + fileName
	fl.CheckPathAndCreateIfNotExist(localPath)
	_, ftp.err = ftp.Dt.Retr(sourcePath, func(r io.Reader) error {
		fileName = strings.Replace(fileName, ":", "", -1)
		destinationPath := CheckLastStringPath(localPath) + fileName

		// var hasher = sha256.New()
		_, errp := os.Stat(destinationPath)
		// fmt.Println(a, errp)
		if errp == nil {
			log.Logf("Copy Abandon " + fileName + " to " + localPath + " File has detected in local drive")
			// return nil
		}

		out, err := os.Create(destinationPath)
		if _, err = io.Copy(out, r); err != nil {
			// idxFailed++
			log.Logf("Failed copy " + fileName + " to " + localPath + ", because: " + err.Error())
			return err
		}
		log.Logf("Copy " + fileName + " to " + localPath + " Success")
		return nil
	})
	fmt.Println(ftp.err)
	return ftp.err
}

func (ftp *FtpData) copyMultipleFilesToLocal(source []string, localPath string, options ...map[string]interface{}) error {
	var option map[string]interface{}
	var getHashtagFolder []string
	if Isset(options, 0) {
		option = options[0]
	}
	// log.Logf("", options, option)
	fmt.Println(source)
	for _, filePath := range source {

		var subLists []string
		if strings.Contains(filePath, "/") {
			subLists = strings.Split(filePath, "/")
		} else if strings.Contains(filePath, "\\") {
			subLists = strings.Split(filePath, "\\")
		}
		fileName := subLists[len(subLists)-1]
		filePattern := subLists[:len(subLists)-1]
		filePath = strings.Join(filePattern, "/")

		/* Check if Customs Folder Destination */
		tmpPath := localPath
		NamaFileWithoutTxt := ReplaceByArr([]string{".txt", ".TXT", ".xls", ".XLS", ".xlsx", ".doc", ".docx", ".pdf"}, "", fileName)

		separator := "-"
		if _, exist, vm := IsMapKeyExists(option, "separator"); exist {
			separator = vm.(string)
		}
		splitFileName := strings.Split(NamaFileWithoutTxt, separator)
		// os.Exit(2)
		if option != nil {

			// log.Logf("option tidak kosong, dengan localPath " + tmpPath)

			if strings.Contains(localPath, "#") {

				getHashtagFolder = strings.Split(localPath, "#")
				for k, val := range getHashtagFolder {
					if k == 0 {
						continue
					}
					// log.Logf("val hashtag "+val+" -> ", option)
					if _, exist, vm := IsMapKeyExists(option, "#"+strings.Replace(val, "/", "", -1)); exist {
						var newVM = vm.(string)
						idx, _ := strconv.Atoi(newVM)

						if !Isset(splitFileName, idx) {
							continue
						}
						tmpPath = strings.Replace(tmpPath, "#"+val, splitFileName[idx]+"/", -1)
					} else {
						return errors.New("Must Have Path Replace Pattern in your localPath")
					}
				}
			}
		}
		if strings.Contains(tmpPath, "#") {
			return errors.New("Please remove or replace your '#' on your path pattern, replace with 'options' parameter")
		}
		// if strings.Contains(tmpPath, ":") {
		// 	return errors.New("Please remove or replace your '#' on your path pattern, replace with 'options' parameter")
		// }
		if errStore := ftp.copyFileToLocal(filePath, tmpPath, fileName); errStore != nil {
			log.Logf("Error: " + errStore.Error())
			return errStore
		}

		// if _, exist, vm := IsMapKeyExists(option, "saveToDb"); exist {
		// 	dbObj := vm.(DbConnection)
		// 	_, checkFile, _ := dbObj.GetDetailData("master_file", "Filename = '"+fileName+"'")
		// 	if checkFile != nil {
		// 		continue
		// 	}
		// 	// success copy and insert to db
		// 	var genDate, genDateFormat string
		// 	if _, exist, vm := IsMapKeyExists(option, "GenerateDate"); exist {
		// 		value := vm.(string)
		// 		pecahIndex := strings.Split(value, "|")
		// 		for k, idx := range pecahIndex {
		// 			nIdx, _ := strconv.Atoi(idx)
		// 			if k > 0 {
		// 				genDate += " "
		// 			}
		// 			if idx != "" && idx != "-1" && !strings.Contains(fileName, "buyback") {
		// 				genDate += splitFileName[nIdx]
		// 			}
		// 		}
		// 		if strings.Contains(fileName, "buyback") {
		// 			pecahIndex = strings.Split(NamaFileWithoutTxt, "-")
		// 			genDate = pecahIndex[1] + " " + pecahIndex[2]
		// 		}
		// 		genDateFormat, _ = GetDate("y-m-d h:i:s", strings.TrimSpace(genDate))
		// 		if genDateFormat == "" {
		// 			gettingDate := strings.Split(NamaFileWithoutTxt, "_")
		// 			genDateFormat, _ = GetDate("y-m-d h:i:s", strings.TrimSpace(gettingDate[len(gettingDate)-1]))
		// 		}
		// 	}

		// 	/* dbColumn must have this col */
		// 	var willBeInsert = map[string]string{
		// 		"Filename":      fileName,
		// 		"GeneratedDate": genDateFormat,
		// 		"CreatedBy":     "system",
		// 		"Source":        "FTP",
		// 	}

		// 	var res int64

		// 	//table_name must have "master_file" words
		// 	res, ftp.err = dbObj.InsertData("master_file", willBeInsert)
		// 	if res <= 0 {
		// 		log.Logf("Failed When Insert " + fileName + " with Error " + ftp.err.Error())
		// 	} else {
		// 		log.Logf("Succes Insert " + fileName + " to Database")
		// 	}

		// }
	}
	return nil
}

func CheckLastStringPath(path string) string {
	strBytes := []byte(path)
	lastIndex := bytes.LastIndexByte(strBytes, byte('/'))
	if lastIndex != len(strBytes)-1 {
		path += "/"
	}
	return path
}

/*
	DownloadWholeFolder is used to Download all files in the source folder into 1 destination path
	Created By Raditya Pratama
	Required Config Param:
	1. url string
	2. username string
	3. password string
	4. sourcePath string -> folder Path on FTP
	5. destPath string

	Optional Config Parameter:
	1. options map[string]interface{} -> is used when destPath contains '#' in the pattern
		a. separator string -> default is '-'
		b. #........,
		c. saveToDb DbConnection -> used for saved each process to database tables (master_file)
		d. GenerateDate string -> used for filled GeneratedDate column in 'master_file' tables, if not exists default is null

	example of using 'options' parameter:

	destPath 				: files/abc/def/#tgl/#waktu/#hari
	fileName-ex-in-folder	: ABC-Monday-01012019-111053.TXT
	indexAfterSplit			:  0	1		2		3

	using : options: map[string]interface{} {
		"#tgl" 			: "2", -> index dari tgl di Filename
		"#waktu" 		: "3", -> index dari waktu di Filename
		"#hari" 		: "1", -> index dari hari di Filename
		"saveToDb" 		: dbConnection, -> dbObjInstance
		"GenerateDate" 	: "2|3", -> index dari tgl dan waktu di FileName, dengan pemisah '|'
	}
*/
func (ftp *FtpData) DownloadWholeFolder(config map[string]interface{}) (error, []string) {
	sourcePath, destinationPath, initErr := ftp.initConfig(config)
	var rtr []string
	if initErr != nil {
		return initErr, rtr
	}
	// automatic close ftp if operation is done
	defer ftp.closeFtp()
	var fullPathFiles []string
	url := config["url"].(string)
	uname := config["username"].(string)
	pswd := config["password"].(string)
	fullPathFiles, rtr, ftp.err = ftp.GetFtpFiles(url, uname, pswd, sourcePath)
	// fmt.Println(fullPathFiles)
	if ftp.err != nil {
		return ftp.err, rtr
	}
	var options map[string]interface{}
	if _, exist, val := IsMapKeyExists(config, "options"); exist {
		options = val.(map[string]interface{})
	}

	return ftp.copyMultipleFilesToLocal(fullPathFiles, destinationPath, options), rtr
}

/*
	DownloadWholeFolder is used to Download all files in the source folder into 1 destination path
	Created By Raditya Pratama
	Required Config Param:
	1. url string
	2. username string
	3. password string
	4. sourcePath string -> folder Path on FTP
	5. destPath string

	Optional Config Parameter:
	1. saveToDb bool -> to make sure your downloaded file will save to DB or not
	2. dbObj DbConnection -> database Object to connect database
	3. separator string -> used to separate filename
	4. pathReplace-x -> to replace your path pattern
	   path ex : /path/to/file/#tgl/#waktu
	   filename ex : ABC-123455-21032019-111054.txt
	                  0     1        2       3
	   use : pathReplace1 : 2
	   use : pathReplace2 : 3
	5. indexDate/indexTime -> index of slice array after split filename by separator,
	   filename ex : ABC-123455-21032019-111054.txt
					  0     1       2       3
		using : indexDate : 2
		using : indexTime : 3
*/
// func (ftp *FtpData) DownloadWholeFolderOld(config map[string]interface{}) error {
// 	sourcePath, destinationPath, initErr := ftp.initConfig(config)

// 	if initErr != nil {
// 		return initErr
// 	}
// 	// automatic close ftp if operation is done
// 	defer ftp.closeFtp()

// 	var fullPathFiles []string
// 	url := config["url"].(string)
// 	uname := config["username"].(string)
// 	pswd := config["password"].(string)
// 	fullPathFiles, ftp.err = ftp.GetFtpFiles(url, uname, pswd, sourcePath)

// 	saveToDb := false
// 	if _, exist, val := IsMapKeyExists(config, "saveToDb"); exist {
// 		saveToDb = val.(bool)
// 	}

// 	separator := "-"
// 	if _, exist, val := IsMapKeyExists(config, "separator"); exist {
// 		separator = val.(string)
// 	}

// 	var db DbConnection
// 	if saveToDb {
// 		if resCode, exist, val := IsMapKeyExists(config, "dbObj", true); exist {
// 			if resCode == -1 {
// 				return errors.New("dbObj is Required Parameter")
// 			}
// 			db = val.(DbConnection)
// 		}
// 	}
// 	var idxSuccess, idxFailed int
// 	idxFailed, idxSuccess = 0, 0

// 	getHashtagFolder := strings.Split(destinationPath, "#")
// 	var fl = InitNewFileLib()
// 	for _, file := range files {

// 		file = strings.Replace(file, "	", " ", -1)
// 		file = ReplaceByArr([]string{"\t", "\r", "\n"}, "", file)

// 		file = ReplaceByArr([]string{"  ", "   ", "    ", "     "}, " ", file)
// 		splStr := strings.Split(file, " ")
// 		NamaFile := strings.TrimSpace(splStr[len(splStr)-1])
// 		if NamaFile == "." || NamaFile == ".." {
// 			continue
// 		}
// 		NamaFileWithoutTxt := ReplaceByArr([]string{".txt", ".TXT", ".xls", ".XLS", ".xlsx", ".doc", ".docx", ".pdf"}, "", NamaFile)

// 		folderPathPattern := destinationPath
// 		ListStringDate := strings.Split(NamaFileWithoutTxt, separator)
// 		if len(getHashtagFolder) > 1 {
// 			for index := 1; index < len(getHashtagFolder); index++ {
// 				// log.Logf("replaceString(%s, #%s, %s/)", destinationPath, getHashtagFolder[index], ListStringDate[config["pathReplace"+strconv.Itoa(index)].(int)])
// 				folderPathPattern = strings.Replace(destinationPath, "#"+getHashtagFolder[index], ListStringDate[config["#"+getHashtagFolder[index]].(int)]+"/", -1)
// 			}
// 		}
// 		// log.Logf("destPathnya: %s", folderPathPattern)
// 		var indexDate, indexTime int
// 		if _, exist, val := IsMapKeyExists(config, "indexDate"); exist {
// 			indexDate = val.(int)
// 		}
// 		if _, exist, val := IsMapKeyExists(config, "indexTime"); exist {
// 			indexTime = val.(int)
// 		}
// 		StringDate := ListStringDate[indexDate]
// 		StringTime := ListStringDate[indexTime]
// 		GenerateDate := StringDate + " " + StringTime

// 		/* Pengecekan apakah setiap subfolder is Exists ? */
// 		fl.CheckPathAndCreateIfNotExist(folderPathPattern)

// 		ftp.Dt.Retr(NamaFile, func(r io.Reader) error {
// 			filePath := folderPathPattern + "/" + NamaFile

// 			// var hasher = sha256.New()
// 			out, err := os.Create(filePath)
// 			if _, err = io.Copy(out, r); err != nil {
// 				idxFailed++
// 				log.Logf("Failed, because: %s", err.Error())
// 				// errStr = "Error when copy data %s" + NamaFile
// 				// errCode = dt.ErrFTPCopyFile
// 			} else {
// 				if saveToDb {
// 					//success copy and insert to db
// 					genDateFormat, _ := GetDate("y-m-d h:i:s", GenerateDate)

// 					/* dbColumn must have this col */
// 					var willBeInsert = map[string]string{
// 						"Filename":      NamaFile,
// 						"GeneratedDate": genDateFormat,
// 						"CreatedBy":     "system",
// 						"Source":        "FTP",
// 					}

// 					var res int64

// 					//table_name must have "master_file" words
// 					res, err = db.InsertData("master_file", willBeInsert)
// 					if res <= 0 {
// 						log.Logf("Failed When Insert %s with Error %s", NamaFile, err.Error())
// 					} else {
// 						log.Logf("Succes Insert %s to Database", NamaFile)
// 					}
// 					idxSuccess++
// 				}
// 			}

// 			return nil
// 		})

// 	}

// 	return nil
// }

func (ftp *FtpData) downloadFileToLocal(source, destination string) error {
	return nil
}

func (ftp *FtpData) connectFTP(url, uname, passwd string) error {
	err := ftp.openFtp(url)
	if err != nil {
		return err
	}
	log.Logf("Connected")
	err = ftp.auth(uname, passwd)
	if err != nil {
		return err
	}
	log.Logf("Account Authorized")

	return nil
}

func (ftp *FtpData) changeDir(path string) error {
	if ftp.err = ftp.Dt.Cwd(path); ftp.err != nil {
		// panic(err)
		log.Logf("Cannot Change Directory, please specify a correct path")
		return ftp.err
	}
	log.Logf("Dir changed to %s", path)
	return nil
}

func (ftp *FtpData) DeleteFile(url, uname, pswd, path string, filename string) error {
	if ftp.err = ftp.connectFTP(url, uname, pswd); ftp.err != nil {
		log.Logf("Error when Connect to FTP")
		return ftp.err
	}
	defer ftp.closeFtp()
	ftp.changeDir(path)
	err := ftp.Dt.Dele(filename)
	if err != nil {
		return err
	}
	// if ftp.err = ftp.Dt.Dele(path); ftp.err != nil {
	// 	// panic(err)
	// 	log.Logf("Cannot Delete File")
	// 	return ftp.err
	// }
	// log.Logf("Dir changed to %s", path)
	return nil
}

func (ftp *FtpData) closeFtp() error {
	log.Logf("FTP Closed")
	return ftp.Dt.Close()
}

/*
	FTP Library
	======================== END ==================================
*/
