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
func Import(file string, args []string) error {
	importCmd := exec.Command("/usr/local/munki/munkiimport", "-v", "-n")
	for _, arg := range args {
		importCmd.Args = append(importCmd.Args, arg)
	}
	importCmd.Args = append(importCmd.Args, file)
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
	var args []string
	// TODO append r.URL.Query() params to args
	packageName := r.URL.Path[len("/import/"):]
	saveMunkiPkg(r.Body, packageName)
	err := Import("tmp/"+packageName, args)
	if err != nil {
		log.Println(err)
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		newPackageHandler(w, r)
	case "POST":
		formHandler(w, r)
	}
}

// handle form input
func formHandler(w http.ResponseWriter, r *http.Request) {
	args, err := processFormValues(r)
	munkiPkg, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
	}
	packageName := header.Filename
	err = saveMunkiPkg(munkiPkg, packageName)
	if err != nil {
		log.Println(err)
	}
	err = Import("tmp/"+packageName, args)
	if err != nil {
		log.Println(err)
	}

}

// turn form inputs into munkiimport and makepkginfo params.
func processFormValues(r *http.Request) ([]string, error) {
	var args []string
	err := r.ParseMultipartForm(0)
	if err != nil {
		return nil, err
	}
	params := r.MultipartForm.Value
	for key, value := range params {
		if key != "null" {
			kv := "--" + key + "=" + value[0]
			args = append(args, kv)
		}
	}
	return args, nil
}

// save an upload to tmp folder
func saveMunkiPkg(munkiPkg io.Reader, packageName string) error {
	f, err := os.Create("tmp/" + packageName)
	if err != nil {
		return err
	}
	defer f.Close()
	// Save file to disk.
	_, err = io.Copy(f, munkiPkg)
	if err != nil {
		return err
	}
	return nil
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
	http.Handle("/", http.FileServer(http.Dir("html")))
	http.ListenAndServe(":8080", nil)
}
