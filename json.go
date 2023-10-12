package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Idp struct {
	MetadataB64                     string `json:"metadataB64,omitempty"`
	PartnerName                     string `json:"partnerName,omitempty"`
	NameIDFormat                    string `json:"nameIDFormat,omitempty"`
	SsoProfile                      string `json:"ssoProfile,omitempty"`
	ProviderID                      string `json:"providerID,omitempty"`
	AssertionConsumerURL            string `json:"assertionConsumerURL,omitempty"`
	LogoutRequestURL                string `json:"logoutRequestURL,omitempty"`
	LogoutResponseURL               string `json:"logoutResponseURL,omitempty"`
	AdminManualCreation             string `json:"adminManualCreation,omitempty"`
	DisplaySigningCertDN            string `json:"displaySigningCertDN,omitempty"`
	DisplaySigningCertIssuerDN      string `json:"displaySigningCertIssuerDN,omitempty"`
	DisplaySigningCertStart         string `json:"displaySigningCertStart,omitempty"`
	DisplaySigningCertExpiration    string `json:"displaySigningCertExpiration,omitempty"`
	DisplayEncryptionCertDN         string `json:"displayEncryptionCertDN,omitempty"`
	DisplayEncryptionCertIssuerDN   string `json:"displayEncryptionCertIssuerDN,omitempty"`
	DisplayEncryptionCertStart      string `json:"displayEncryptionCertStart,omitempty"`
	DisplayEncryptionCertExpiration string `json:"displayEncryptionCertExpiration,omitempty"`
}

type NameIDFormat string

const (
	NameIDFormatEmailaddress NameIDFormat = "emailaddress"
	NameIDFormatUnspecified  NameIDFormat = "unspecified"
)

type IDPPartner struct {
	MetadataB64     string       `json:"metadataB64,omitempty"`
	MetadataURL     string       `json:"metadataURL,omitempty"`
	PartnerName     string       `json:"partnerName,omitempty"`
	Description     string       `json:"description,omitempty"`
	PartnerType     string       `json:"partnerType,omitempty"`
	TenantName      string       `json:"tenantName,omitempty"`
	TenantURL       string       `json:"tenantURL,omitempty"`
	NameIDFormat    NameIDFormat `json:"nameIDFormat,omitempty"`
	SsoProfile      string       `json:"ssoProfile,omitempty"`
	AttributeLDAP   string       `json:"attributeLDAP,omitempty"`
	AttributeSAML   string       `json:"attributeSAML,omitempty"`
	FaWelcomePage   string       `json:"faWelcomePage,omitempty"`
	GenerateNewKeys string       `json:"generateNewKeys,omitempty"`
	ValidityNewKeys string       `json:"validityNewKeys,omitempty"`
	Preverify       bool         `json:"preverify,omitempty"`
	ProviderID      string       `json:"providerID,omitempty"`
	SsoURL          string       `json:"ssoURL,omitempty"`
}

type Partner struct {
	PartnerNameOut string `json:"PartnerNameOut"`
	PartnerNameIN  string `json:"PartnerNameIN"`
	Description    string `json:"Description"`
}

type Data struct {
	URL      string             `json:"url"`
	URL_out  string             `json:"url_out"`
	Username string             `json:"username"`
	Password string             `json:"password"`
	Headers  map[string]string  `json:"headers"`
	Partners map[string]Partner `json:"Partners"`
}

func readFromFile(filename string) []string {
	// Read the URL from a file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")
	return lines
}

func getDataFromURL(url, username, password string) io.ReadCloser {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(url)
	// Set the request headers
	req.Header.Set("Content-Type", "application/json")

	// Set the basic authentication credentials
	req.SetBasicAuth(username, password)
	fmt.Println(username, password)
	// Send the HTTP request
	client := &http.Client{Timeout: time.Second * 2}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(resp)
	fmt.Println("++++")
	fmt.Println(resp.Body)
	if err != nil {
		// Define a JSON string
		jsonString := `{"MetadataB64": "error", "PartnerName": "IDP"}`

		// Create an io.ReadCloser that reads from the JSON string
		reader := ioutil.NopCloser(strings.NewReader(jsonString))
		//defer reader.Close()
		return reader
	}
	return resp.Body
}

