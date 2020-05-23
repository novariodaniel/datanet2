package lib

import (
	"bytes"
	"errors"
	"fmt"
	log "projects/datanet2/logging"
	"strings"

	//"fmt"
	"io"
	"os"

	nsftp "github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SftpLibraries interface {
}

type SftpLibrary struct {
	host, uname, pass, port, srcPath, destPath, fileName string
	conn                                                 *ssh.Client
	client                                               *nsftp.Client
	config                                               *ssh.ClientConfig
	session                                              *ssh.Session
	stdout, stderr                                       bytes.Buffer
	err                                                  error
}

func (sftp *SftpLibrary) setConfig() {
	sftp.config = &ssh.ClientConfig{
		User:            sftp.uname,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(sftp.pass),
		},
	}
}

func InitNewSFTP() SftpLibrary {
	return SftpLibrary{}
}

func (sftp *SftpLibrary) connectSftp() {
	sftp.setConfig()
	sftp.conn, sftp.err = ssh.Dial("tcp", sftp.host+":"+sftp.port, sftp.config)
}

func (sftp *SftpLibrary) createNewClient() {
	sftp.client, sftp.err = nsftp.NewClient(sftp.conn)
}
func (sftp *SftpLibrary) closeSSH() {
	log.Logf("SSH closed")
	sftp.conn.Close()
}

func (sftp *SftpLibrary) closeSFTP() {
	log.Logf("SFTP closed")
	sftp.client.Close()
}

func (sftp *SftpLibrary) storeFileToServer(sourcePath, destPath, filename string) error {
	sftp.client, sftp.err = nsftp.NewClient(sftp.conn)
	var destFile *nsftp.File
	var srcFile *os.File
	sourcePath = CheckLastStringPath(sourcePath)
	destPath = CheckLastStringPath(destPath)
	destFile, sftp.err = sftp.client.Create(destPath + filename)
	if sftp.err != nil {
		return errors.New("Error when Create File in Server, because: " + sftp.err.Error())
	}
	defer destFile.Close()
	var flb = InitNewFileLib()

	srcFile, sftp.err = flb.OpenExistingFile(sourcePath + filename)
	if _, sftp.err = io.Copy(destFile, srcFile); sftp.err != nil {
		// idxFailed++
		return errors.New("Failed copy " + filename + " to " + destPath + ", because: " + sftp.err.Error())
	}
	defer srcFile.Close()
	log.Logf("Copy " + filename + " to " + destPath + " Success")
	return nil
}

