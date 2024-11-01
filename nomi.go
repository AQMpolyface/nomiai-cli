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
)

type Message struct ***REMOVED***
	MessageText string `json:"messageText"`
***REMOVED***
type ReplyText struct ***REMOVED***
	ReplyMessage struct ***REMOVED***
		Text string `json:"text"`
	***REMOVED*** `json:"replyMessage"`
***REMOVED***
type GetNomisResponse struct ***REMOVED***
	Nomis []Nomi `json:"nomis"`
***REMOVED***
type Nomi struct ***REMOVED***
	UUID             string    `json:"uuid"`
	Gender           string    `json:"gender"`
	Name             string    `json:"name"`
	Created          time.Time `json:"created"`
	RelationshipType string    `json:"relationshipType"`
***REMOVED***

type Config struct ***REMOVED***
	Apikey string       `json:"apiKey"`
	Nomi   []NomiConfig `json:"nomis"`
***REMOVED***
type NomiConfig struct ***REMOVED***
	Name   string `json:"name"`
	Id     string `json:"id"`
	Gender string `json:"gender"`
***REMOVED***

var filePath string = "config.json"
var apikey string
var nomiId string
var nomiName string

func main() ***REMOVED***
	if fileExists(filePath) ***REMOVED***
		jsonData, err := os.ReadFile(filePath)
		if err != nil ***REMOVED***
			fmt.Println("error reading ", filePath)
			return
		***REMOVED***
		var currentData Config

		err = json.Unmarshal(jsonData, &currentData)
		if err != nil ***REMOVED***
			fmt.Println("error unmarshaling data: ", err)
			return
		***REMOVED***
		apikey = currentData.Apikey

			fmt.Print("Enter the name of the Nomi to chat with:\n")
		for ***REMOVED***
			for _, nomi := range currentData.Nomi ***REMOVED***
				fmt.Printf("Name: %s, Gender: %s\n", nomi.Name, nomi.Gender)
			***REMOVED***

			var userInput string

			fmt.Scan(&userInput)

			found := false
			for _, nomi := range currentData.Nomi ***REMOVED***
				if userInput == nomi.Name ***REMOVED***
					nomiName = nomi.Name
					nomiId = nomi.Id
					startChatting()
					found = true
					break
				***REMOVED***
			***REMOVED***

			if !found ***REMOVED***
				fmt.Printf("No Nomi with the name '%s' was found, please try again.\n", userInput)
			***REMOVED***
		***REMOVED***

	***REMOVED*** else ***REMOVED***
		fmt.Println("No config file found, please paste here your api key to generate config:")
		fmt.Scan(&apikey)
		//fmt.Println(apikey)
		generateConfig(apikey)
	***REMOVED***
***REMOVED***
//function generateConfig generates the config file
func generateConfig(tempApiKey string) ***REMOVED***
	req, err := http.NewRequest(http.MethodGet, "https://api.nomi.ai/v1/nomis", nil)
	if err != nil ***REMOVED***
		fmt.Println("error making the request, ", err)
		return
	***REMOVED***
	req.Header.Add("Authorization", tempApiKey)

	response, err := http.DefaultClient.Do(req)
	if err != nil ***REMOVED***
		fmt.Println("error making the request: ", err)
		return
	***REMOVED***
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil ***REMOVED***
		fmt.Println("error reading the response: ", err)
		return
	***REMOVED***
	fmt.Println(string(b))
	var res GetNomisResponse
	err = json.Unmarshal(b, &res)
	if err != nil ***REMOVED***
		fmt.Println("error unmarshalling the response: ", err)
		return
	***REMOVED***
	var nomiSlice strings.Builder
	for i, nomi := range res.Nomis ***REMOVED***
		nomiSlice.WriteString(fmt.Sprintf(`***REMOVED***  
		"name": "%s",
		"id": "%s",
		"relationshipType": "%s",
		"gender": "%s"
	***REMOVED***`, nomi.Name, nomi.UUID, nomi.RelationshipType, nomi.Gender))

		// Add a comma and newline after each item except the last one
		if i < len(res.Nomis)-1 ***REMOVED***
			nomiSlice.WriteString(",\n")
		***REMOVED***
	***REMOVED***
	finalString := nomiSlice.String()
	fmt.Println(finalString)
	towrite := fmt.Sprintf(`***REMOVED***
	"apiKey": "%s",
***REMOVED***
    %s
***REMOVED***
***REMOVED***`, apikey, finalString)
	os.WriteFile(filePath, []byte(towrite), 0644)
***REMOVED***

func startChatting() ***REMOVED***
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to Nomi chat with ", nomiName)

	for ***REMOVED***
		fmt.Print("You> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" ***REMOVED***
			continue
		***REMOVED***
		messageBody := Message***REMOVED***
			MessageText: input,
		***REMOVED***
		url := fmt.Sprintf("https://api.nomi.ai/v1/nomis/%s/chat", nomiId)

		reqBody, err := json.Marshal(messageBody)
		if err != nil ***REMOVED***
			fmt.Println("errorr marshaling the body:", err)
			return
		***REMOVED***

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		if err != nil ***REMOVED***
			fmt.Println("error building post request:", err)
			return
		***REMOVED***
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", apikey)

		response, err := http.DefaultClient.Do(req)
		if err != nil ***REMOVED***
			fmt.Println("error making post request:", err)
			return
		***REMOVED***
		defer response.Body.Close()

		b, err := io.ReadAll(response.Body)
		if err != nil ***REMOVED***
			fmt.Println("error reading body:", err)
			return
		***REMOVED***
		var data ReplyText
		err = json.Unmarshal(b, &data)
		if err != nil ***REMOVED***
			fmt.Println("error unmarshaling reply data:", err)
			return
		***REMOVED***
		//ai response
		fmt.Println(nomiName + ">" + data.ReplyMessage.Text)
	***REMOVED***
***REMOVED***
//check if the file exists
func fileExists(filePath string) bool ***REMOVED***
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) ***REMOVED***
		return false
	***REMOVED***
	return err == nil
***REMOVED***
