package server

import (
	"fmt"
	"net/http"
	"text/template"

	function "server/ascii"
)

type ascii struct {
	AsciiArt string
	Error    string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	home, err := template.New("index.html").ParseFiles("./template/index.html")
	if err != nil {
		ErrorPageHandler(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	w.WriteHeader(http.StatusOK)

	home.Execute(w, nil)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("about.html").ParseFiles("./template/about.html")

	if err != nil {
		ErrorPageHandler(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	w.WriteHeader(http.StatusOK)

	tmpl.Execute(w, nil)
}

func InstructionsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("instructions.html").ParseFiles("./template/instructions.html")
	if err != nil {
		ErrorPageHandler(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	w.WriteHeader(http.StatusOK)

	tmpl.Execute(w, nil)
}

func ErrorPageHandler(w http.ResponseWriter, statusCode int, message string) {
	tmpl, err := template.New("error.html").ParseFiles("./template/error.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	data := ascii{Error: message}
	tmpl.Execute(w, data)
}

func ArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPageHandler(w, http.StatusBadRequest, "400 Bad Request")
		return
	}

	r.ParseForm()
	text := r.FormValue("input")

	banner := r.FormValue("bannerfile")
	ascii_Art, err := function.Input(text, banner)

	if text == "" || banner == "" {
		ErrorPageHandler(w, http.StatusNotFound, "404 Not Found")
		return
	}

	banners := []string{"standard", "thinkertoy", "shadow", "ac"}
	for i := range banners {
		if banner != banners[i] && i == len(banners)-1 {
			ErrorPageHandler(w, http.StatusBadRequest, "400 Bad Request")
			return
		} else if banner == banners[i] {
			break
		}
	}
	if err != nil {
		if err.Error() == "file not found" {

			ErrorPageHandler(w, http.StatusNotFound, "404 Not Found")
			return
		}

		ErrorPageHandler(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	data := ascii{AsciiArt: ascii_Art}

	home := template.Must(template.ParseFiles("template/index.html"))
	err2 := home.Execute(w, data)
	if err2 != nil {
		fmt.Println(err2)
	}
}
