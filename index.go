// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START gae_go111_app]

// Sample helloworld is an App Engine app.
package main

// [START import]
import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

// [END import]
// [START main_func]
type Project struct {
	Title string
	URL   string
}

func main() {
	http.HandleFunc("/", indexHandler)

	// [START setting_port]
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	// [END setting_port]
}

// [END main_func]

// [START indexHandler]

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	var projects []Project
	projects = append(projects, Project{"LGTM", "https://lgtm.lol"})
	projects = append(projects, Project{"ToDo", "https://todo.hackerth.com"})
	projects = append(projects, Project{"Play", "https://play.hackerth.com"})

	fp := path.Join("views", "index.html")
	header := path.Join("views", "header.html")
	footer := path.Join("views", "footer.html")
	tmpl, err := template.ParseFiles(fp, header, footer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, projects); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// [END indexHandler]
// [END gae_go111_app]
