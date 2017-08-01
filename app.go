package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"fmt"
	"path/filepath"

	"io/ioutil"

	"github.com/gorilla/mux"
)

type Schema struct {
	Doc    string
	Org    string
	Fields map[string]Field
}

type Field struct {
	FieldType string
}

var schemaMap map[string]Schema

func CreateDocsEndpoint(w http.ResponseWriter, req *http.Request) {
	//Get Doc Org and load schema
	docId := mux.Vars(req)["doc"]
	orgId := mux.Vars(req)["org"]
	schema := schemaMap[fmt.Sprintf("%s/%s", orgId, docId)]

	var docs interface{}
	_ = json.NewDecoder(req.Body).Decode(&docs)
	log.Println(schema.Fields["nome"].FieldType)
	json.NewEncoder(w).Encode(docs)
}

func main() {
	getSchema()
	router := mux.NewRouter()
	router.HandleFunc("/docs/{org}/{doc}", CreateDocsEndpoint).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getSchema() {
	schemaMap = make(map[string]Schema)
	searchDir := "schema/"

	fileList := []string{}
	_ = filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	for _, file := range fileList {
		var schema Schema
		dat, _ := ioutil.ReadFile(file)
		_ = json.Unmarshal(dat, &schema)
		schemaMap[fmt.Sprintf("%s/%s", schema.Org, schema.Doc)] = schema
		fmt.Println(schema)
	}
}
