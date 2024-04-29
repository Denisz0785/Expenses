package server

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

// validateIdExpense validate id is not empty and convert id's type string to int
func validateIdExpense(id string) (int, error) {
	var expenseID int
	var err error
	if id != "" {
		expenseID, err = strconv.Atoi(id)
		if err != nil {
			return -1, errors.New("incorrect id")
		}
	} else {
		return -1, errors.New("incorrect id")
	}
	return expenseID, nil
}

// checkExtension cheks type of file from request
func checkExtension(fileName string) (typeFile string, err error) {
	src := strings.Split(fileName, ".")
	//check name of file has extension and create a new file with random characters appended to the name
	if len(src) == 2 {
		extension := strings.ToLower(src[1])
		switch extension {
		case "doc", "pdf", "txt":
			typeFile = "document"
		case "jpg", "jpeg", "png", "gif", "raw", "svg", "bmp", "ico", "tiff", "webp":
			typeFile = "image"
		case "mp4", "webm", "mov", "avi", "flv", "wmv", "mkv", "mpeg", "3gp", "ogv":
			typeFile = "video"
		default:
			log.Println("incorrect file's type")
			return "", errors.New("incorrect file type")
		}
	}
	return typeFile, nil
}
