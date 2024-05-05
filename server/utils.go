package server

import (
	"errors"
	dto "expenses/dto_expenses"

	"github.com/gin-gonic/gin"

	"log"
	"net/http"
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

func validateUserData(req *http.Request) (*dto.User, error) {
	user := dto.User{}

	if name := req.FormValue("name"); name != "" {
		user.Name = name
	} else {
		return nil, errors.New("incorrect name")
	}

	if surname := req.FormValue("surname"); surname != "" {
		user.Surname = surname
	} else {
		return nil, errors.New("incorrect surname")
	}

	if login := req.FormValue("login"); login != "" {
		user.Login = login
	} else {
		return nil, errors.New("incorrect login")
	}

	if pass := req.FormValue("pass"); pass != "" {
		user.Pass = pass
	} else {
		return nil, errors.New("incorrect pass")
	}

	if email := req.FormValue("email"); email != "" {
		user.Email = email
	} else {
		return nil, errors.New("incorrect email")
	}
	return &user, nil
}

func getUserId(c *gin.Context) (userId int, err error) {
	checkId, ok := c.Get("id")
	if !ok {
		log.Println("error get id")
		return 0, errors.New("incorrect id")
	}
	id, ok := checkId.(int)
	if !ok {
		log.Println("incorrect type id")
		return 0, errors.New("incorrect type id")
	}
	return id, nil
}
