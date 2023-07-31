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
	obj := IDPPartner{MetadataB64: "ww", PartnerName: "qq", Description: "rrr"}

	postToEndpoint(username, password, url, obj)
	t.Error()
}

func TestReadJsonFromFile(t *testing.T) {
	data := new(Data)
	readJsonFromFile("text.json", data)
	t.Error()
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
	postMultiple(username, password, *data)
	t.Error()
}
