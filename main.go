package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// Artists represents the artist data structure
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

// Relations represents the concert dates and locations data
type Relations struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

// RelationsResponse represents the API response structure
type RelationsResponse struct {
	Index []Relations `json:"index"`
}

// ErrorPage represents the data structure for error information
type ErrorPage struct {
	Code    int
	Message string
	Is400   bool
	Is404   bool
	Is500   bool
}

// fetchData makes an HTTP GET request and decodes the JSON response
func fetchData(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making GET request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response code: %d", response.StatusCode)
	}

	return json.NewDecoder(response.Body).Decode(target)
}

// handleError handles error responses consistently across handlers
func handleError(w http.ResponseWriter, tmpl *template.Template, code int, message string) {
	errorPage := ErrorPage{
		Code:    code,
		Message: message,
		Is400:   code == http.StatusBadRequest || code == http.StatusMethodNotAllowed,
		Is404:   code == http.StatusNotFound,
		Is500:   code == http.StatusInternalServerError,
	}
	w.WriteHeader(code)
	if err := tmpl.Execute(w, errorPage); err != nil {
		log.Printf("Error executing error template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	// Parse templates
	templates := make(map[string]*template.Template)
	templateFiles := map[string]string{
		"index":  "templates/index.html",
		"error":  "templates/error.html",
		"about":  "templates/about.html",
		"readme": "templates/readme.html",
	}

	for name, file := range templateFiles {
		tmpl, err := template.ParseFiles(file)
		if err != nil {
			log.Fatalf("Error parsing template %s: %v", name, err)
		}
		templates[name] = tmpl
	}

	// Fetch and prepare data
	var artists []Artists
	var relationsResponse RelationsResponse

	// Fetch relations data
	if err := fetchData("https://groupietrackers.herokuapp.com/api/relation", &relationsResponse); err != nil {
		log.Printf("Error fetching relations: %v", err)
	}

	// Fetch artists data
	if err := fetchData("https://groupietrackers.herokuapp.com/api/artists", &artists); err != nil {
		log.Printf("Error fetching artists: %v", err)
	}

	// Map relations to artists
	relationsMap := make(map[int]Relations)
	for _, relation := range relationsResponse.Index {
		relationsMap[relation.ID] = relation
	}
	for i := range artists {
		if relation, found := relationsMap[artists[i].ID]; found {
			artists[i].DatesLocations = relation
		}
	}

	// Add The Weeknd to artists
	theWeeknd := Artists{
		Image:        "/static/assets/xo.jpeg",
		ID:           54,
		Name:         "The Weeknd",
		Members:      []string{"Abel Tesfaye"},
		CreationDate: 2009,
		FirstAlbum:   "House of baloons",
		DatesLocations: Relations{
			ID: 54,
			DatesLocations: map[string][]string{
				"new_york_usa":   {"27-11-2016", "26-11-2016"},
				"toronto_canada": {"05-09-2016", "04-09-2016"},
				"oujda_morocco": {"02-12-2016", "01-12-2016"},
			},
		},
	}
	artists = append([]Artists{theWeeknd}, artists...)

	// Define route handlers
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			handleError(w, templates["error"], http.StatusNotFound, "Page not found")
			return
		}

		if r.Method != http.MethodGet {
			handleError(w, templates["error"], http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		if err := templates["index"].Execute(w, artists); err != nil {
			log.Printf("Error executing index template: %v", err)
			handleError(w, templates["error"], http.StatusInternalServerError, "Internal server error")
		}
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/about" {
			handleError(w, templates["error"], http.StatusNotFound, "Page not found")
			return
		}

		if r.Method != http.MethodGet {
			handleError(w, templates["error"], http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		if err := templates["about"].Execute(w, nil); err != nil {
			log.Printf("Error executing about template: %v", err)
			handleError(w, templates["error"], http.StatusInternalServerError, "Internal server error")
		}
	})

	http.HandleFunc("/readme", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/readme" {
			handleError(w, templates["error"], http.StatusNotFound, "Page not found")
			return
		}

		if r.Method != http.MethodGet {
			handleError(w, templates["error"], http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		if err := templates["readme"].Execute(w, nil); err != nil {
			log.Printf("Error executing readme template: %v", err)
			handleError(w, templates["error"], http.StatusInternalServerError, "Internal server error")
		}
	})

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("templates"))))

	// Start server
	port := ":8080"
	fmt.Printf("Server started at http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
