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
	url := "localhost:8000"
	obj := IDPPartner{MetadataB64: "", PartnerName: "", Description: ""}

	postToEndpoint(username, password, url, obj)
	t.Error()
}

func TestReadJsonFromFile(t *testing.T) {

	readJsonFromFile("text.json")
	t.Error()
}

func TestGetData(t *testing.T) {

	getData()
	t.Error()
}
