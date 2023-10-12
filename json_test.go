package main

import (
	"testing"
)

func TestReadFromFile(t *testing.T) {
	filename := "test"
	readFromFile(filename)
	t.Error()
}

func TestPostToEndpoint(t *testing.T) {
	username := "testuser"
	password := "testpw"
	url := "https://localhost:8000"
	obj := IDPPartner{MetadataB64: "ww", PartnerName: "qq", Description: "rrr"}

	_, err := postToEndpoint(username, password, "", url, obj)

	t.Error(err)
}

func TestReadJsonFromFile(t *testing.T) {
	data := new(Data)
	err := readJsonFromFile("text.json", data)
	if err != nil {
		t.Error(err)
	}
}

func TestGetData(t *testing.T) {

	//getData()
	t.Error()
}

func TestPostMultiple(t *testing.T) {
	username := "testuser"
	password := "testpw"
	data := new(Data)
	readJsonFromFile("text.json", data)
	postMultiple(username, password, "", *data)
	t.Error()
}

func TestParseJSONFilesInFolder(t *testing.T) {
	path := "./"
	w, err := parseJSONFilesInFolder(path)
	if err != nil {
		t.Error(err)
	}
	//getData()
	t.Error(w)
}
