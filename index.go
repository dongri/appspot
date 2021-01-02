package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./public/img"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./public/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./public/js"))))

	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// Issue ...
type Issue struct {
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// indexHandler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	url := "https://api.github.com/repos/dongri/appspot/issues"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "token "+os.Getenv("GITHUB_TOKEN"))

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP Request error:", err)
		return
	}

	statusCode := resp.StatusCode
	if statusCode != 200 {
		fmt.Println("HTTP Status error:", statusCode)
		return
	}
	byteArray, _ := ioutil.ReadAll(resp.Body)
	jsonBytes := ([]byte)(byteArray)
	var issues []Issue

	if err := json.Unmarshal(jsonBytes, &issues); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return
	}

	index := path.Join("views", "index.html")
	header := path.Join("views", "header.html")
	footer := path.Join("views", "footer.html")
	tmpl, err := template.ParseFiles(index, header, footer)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"issues": issues,
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
