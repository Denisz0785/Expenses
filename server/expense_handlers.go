// package server define server, router and handlers for api
package server

import (
	"context"
	dto "expenses/dto_expenses"
	"expenses/repository"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Server contains a link to the database connection in the repository
type Server struct {
	repo *repository.ExpenseRepo
}

// NewServer creates new struct Server
func NewServer(c *repository.ExpenseRepo) *Server {
	return &Server{repo: c}
}

// GetExpenseHandler for get all types of expenses by user's id or login or name
func (r *Server) GetExpenseHandler(c *gin.Context) {
	ctx := context.Background()
	user := &dto.TypesExpenseUserParams{}

	if err := c.BindJSON(user); err != nil {
		log.Printf("error get user data:%s", err.Error())
		c.JSON(http.StatusBadRequest, "incorrect user data")
		return
	}

	// []titleExpense keep title of expenses of  user
	titleExpense, err := r.repo.GetTypesExpenseUser(ctx, user)
	if err != nil {
		log.Printf("error get type expense:%s", err.Error())
		c.JSON(http.StatusInternalServerError, "error get type expense")
		return
	}
	// send expense type
	c.JSON(http.StatusOK, titleExpense)

}

// UploadFile uploads file from user, write it to storage of server and write name of the file to database
func (r *Server) UploadFile(c *gin.Context) {
	ctx := context.Background()
	var newFileName, typeFile, absolutePath string
	var userId int
	//check that expense's id from request exists in a database

	expenseID, err := validateIdExpense(c.Param("id"))

	if err != nil {
		log.Printf("error validate id expense:%s", err.Error())
		c.JSON(http.StatusBadRequest, "incorrect expense Id")
		return
	}
	//check that expense's Id exist in a database
	exist, err := r.repo.IsExpenseExists(ctx, expenseID)
	if err != nil || !exist {
		log.Printf("error expense existing:%s", err.Error())
		c.JSON(http.StatusBadRequest, "incorrect expense Id")
		return
	}
	//get user's Id frm database by expense's id
	userId, err = r.repo.GetUserId(ctx, expenseID)
	if err != nil {
		log.Printf("error user existing:%s", err.Error())
		c.JSON(http.StatusBadRequest, "incorrect user Id")
		return
	}
	userIdSring := strconv.Itoa(userId)
	workDir, err := os.Getwd()
	if err != nil {
		log.Println("error of getting path", err)
		return
	}

	//get file info from requests
	fileHeader, err := c.FormFile("files")
	if err != nil {
		log.Printf("incorrect file name:%s", err.Error())
		c.JSON(http.StatusBadRequest, "incorrect file name")
		return
	} else {
		file, err := fileHeader.Open()
		if err != nil {
			log.Printf("error open file:%s", err.Error())
			c.JSON(http.StatusBadRequest, "error open file")
			return
		}
		defer file.Close()
		fileName := fileHeader.Filename
		typeFile, err = checkExtension(fileName)
		if err != nil {
			log.Printf("incorrect file type:%s", err.Error())
			c.JSON(http.StatusBadRequest, "incorrect file type")
			return
		}

		//create path to save the file to server storage
		path := fmt.Sprint(workDir + "/files/" + userIdSring)
		//check of existing folder with this path, if not exists create a new one
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err = os.Mkdir(path, 0766)
			if err != nil {
				log.Println("error of creating new directory", err)
			}
		}
		src := strings.Split(fileName, ".")
		pattern := fmt.Sprint(src[0] + "-*." + src[1])
		newFile, err := os.CreateTemp(path, pattern)
		absolutePath = newFile.Name()
		//newFileName save a new name of the file
		newFileName = filepath.Base(absolutePath)
		if err == nil {
			//copy file from request to the file of server storage
			//err = c.SaveUploadedFile(file, newFile.Name())
			io.Copy(newFile, file)
		} else {
			log.Printf("error of copy file:%s", err.Error())
			c.JSON(http.StatusInternalServerError, "error save file")
		}
		defer newFile.Close()
	}

	//AddFileExpense write info about file to the database
	err = r.repo.AddFileExpense(c, newFileName, expenseID, typeFile)
	if err != nil {
		err1 := os.Remove(absolutePath)
		if err1 != nil {
			log.Printf("delete file error:%v", err1)
			return
		}
		log.Printf("error write filename to db:%s", err.Error())
		c.JSON(http.StatusInternalServerError, "error save file")
		return
	}
	c.JSON(http.StatusOK, newFileName)

}

// DeleteFile removes file by name from the server's storage and database
func (r *Server) DeleteFile(c *gin.Context) {

	var nameFile string
	var expenseID int
	var userId int
	//check that name of file from request is not empty

	nameFile = c.Query("nameFile")
	if nameFile == "" {
		log.Printf("incorrect file name:%s", nameFile)
		c.JSON(http.StatusBadRequest, "incorrect file name")
		return
	}
	//check that expense's id from request exists in a database
	expenseID, err := validateIdExpense(c.Param("id"))
	if err != nil {
		log.Printf("incorrect id expense:%s", err.Error())
		c.JSON(http.StatusBadRequest, "incorrect id expense")
		return
	}
	exist, err := r.repo.IsExpenseExists(c, expenseID)
	if err != nil || !exist {
		log.Printf("incorrect id expense:%s", err.Error())
		c.JSON(http.StatusBadRequest, "incorrect id expense")
		return
	}
	//get user id from database by expense's id
	userId, _ = r.repo.GetUserId(c, expenseID)
	userIdSring := strconv.Itoa(userId)
	workDir, err := os.Getwd()
	if err != nil {
		log.Println("error of getting path", err)
	}
	//create an absolute path of the file by user's id and name of the file
	path := fmt.Sprint(workDir + "/files/" + userIdSring + "/" + nameFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusBadRequest, "file not found. Check name of the file")
		return
	}
	//remove file from the server storage
	err = os.Remove(path)
	if err != nil {
		log.Printf("delete file error:%v", err)
		return
	}
	//remove file from database
	err = r.repo.DeleteFile(c, nameFile, expenseID)
	if err != nil {
		log.Printf("error delete file:%s", err.Error())
		c.JSON(http.StatusInternalServerError, "error delete file")
		return
	}

	c.JSON(http.StatusOK, "file was deleted")
}
