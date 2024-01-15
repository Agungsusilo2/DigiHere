package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type PostImp struct {
	Id                   int64  `json:"id"`
	ApplicantsName       string `json:"applicant_name"`
	EventName            string `json:"event_name"`
	Date                 string `json:"date"`
	EventVenues          string `json:"event_venues"`
	RequirementMaterials string `json:"requirement_materials"`
}

var BASE_URL = "http://localhost:9888/api/v1"

func (p *PostImp) Index(w http.ResponseWriter, r *http.Request) {
	var posts []PostImp

	response, err := http.Get(BASE_URL + "/applicants")
	if err != nil {
		log.Print(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&posts); err != nil {
		log.Print(err)
	}

	data := map[string]interface{}{
		"posts": posts,
	}

	temp, _ := template.ParseFiles("views/index.gohtml")
	err = temp.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func (p *PostImp) Create(w http.ResponseWriter, r *http.Request) {
	var post PostImp
	var data map[string]interface{}

	id := r.URL.Query().Get("id")

	fmt.Println(id)

	if id != "" {
		res, err := http.Get(BASE_URL + "/applicants/" + id)
		if err != nil {
			log.Print(err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				panic(err)
			}
		}(res.Body)

		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&post); err != nil {
			log.Print(err)
		}

		data = map[string]interface{}{
			"post": post,
		}
	}

	temp, _ := template.ParseFiles("views/create.gohtml")
	err := temp.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func (p *PostImp) Store(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	id := r.FormValue("post_id")
	name := r.FormValue("post_name")
	event := r.FormValue("post_event")
	date := r.FormValue("post_date")
	venue := r.FormValue("post_venue")

	file, handler, err := r.FormFile("post_material")
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to process the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uniqueID := uuid.New()
	filename := strings.Replace(uniqueID.String(), "-", "", -1)
	fileExt := strings.Split(handler.Filename, ".")[1]
	uploadedFileName := fmt.Sprintf("%s.%s", filename, fileExt)

	uploadedFile, err := os.Create("uploads/" + uploadedFileName)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	_, err = io.Copy(uploadedFile, file)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Update or create the applicant in the API with the new file information
	idInt, _ := strconv.ParseInt(id, 10, 64)
	newPost := PostImp{
		Id:                   idInt,
		ApplicantsName:       name,
		EventName:            event,
		Date:                 date,
		EventVenues:          venue,
		RequirementMaterials: uploadedFileName,
	}

	jsonValue, _ := json.Marshal(newPost)
	buff := bytes.NewBuffer(jsonValue)

	var req *http.Request

	if id != "" {
		fmt.Println("Proses update")
		req, err = http.NewRequest(http.MethodPatch, BASE_URL+"/applicants/"+id, buff)
	} else {
		fmt.Println("Proses create")
		req, err = http.NewRequest(http.MethodPost, BASE_URL+"/applicants", buff)
	}

	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusCreated || res.StatusCode == http.StatusOK {
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
		return
	} else {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (p *PostImp) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	req, err := http.NewRequest(http.MethodDelete, BASE_URL+"/applicants/"+id, nil)
	if err != nil {
		log.Print(err)
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Print(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(res.Body)

	fmt.Println(res.StatusCode)
	fmt.Println(res.Status)

	if res.StatusCode == 200 {
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	}
}

func (p PostImp) DownloadFile(writer http.ResponseWriter, request *http.Request) {
	file := request.URL.Query().Get("file")
	if file == "" {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(writer, "Bad Request")
		return
	}
	writer.Header().Set("Content-Disposition", "attachment; filename="+file)
	http.ServeFile(writer, request, "./uploads/"+file)
}