func getData(username, password string, data Data) {
	for i, v := range data.Partners {
		fmt.Println(i, v)
		fmt.Println("#####")
		reader := getDataFromURL(data.URL+"/"+v.PartnerNameOut, username, password)
		defer reader.Close()
		dataParsed := new(Idp)

		if err := json.NewDecoder(reader).Decode(dataParsed); err != nil {
			panic(err)
		}
		fmt.Println(dataParsed)
		IDPPartner := new(IDPPartner)
		IDPPartner.MetadataB64 = cleanMetadata(dataParsed.MetadataB64)
		IDPPartner.PartnerName = v.PartnerNameIN
		IDPPartner.Description = fmt.Sprintf("Participent: %s", v.Description)
		saveToFile(v.PartnerNameIN, IDPPartner)

		//postToEndpoint(data.Username, data.Password, data.URL_out, IDPPartner)
	}

}

func cleanMetadata(MetadataB64 string) string {
	MetadataB64 = strings.ReplaceAll(MetadataB64, "\n", "")
	MetadataB64 = strings.ReplaceAll(MetadataB64, "\t", "")
	MetadataB64 = strings.ReplaceAll(MetadataB64, "\\\"", "\"")
	MetadataB64 = base64.StdEncoding.EncodeToString([]byte(MetadataB64))
	return MetadataB64
}

func saveToFile(filename string, data interface{}) {
	// Open the file for writing
	file, err := os.Create(filename + ".json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode the data as JSON and write it to the file
	if err := json.NewEncoder(file).Encode(data); err != nil {
		panic(err)
	}
}

func readJsonFromFile(filename string, obj interface{}) error {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the file contents
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the JSON data
	err = json.Unmarshal(data, obj)
	if err != nil {
		log.Fatal(err)
	}

	// Print the data
	//fmt.Printf("%+v\n", d)

	return err

}

func postToEndpoint(username, password, method, url string, obj IDPPartner) (*http.Response, error) {
	// Marshal the object to JSON
	data, err := json.Marshal(obj)
	if err != nil {
		log.Fatal(err)
	}
	// Create a new HTTP request

	req, err := http.NewRequest(method, url+"/"+obj.PartnerName, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	// Set the Content-Type header to application/json
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp, err
}

func postMultiple(username, password, method string, data Data) {
	partners, err := parseJSONFilesInFolder("./")
	if err != nil {
		log.Fatal(err)
	}
	for _, IDPPartner := range partners {
		if err != nil {
			log.Fatal(err)
		}
		postToEndpoint(username, password, method, data.URL_out, IDPPartner)
	}
}

func parseJSONFilesInFolder(folderPath string) ([]IDPPartner, error) {
	var data []IDPPartner

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil // Skip directories
		}
		if filepath.Ext(path) == ".json" {
			// Read and parse the JSON file
			fileContent, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var item IDPPartner
			if err := json.Unmarshal(fileContent, &item); err != nil {
				return err
			}

			data = append(data, item)
		}
		return nil
	})

	return data, err
}

func MethodSwitch() {
	method := os.Args[1]
	filename := os.Args[2]
	username := os.Args[3]
	password := os.Args[4]

	w, err := regexp.MatchString(`.json`, filename)
	if !w && err == nil {
		filename = "text.json"
	}
	data := new(Data)

	err = readJsonFromFile("text.json", data)
	if err == nil {
		log.Println(err)
	}
	fmt.Println(data)
	switch method {
	case "POST":
		postMultiple(username, password, method, *data)
	case "PUT":
		postMultiple(username, password, method, *data)
	case "GET":
		getData(username, password, *data)
	default:
		fmt.Println("Method is neither POST nor GET")
	}

}
