package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func upload(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		io.WriteString(w,
			`<html>
			<head>
				<title>Upload!</title>
			</head>
			<body>
				<form enctype="multipart/form-data" action="/" method="post">
					<input type="file" name="fileupload" onchange="submit();">
				</form>
			</body>
			</html>`)
	} else if req.Method == http.MethodPost {
		if err := req.ParseMultipartForm(20 << 20); err != nil {
			log.Println(err)
			return
		}
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
		log.Printf("Received %s\n", handler.Filename)
		http.Redirect(w, req, "/", http.StatusFound)
	}
}

func main() {
	http.HandleFunc("/", upload)
	log.Println("Receiving files to current directory...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
