package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// function Import runs munkiimport command.
func Import(file string) error {
	importCmd := exec.Command("/usr/local/munki/munkiimport", "-v", "-n", file)
	var out bytes.Buffer
	importCmd.Stdout = &out
	err := importCmd.Run()
	if err != nil {
		return err
	}
	fmt.Println(out.String())
	return nil
}

// handle PUT requests.
// copies the binary to a tmp folder and calls munkiimport command.
func newPackageHandler(w http.ResponseWriter, r *http.Request) {
	packageName := r.URL.Path[len("/import/"):]
	f, err := os.Create("tmp/" + packageName)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	// Save file to disk.
	fileLength, err := io.Copy(f, r.Body)
	if err != nil {
		log.Println(err)
	}
	// handle file size mismatch
	if fileLength != r.ContentLength {
		log.Println("mismatched size")
	}

	// import munki pkg
	err = Import("tmp/" + packageName)
	if err != nil {
		log.Println(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		newPackageHandler(w, r)
	}
}

// create directory where uploads are saved to.
func createTmpDir() error {
	if _, err := os.Stat("tmp"); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir("tmp", 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	// create tmp directory
	err := createTmpDir()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/import/", handler)
	http.ListenAndServe(":8080", nil)
}
