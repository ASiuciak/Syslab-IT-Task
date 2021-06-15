package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

/* Funkcja pomocnicza, zwracająca listę nazw
 * plików znajdujących się w katalogu files. */
func readFiles() []string {
	var all_files []string

	/* Czytamy wszystkie pliki z danego katalogu
	 * (ale bez plików z podkatalogów) */
	dir := "./files"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	// Dodajemy do listy odczytane pliki, pomijamy katalogi
	for _, file := range files {
		if !file.IsDir() {
			all_files = append(all_files, file.Name())
		}
	}
	return all_files
}

func listFiles(w http.ResponseWriter, r *http.Request) {
	var files []string = readFiles()
	var all = strings.Join(files[:], "\n")
	w.Write([]byte(all))
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")

	// Zwracamy odpowiedź 500, jeśli nie udało się wczytać pliku
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - could not load the file"))
		fmt.Println(err)
		return
	}
	defer file.Close()
	dst, err := os.Create(filepath.Join("./files", filepath.Base(handler.Filename)))

	// Zwracamy odpowiedź 500, jeśli nie udało się umieścić pliku w /files.
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - could not add file to files directory"))
		fmt.Println(err)
		return
	}
	defer dst.Close()
}

func main() {
	http.HandleFunc("/list", listFiles)
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}
