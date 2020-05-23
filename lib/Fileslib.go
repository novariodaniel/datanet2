package lib

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	log "projects/datanet2/logging"
	"strings"
)

/*
	Files Library for File Operations
	Created By Raditya Pratama
	1 April 2019
*/

type FileLibraries interface {
	/* Public Modifier Access */
	InitNewFileLib() FileData
	CreateFile(path, fileName string) error
	WriteFile(path, fileName string, contents []string) error
	CreateAndWriteFile(path, fileName string, contents []string) error
	GetListsOfFileInsideFolder(folderPath string) (files []string, err error)
	CheckPathAndCreateIfNotExist(fullPath string) error
	CloseFile() error
	SaveFile() error
	OpenFile(fullPath string) (*os.File, error)
	OpenExistingFile(fullPath string) (*os.File, error)
	WriteExistingFile(fullPath string, contents []string) error
	UploadFile(uploadedFile multipart.File, newPath, fileName string) error

	/* Private Modifier Access */
	checkPath(fullPath string) bool
	writeFile(contents []string) error
	copyFile(uploadedFile multipart.File) error
}

type FileData struct {
	obj *os.File
	err error
}

// initiate new File Library
func InitNewFileLib() FileData {
	return FileData{}
}

func (fd *FileData) CheckPathAndCreateIfNotExist(fullPath string) error {
	subFolderLists := strings.Split(fullPath, "/")
	tmpFolder := ""
	for k, subFolder := range subFolderLists {
		if k > 0 {
			tmpFolder += string(filepath.Separator)
		}
		tmpFolder += subFolder
		if !fd.checkPath(tmpFolder) {
			os.Mkdir("."+string(filepath.Separator)+tmpFolder, 0777)
		}
	}

	return nil
}

func (fd *FileData) checkPath(fullPath string) bool {
	_, fd.err = os.Stat(fullPath)
	return !os.IsNotExist(fd.err)
}

func (fd *FileData) GetListsOfFileInsideFolder(folderPath string) (files []string, err error) {
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	files = files[1:]
	return
}

func (fd *FileData) CreateAndWriteFile(path, fileName string, contents []string) error {
	// Check is the path folder is Exists ?
	path = CheckLastStringPath(path)

	fd.CheckPathAndCreateIfNotExist(path)
	fd.CreateFile(path, fileName)
	defer fd.CloseFile()
	return fd.writeFile(contents)
}

func (fd *FileData) UploadFile(uploadedFile multipart.File, newPath, fileName string) error {
	// Check is the path folder is Exists ?
	fd.CheckPathAndCreateIfNotExist(newPath)

	fd.CreateFile(newPath, fileName)
	defer fd.CloseFile()
	return fd.copyFile(uploadedFile)
}

func (fd *FileData) copyFile(uploadedFile multipart.File) error {
	_, err := io.Copy(fd.obj, uploadedFile)
	if err != nil {
		return err
	}
	return nil
}

func (fd *FileData) CloseFile() error {
	log.Logf("File Closed")
	return fd.obj.Close()
}

func (fd *FileData) SaveFile() error {
	if fd.err = fd.obj.Sync(); fd.err != nil {
		log.Logf("File failed to save, because: %s", fd.err.Error())
		return fd.err
	}
	log.Logf("File Saved")
	return nil
}

func (fd *FileData) CreateFile(path, fileName string) error {
	fd.obj, fd.err = os.Create(path + fileName)
	if fd.err != nil {
		log.Logf("File %s failed to created, because: %s", fd.err.Error())
		return fd.err
	}
	fd.obj.Chmod(0777)
	log.Logf("File %s success to Created", fileName)
	return nil
}

/* open File for existing File */
func (fd *FileData) OpenExistingFile(fullPath string) (*os.File, error) {
	fd.obj, fd.err = os.OpenFile(fullPath, os.O_APPEND|os.O_RDWR, 0755)
	if fd.err != nil {
		log.Logf("Cannot Open File, because: %s", fd.err.Error())
		return nil, fd.err
	}
	return fd.obj, nil
}

/* open File for Create New */
func (fd *FileData) OpenFile(fullPath string) (*os.File, error) {
	fd.obj, fd.err = os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0666)
	if fd.err != nil {
		log.Logf("Cannot Open File, because: %s", fd.err.Error())
		return nil, fd.err
	}
	return fd.obj, nil
}

func (fd *FileData) WriteExistingFile(fullPath string, contents []string) error {
	if !fd.checkPath(fullPath) {
		return errors.New("File Not Found !")
	}
	fd.OpenExistingFile(fullPath)
	defer fd.CloseFile()
	fd.writeFile(contents)
	return nil
}

func (fd *FileData) RemoveFile(path string) error {
	// delete file
	var err = os.Remove(path)
	if err != nil {
		log.Logf("Error when Delete Data: " + err.Error())
		return err
	}

	log.Logf("File Deleted")
	return nil
}

func (fd *FileData) writeFile(contents []string) error {
	for _, content := range contents {
		textContent := content + "\n"
		_, fd.err = fmt.Fprint(fd.obj, textContent)
		if fd.err != nil {
			log.Logf("Content " + textContent + " cannot be write, because: " + fd.err.Error())
			return fd.err
		}
	}
	fd.SaveFile()
	log.Logf("File is Successful to write")
	return nil
}
