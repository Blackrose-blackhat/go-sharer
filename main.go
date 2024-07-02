package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
	"encoding/json"
    "github.com/gorilla/mux"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "./index.html")
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    filePath := filepath.Join("uploads", header.Filename)
    outFile, err := os.Create(filePath)
    if err != nil {
        http.Error(w, "Error creating the file", http.StatusInternalServerError)
        return
    }
    defer outFile.Close()

    _, err = io.Copy(outFile, file)
    if err != nil {
        http.Error(w, "Error saving the file", http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, "File uploaded successfully")
}

func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fileName := vars["filename"]
    filePath := filepath.Join("uploads", fileName)

    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }

    http.ServeFile(w, r, filePath)
}
func listFilesHandler(w http.ResponseWriter, r *http.Request) {
    files, err := os.ReadDir("uploads")
    if err != nil {
        http.Error(w, "Error reading directory", http.StatusInternalServerError)
        return
    }

    var fileNames []string
    for _, file := range files {
        if !file.IsDir() {
            fileNames = append(fileNames, file.Name())
        }
    }

    json.NewEncoder(w).Encode(fileNames)
}


func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", homeHandler)
    r.HandleFunc("/upload", uploadFileHandler).Methods("POST")
    r.HandleFunc("/download/{filename}", downloadFileHandler).Methods("GET")
	r.HandleFunc("/uploads", listFilesHandler).Methods("GET")

    // Serve static files
    fs := http.FileServer(http.Dir("static"))
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

    fmt.Println("Starting server on :8080")
    http.ListenAndServe(":8080", r)
}
