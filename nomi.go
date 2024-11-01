package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	//	"github.com/vhalmd/nomi-go-sdk"
)

type Message struct {
	MessageText string `json:"messageText"`
}
type ReplyText struct {
	ReplyMessage struct {
		Text string `json:"text"`
	} `json:"replyMessage"`
}
type GetNomisResponse struct {
	Nomis []Nomi `json:"nomis"`
}
type Nomi struct {
	UUID             string    `json:"uuid"`
	Gender           string    `json:"gender"`
	Name             string    `json:"name"`
	Created          time.Time `json:"created"`
	RelationshipType string    `json:"relationshipType"`
}

type Config struct {
	Apikey string       `json:"apiKey"`
	Nomi   []NomiConfig `json:"nomis"`
}
type NomiConfig struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	Gender string `json:"gender"`
}

var filePath string = "config.json"
var apikey string
var nomiId string
var nomiName string

func main() {
	if fileExists(filePath) {
		jsonData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("error reading ", filePath)
			return
		}
		var currentData Config

		err = json.Unmarshal(jsonData, &currentData)
		if err != nil {
			fmt.Println("error unmarshaling data: ", err)
			return
		}
		apikey = currentData.Apikey

			fmt.Print("Enter the name of the Nomi to chat with:\n")
		for {
			for _, nomi := range currentData.Nomi {
				fmt.Printf("Name: %s, Gender: %s\n", nomi.Name, nomi.Gender)
			}

			var userInput string

			fmt.Scan(&userInput)

			found := false
			for _, nomi := range currentData.Nomi {
				if userInput == nomi.Name {
					nomiName = nomi.Name
					nomiId = nomi.Id
					startChatting()
					found = true
					break
				}
			}

			if !found {
				fmt.Printf("No Nomi with the name '%s' was found, please try again.\n", userInput)
			}
		}

	} else {
		fmt.Println("No config file found, please paste here your api key to generate config:")
		fmt.Scan(&apikey)
		//fmt.Println(apikey)
		generateConfig(apikey)
	}
}
func generateConfig(tempApiKey string) {
	req, err := http.NewRequest(http.MethodGet, "https://api.nomi.ai/v1/nomis", nil)
	if err != nil {
		fmt.Println("error making the request, ", err)
		return
	}
	req.Header.Add("Authorization", tempApiKey)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error making the request: ", err)
		return
	}
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error reading the response: ", err)
		return
	}
	fmt.Println(string(b))
	var res GetNomisResponse
	err = json.Unmarshal(b, &res)
	if err != nil {
		fmt.Println("error unmarshalling the response: ", err)
		return
	}
	var nomiSlice strings.Builder
	for i, nomi := range res.Nomis {
		nomiSlice.WriteString(fmt.Sprintf(`{  
		"name": "%s",
		"id": "%s",
		"relationshipType": "%s",
		"gender": "%s"
	}`, nomi.Name, nomi.UUID, nomi.RelationshipType, nomi.Gender))

		// Add a comma and newline after each item except the last one
		if i < len(res.Nomis)-1 {
			nomiSlice.WriteString(",\n")
		}
	}
	finalString := nomiSlice.String()
	fmt.Println(finalString)
	towrite := fmt.Sprintf(`{
	"apiKey": "%s",
	"nomis": [
    %s
    ]
}`, apikey, finalString)
	os.WriteFile(filePath, []byte(towrite), 0644)
}

func startChatting() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Nomi chat with ", nomiName)

	for {
		fmt.Print("You> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}
		messageBody := Message{
			MessageText: input,
		}
		url := fmt.Sprintf("https://api.nomi.ai/v1/nomis/%s/chat", nomiId)

		reqBody, err := json.Marshal(messageBody)
		if err != nil {
			fmt.Println("errorr marshaling the body:", err)
			return
		}

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Println("error building post request:", err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", apikey)

		response, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("error making post request:", err)
			return
		}
		defer response.Body.Close()

		b, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("error reading body:", err)
			return
		}
		var data ReplyText
		err = json.Unmarshal(b, &data)
		if err != nil {
			fmt.Println("error unmarshaling reply data:", err)
			return
		}
		fmt.Println(nomiName + ">" + data.ReplyMessage.Text)
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
