package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

type Artists struct {
	Image          string   `json:"image"`
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Members        []string `json:"members"`
	CreationDate   int      `json:"creationDate"`
	FirstAlbum     string   `json:"firstAlbum"`
	RelationsURL   string   `json:"relations"`
	DatesLocations Relations
}

type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func fetchData(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

func main() {
	artists := []Artists{}
	err := fetchData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err != nil {
		fmt.Println("Error fetching artists")
		return
	}

	for i := 0; i < len(artists); i++ {
		var relation Relations
		err := fetchData(artists[i].RelationsURL, &relation)
		if err != nil {
			fmt.Println("Error fetching relation for artist", artists[i].ID)
		} else {
			artists[i].DatesLocations = relation
		}
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl2 := template.Must(template.ParseFiles("templates/error.html"))
	about := template.Must(template.ParseFiles("templates/about.html"))
	readme := template.Must(template.ParseFiles("templates/readme.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			data := "Error " + strconv.Itoa(http.StatusMethodNotAllowed) + ":  Method not allowed."
			w.WriteHeader(http.StatusMethodNotAllowed)
			tmpl2.Execute(w, data)
			return
		}

		if r.URL.Path != "/" {
			data := "Error " + strconv.Itoa(http.StatusNotFound) + ":  Page not found."
			w.WriteHeader(http.StatusNotFound)
			tmpl2.Execute(w, data)
			return
		}

		tmpl.Execute(w, artists)
	})
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			data := "Error " + strconv.Itoa(http.StatusMethodNotAllowed) + ":  Method not allowed."
			w.WriteHeader(http.StatusMethodNotAllowed)
			tmpl2.Execute(w, data)
			return
		}

		if r.URL.Path != "/about" {
			data := "Error " + strconv.Itoa(http.StatusNotFound) + ":  Page not found."
			w.WriteHeader(http.StatusNotFound)
			tmpl2.Execute(w, data)
			return
		}

		about.Execute(w, nil)
	})
	http.HandleFunc("/readme", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			data := "Error " + strconv.Itoa(http.StatusMethodNotAllowed) + ":  Method not allowed."
			w.WriteHeader(http.StatusMethodNotAllowed)
			tmpl2.Execute(w, data)
			return
		}

		if r.URL.Path != "/readme" {
			data := "Error " + strconv.Itoa(http.StatusNotFound) + ":  Page not found."
			w.WriteHeader(http.StatusNotFound)
			tmpl2.Execute(w, data)
			return
		}

		readme.Execute(w, nil)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates"))))
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
