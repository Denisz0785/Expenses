// package server define server, router and handlers for api
package server

import (
	"context"
	"encoding/json"
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
func (r *Server) GetExpenseHandler(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	user := &dto.TypesExpenseUserParams{}
	if req.Method == http.MethodPost {
		if checkId := req.FormValue("id"); checkId != "" {
			id, err := strconv.Atoi(checkId)
			if err != nil {
				fmt.Fprintf(w, "Uncorrect value of id%v", err)
				return
			}
			user.Id = id
		}
		login := req.FormValue("login")
		user.Login = login
		name := req.FormValue("name")
		user.Name = name
		fmt.Println(user)
		// []titleExpense keep title of expenses of  user
		titleExpense, err := r.repo.GetTypesExpenseUser(ctx, user)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// write to w json data of expenseList
		err = json.NewEncoder(w).Encode(titleExpense)
		if err != nil {
			fmt.Println("Error of marshalig to json")
		}
		w.Header().Set("Content-Type", "application/json")
		return
	} else {
		io.WriteString(w, "enter a post method")
		return
	}
}

// UploadFile uploads file from user, write it to storage of server and write name of the file to database
func (r *Server) UploadFile(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	var newFileName, typeFile, absolutePath string
	var userId int
	//check that expense's id from request exists in a database
	expenseID, err := validateIdExpense(req.FormValue("id"))
	if err != nil {
		http.Error(w, "incorrect expense Id", http.StatusBadRequest)
		return
	}
	//check that expense's Id exist in a database
	exist, err := r.repo.IsExpenseExists(ctx, expenseID)
	if err != nil || !exist {
		http.Error(w, "incorrect expense Id", http.StatusBadRequest)
		return
	}
	//get user's Id frm database by expense's id
	userId, err = r.repo.GetUserId(ctx, expenseID)
	if err != nil {
		fmt.Fprintln(w, "incorrect expense Id")
		return
	}
	userIdSring := strconv.Itoa(userId)
	workDir, err := os.Getwd()
	if err != nil {
		log.Print("error of getting path", err)
		return
	}

	//get file info from requests
	file, header, err := req.FormFile("files")
	if err != nil {
		http.Error(w, "incorrect file name", http.StatusBadRequest)
		return
	} else {
		fileName := header.Filename
		typeFile, err = checkExtension(fileName)
		if err != nil {
			fmt.Fprintln(w, "incorrect file type")
			return
		}

		//create path to save the file to server storage
		path := fmt.Sprint(workDir + "/files/" + userIdSring)
		//check of existing folder with this path, if not exists create a new one
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err = os.Mkdir(path, 0766)
			if err != nil {
				log.Print("error of creating new directory", err)
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
			io.Copy(newFile, file)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer newFile.Close()
	}
	defer file.Close()

	//AddFileExpense write info about file to the database
	err = r.repo.AddFileExpense(ctx, newFileName, expenseID, typeFile)
	if err != nil {
		err1 := os.Remove(absolutePath)
		if err1 != nil {
			log.Printf("delete file error:%v", err1)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "file added. File name is:%s", newFileName)
}

// DeleteFile removes file by name from the server's storage and database
func (r *Server) DeleteFile(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	var nameFile string
	var expenseID int
	var userId int
	//check that name of file from request is not empty
	if nameFile = req.FormValue("nameFile"); nameFile == "" {
		http.Error(w, "incorrect file name", http.StatusBadRequest)
		return
	}
	//check that expense's id from request exists in a database
	expenseID, err := validateIdExpense(req.FormValue("id"))
	if err != nil {
		http.Error(w, "incorrect expense Id", http.StatusBadRequest)
		return
	}
	exist, err := r.repo.IsExpenseExists(ctx, expenseID)
	if err != nil || !exist {
		http.Error(w, "incorrect expense Id", http.StatusBadRequest)
		return
	}
	//get user id from database by expense's id
	userId, _ = r.repo.GetUserId(ctx, expenseID)
	userIdSring := strconv.Itoa(userId)
	workDir, err := os.Getwd()
	if err != nil {
		log.Print("error of getting path", err)
	}
	//create an absolute path of the file by user's id and name of the file
	path := fmt.Sprint(workDir + "/files/" + userIdSring + "/" + nameFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, "file not found. Check name of the file", http.StatusBadRequest)
		return
	}
	//remove file from the server storage
	err = os.Remove(path)
	if err != nil {
		log.Printf("delete file error:%v", err)
		return
	}
	//remove file from database
	err = r.repo.DeleteFile(ctx, nameFile, expenseID)
	if err != nil {
		http.Error(w, "error of deleting file", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "file was deleted")
}
