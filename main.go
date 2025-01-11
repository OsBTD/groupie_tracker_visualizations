package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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

type RelationsResponse struct {
	Index []Relations `json:"index"`
}

func fetchData(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making GET request: %w", err)
	}
	defer response.Body.Close()

	data := json.NewDecoder(response.Body)

	return data.Decode(target)
}

func main() {
	relationsResponse := RelationsResponse{}
	artists := []Artists{}

	err := fetchData("https://groupietrackers.herokuapp.com/api/relation", &relationsResponse)
	if err != nil {
		fmt.Printf("Error fetching relations: %v\n", err)
		return
	}

	err2 := fetchData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	if err2 != nil {
		fmt.Printf("Error fetching artists: %v\n", err2)
		return
	}

	relationsMap := make(map[int]Relations)
	for _, relation := range relationsResponse.Index {
		relationsMap[relation.ID] = relation
	}
	for i := range artists {
		if relation, found := relationsMap[artists[i].ID]; found {
			artists[i].DatesLocations = relation
		}
	}
	TheWeeknd := Artists{
		Image:        "/static/assets/xo.jpeg",
		ID:           54,
		Name:         "The Weeknd",
		Members:      []string{"Abel Tesfaye"},
		CreationDate: 2009,
		FirstAlbum:   "House of baloons",
		DatesLocations: Relations{
			ID:             54,
			DatesLocations: map[string][]string{"new_york_usa": {"27-11-2016", "26-11-2016"}, "toronto_canada": {"05-09-2016", "04-09-2016"}, "oujda_morocco": {"02-12-2016", "01-12-2016"}},
		},
	}
	artists = append([]Artists{TheWeeknd}, artists...)
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