func (sftp *SftpLibrary) downloadFileToLocal(sourcePath, destPath, filename string) error {
	sftp.client, sftp.err = nsftp.NewClient(sftp.conn)
	if sftp.err != nil {
		fmt.Println(sftp.err)
	}
	// defer sftp.client.Close()
	var srcFile *nsftp.File
	// var srcFile *os.File
	sourcePath = CheckLastStringPath(sourcePath)
	destPath = CheckLastStringPath(destPath)

	srcFile, sftp.err = sftp.client.Open(sourcePath + filename)
	if sftp.err != nil {
		return errors.New("Error when Create File in Server, because: " + sftp.err.Error())
	}
	defer srcFile.Close()
	var flb = InitNewFileLib()
	sftp.err = flb.CreateFile(destPath, filename)
	if _, sftp.err = io.Copy(srcFile, flb.obj); sftp.err != nil {
		// idxFailed++
		return errors.New("Failed copy " + filename + " to " + destPath + ", because: " + sftp.err.Error())
	}
	defer flb.CloseFile()
	log.Logf("Copy " + filename + " to " + destPath + " Success")
	return nil
}
func (sftp *SftpLibrary) uploadSingleFile(sourcePath, destPath, filename string) error {
	sftp.client, sftp.err = nsftp.NewClient(sftp.conn)
	var destFile *nsftp.File
	var srcFile *os.File
	sourcePath = CheckLastStringPath(sourcePath)
	destPath = CheckLastStringPath(destPath)
	destFile, sftp.err = sftp.client.Create(destPath + filename)
	if sftp.err != nil {
		return errors.New("Error when Create File in Servers, because: " + sftp.err.Error())
	}
	defer destFile.Close()
	var flb = InitNewFileLib()

	srcFile, sftp.err = flb.OpenExistingFile(sourcePath + filename)
	if _, sftp.err = io.Copy(destFile, srcFile); sftp.err != nil {
		// idxFailed++
		return errors.New("Failed copy " + filename + " to " + destPath + ", because: " + sftp.err.Error())
	}
	defer srcFile.Close()
	log.Logf("Copy " + filename + " to " + destPath + " Success")
	defer sftp.closeSSH()
	defer sftp.closeSFTP()
	return nil
}
func (sftp *SftpLibrary) downloadSingleFile(sourcePath, destPath, filename string) error {
	sftp.client, sftp.err = nsftp.NewClient(sftp.conn)
	if sftp.err != nil {
		return (sftp.err)
	}

	var srcFile *nsftp.File
	srcFile, err := sftp.client.Open(sourcePath + filename)

	if err != nil {
		return errors.New("Error when Create File in Server, because: " + sftp.err.Error())
	}

	defer srcFile.Close()
	dstFile, err := os.Create(destPath + filename)
	if err != nil {
		return errors.New("Error when Create File in Server, because: " + sftp.err.Error())
	}
	defer dstFile.Close()
	if _, sftp.err = io.Copy(dstFile, srcFile); sftp.err != nil {
		return errors.New("Failed copy " + filename + " to " + destPath + ", because: " + sftp.err.Error())
	}
	defer sftp.closeSSH()
	defer sftp.closeSFTP()
	return nil
}
func (sftp *SftpLibrary) deleteSingleFile(sourcePath, filename string) error {
	sftp.client, sftp.err = nsftp.NewClient(sftp.conn)
	if sftp.err != nil {
		return (sftp.err)
	}

	err := sftp.client.Remove(sourcePath + filename)

	if err != nil {
		return errors.New("Error when Remove File in Server, because: " + sftp.err.Error())
	}

	defer sftp.closeSSH()
	defer sftp.closeSFTP()
	return nil
}
func (sftp *SftpLibrary) ReadDirectory(sourcePath string) []string {
	sftp.client, sftp.err = nsftp.NewClient(sftp.conn)
	if sftp.err != nil {
		fmt.Println(sftp.err)
	}

	var paths []string
	var nameFile string
	w := sftp.client.Walk(sourcePath)

	for w.Step() {
		if w.Err() != nil {
			continue
		}
		// file[i] = strings.Replace(w.Path(), srcPath, "", -1)
		nameFile = strings.Replace(w.Path(), sourcePath, "", -1)
		if nameFile != "" {
			paths = append(paths, nameFile)
		}
	}
	defer sftp.closeSSH()
	defer sftp.closeSFTP()
	return paths
}
func (sftp *SftpLibrary) newStoreFileToServer(source, dest, filename string) error {
	// clients, _ := sftp2.NewClient()
	return nil
}
func (sftp *SftpLibrary) SetConfig(config map[string]interface{}) error {
	sftp.host = config["url"].(string)
	sftp.uname = config["uname"].(string)
	sftp.pass = config["pass"].(string)
	sftp.port = "22"
	if _, exist, val := IsMapKeyExists(config, "port"); exist {
		sftp.port = val.(string)
	}
	if _, exist, val := IsMapKeyExists(config, "destPath"); exist {
		sftp.destPath = val.(string)
	}
	if _, exist, val := IsMapKeyExists(config, "srcPath"); exist {
		sftp.srcPath = val.(string)
	}
	if _, exist, val := IsMapKeyExists(config, "filename"); exist {
		sftp.fileName = val.(string)
	}
	sftp.connectSftp()
	if sftp.err != nil {
		return errors.New("Error When Connect to Sftp, because: " + sftp.err.Error())
	}

	sftp.createNewClient()
	if sftp.err != nil {
		return errors.New("Error When create new client, because: " + sftp.err.Error())
	}
	return nil
}
func (sftp *SftpLibrary) settingConfig(config map[string]interface{}) error {
	sftp.host = config["url"].(string)
	sftp.uname = config["uname"].(string)
	sftp.pass = config["pass"].(string)
	sftp.port = "22"
	if _, exist, val := IsMapKeyExists(config, "port"); exist {
		sftp.port = val.(string)
	}
	if _, exist, val := IsMapKeyExists(config, "destPath"); exist {
		sftp.destPath = val.(string)
	}
	if _, exist, val := IsMapKeyExists(config, "srcPath"); exist {
		sftp.srcPath = val.(string)
	}
	if _, exist, val := IsMapKeyExists(config, "filename"); exist {
		sftp.fileName = val.(string)
	}
	sftp.connectSftp()
	if sftp.err != nil {
		return errors.New("Error When Connect to Sftp, because: " + sftp.err.Error())
	}
	defer sftp.closeSSH()

	sftp.createNewClient()
	if sftp.err != nil {
		return errors.New("Error When create new client, because: " + sftp.err.Error())
	}
	defer sftp.closeSFTP()
	return nil
}

func (sftp *SftpLibrary) UploadSingleFile(config map[string]interface{}) error {
	sftp.settingConfig(config)
	if sftp.err != nil {
		return sftp.err
	}
	return sftp.storeFileToServer(sftp.srcPath, sftp.destPath, sftp.fileName)
}
func (sftp *SftpLibrary) DownloadSingleFile(config map[string]interface{}) error {
	sftp.settingConfig(config)
	if sftp.err != nil {
		return sftp.err
	}

	return sftp.downloadFileToLocal(sftp.srcPath, sftp.destPath, sftp.fileName)

}
func (sftp *SftpLibrary) UploadFile(config map[string]interface{}) error {
	sftp.SetConfig(config)
	if sftp.err != nil {
		return sftp.err
	}
	return sftp.uploadSingleFile(sftp.srcPath, sftp.destPath, sftp.fileName)
}
func (sftp *SftpLibrary) DownloadFile(config map[string]interface{}) error {
	sftp.SetConfig(config)
	if sftp.err != nil {
		return sftp.err
	}

	return sftp.downloadSingleFile(sftp.srcPath, sftp.destPath, sftp.fileName)

}
func (sftp *SftpLibrary) ReadDir(config map[string]interface{}) []string {
	sftp.SetConfig(config)
	if sftp.err != nil {
		fmt.Println(sftp.err)
	}
	return sftp.ReadDirectory(sftp.srcPath)
}
func (sftp *SftpLibrary) DeleteSingleFile(config map[string]interface{}) error {
	sftp.SetConfig(config)
	if sftp.err != nil {
		fmt.Println(sftp.err)
	}
	return sftp.deleteSingleFile(sftp.srcPath, sftp.fileName)
}

// func CheckLastStringPath(path string) string {
// 	strBytes := []byte(path)
// 	lastIndex := bytes.LastIndexByte(strBytes, byte('/'))
// 	if lastIndex != len(strBytes)-1 {
// 		path += "/"
// 	}
// 	return path
// }
