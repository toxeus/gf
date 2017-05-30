package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var tmpl = template.Must(template.New("list").Parse(
	`<!DOCTYPE html>
		<html>
		<head>
			<title>Share!</title>
		</head>
		<body>
			<form enctype="multipart/form-data" action="/" method="post">
				<input type="file" name="fileupload" onchange="submit();">
			</form>
			<ul>
				{{- range .}}
				<li><a href="{{.}}">{{.}}</a></li>
				{{- end}}
			</ul>
		</body>
		</html>`))

func share(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		if req.RequestURI != "/" {
			fileName := req.RequestURI[1:]
			f, err := os.Open(fileName)
			if err != nil {
				log.Println(err)
				return
			}
			defer f.Close()
			if _, err := io.Copy(w, f); err != nil {
				log.Println(err)
				return
			}
			w.Header().Add("Content-Type", "application/octet-stream")
			log.Printf("Sent %s to %s\n", fileName, req.Host)
		}
		files := []string{}
		walkFunc := func(path string, info os.FileInfo, err error) error {
			if path == "." {
				return nil
			}
			if info.IsDir() {
				return filepath.SkipDir
			}
			files = append(files, info.Name())
			return nil
		}
		if err := filepath.Walk(".", walkFunc); err != nil {
			log.Println(err)
			return
		}
		if err := tmpl.Execute(w, files); err != nil {
			log.Println(err)
			return
		}
	} else if req.Method == http.MethodPost {
		file, handler, err := req.FormFile("fileupload")
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		f, err := os.Create(handler.Filename)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		log.Printf("Received %s from %s\n", handler.Filename, req.Host)
		http.Redirect(w, req, "/", http.StatusFound)
	}
}

func main() {
	http.HandleFunc("/", share)
	log.Println("Receiving files to current directory...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
