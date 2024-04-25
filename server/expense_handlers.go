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

type Server struct {
	repo *repository.ExpenseRepo
}

func NewServer(c *repository.ExpenseRepo) *Server {
	return &Server{repo: c}
}

func (r *Server) GetExpenseHandler(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	user := &dto.TypesExpenseUserParams{}
	if req.Method == http.MethodPost {
		if checkId := req.FormValue("id"); checkId != "" {
			id, err := strconv.Atoi(checkId)
			if err != nil {
				fmt.Fprintf(w, "Uncorrect value of id%v", err)
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

func (r *Server) UploadFile(w http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	var nameNewFile string
	var expenseID int
	var userId int
	if expIdString := req.FormValue("id"); expIdString != "" {
		var err error
		expenseID, err = strconv.Atoi(expIdString)
		if err != nil {
			fmt.Fprintf(w, "Incorrect id of expense:%v", err)
			return
		}
		exist, err := r.repo.IsExpenseExists(ctx, expenseID)
		if err != nil || !exist {
			fmt.Fprintf(w, "Incorrect id of expense:%v", err)
			return
		}
	} else {
		http.Error(w, "incorrect id of expense", http.StatusBadRequest)
		return
	}

	userId, _ = r.repo.GetUserId(ctx, expenseID)
	userIdSring := strconv.Itoa(userId)
	workDir, err := os.Getwd()
	if err != nil {
		log.Print("error of gettinh path", err)
	}
	path := fmt.Sprint(workDir + "/files/" + userIdSring)
	fmt.Println(path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0766)
		if err != nil {
			log.Print("error of creating new directory", err)
		}
	}

	file, header, err := req.FormFile("files")
	if err != nil {
		http.Error(w, "incorrect file's name", http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintf(w, "Name: %v, Size: %v\n", header.Filename, header.Size)
		for k, v := range header.Header {
			fmt.Fprintf(w, "Key: %v, Value: %v\n", k, v)
		}

		src := strings.Split(header.Filename, ".")
		if len(src) == 2 {
			pattern := fmt.Sprint(src[0] + "-*." + src[1])

			newFile, err := os.CreateTemp(path, pattern)
			nameNewFile = newFile.Name()
			if err == nil {
				io.Copy(newFile, file)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			defer newFile.Close()
		} else {
			fmt.Fprintln(w, "incorrect file's name")
			return
		}

	}
	defer file.Close()
	filePath, _ := filepath.Abs(nameNewFile)
	err = r.repo.AddFileExpense(ctx, filePath, expenseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "file added. File name is:%s", filepath.Base(nameNewFile))
}
