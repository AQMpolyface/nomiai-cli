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
	Apikey      string       `json:"apiKey"`
	DefaultName string       `json:"default"`
	Nomi        []NomiConfig `json:"nomis"`
***REMOVED***
type NomiConfig struct ***REMOVED***
	Name   string `json:"name"`
	Id     string `json:"id"`
	Gender string `json:"gender"`
***REMOVED***

const filePath string = "config.json"

var res GetNomisResponse
var currentData Config
var apikey string
var nomiId string
var nomiName string

func main() ***REMOVED***
	if fileExists(filePath) ***REMOVED***
		jsonData, err := os.ReadFile(filePath)
		if err != nil ***REMOVED***
			fmt.Printf("\033[31merror reading %s:%s\033[31m", filePath, err)
			return
		***REMOVED***

		err = json.Unmarshal(jsonData, &currentData)
		if err != nil ***REMOVED***
			fmt.Printf("\033[31merror unmarshaling data:%s\033[31m", err)
			return
		***REMOVED***
		apikey = currentData.Apikey
		//if a name is set in the config.json, it will auitomatically pick it
		if currentData.DefaultName != "" ***REMOVED***
			nomiName = currentData.DefaultName
		***REMOVED*** else ***REMOVED***
			fmt.Println("\033[31mits most likely that the config is malformed. Please delete it and re-run the program\033[0m")
			return
		***REMOVED***

	***REMOVED*** else ***REMOVED***
		fmt.Println("No config file found, please paste here your api key to generate config:")
		fmt.Scan(&apikey)
		//fmt.Println(apikey)

		for ***REMOVED***
			fmt.Println("please pick a default nomi (by their name) to chat with (you can change this later, type :h to know more):")
			listNomiViaApi()
			var userInput string
			fmt.Scan(&userInput)
			//fmt.Println(userInput)
			exists := listNomiViaApi2AndValidate(userInput)

			fmt.Println("exists:", exists)
			if !exists ***REMOVED***
				fmt.Printf("\033[31mNo Nomi with the name '%s' was found, please pick a valid name (case sensitive)\033[0m\n", userInput)
			***REMOVED*** else ***REMOVED***
				generateConfig()
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	startChatting()
***REMOVED***

// function generateConfig generates the config file
func generateConfig() ***REMOVED***

	finalString := listNomiViaApi()
	//fmt.Println(finalString)

	towrite := fmt.Sprintf(`***REMOVED***
	"apiKey": "%s",
    "default": "%s",
***REMOVED***
    %s
***REMOVED***
***REMOVED***`, apikey, nomiName, finalString)
	os.WriteFile(filePath, []byte(towrite), 0644)
***REMOVED***

func startChatting() ***REMOVED***
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("(type :h to see the list of options)")
	fmt.Println("Welcome to Nomi chat with ", nomiName)

	for ***REMOVED***
		fmt.Print("\033[34mYou> \033[0m")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input ***REMOVED***
		case "":
			continue
		case ":help", ":h":
			showHelp()
		case ":pchange", ":p":
			changeDefaultNomi(true)
		case ":q", ":quit":
			quitChat()
		case ":c", ":change":
			changeDefaultNomi(false)
		default:
			responseMessage := sendMesssage(input)
			if responseMessage == "" ***REMOVED***
				fmt.Printf("\033[31mthere was a fatal error, sorry :( \033[31m")
				return
			***REMOVED***
			fmt.Println(responseMessage)
		***REMOVED***
	***REMOVED***
***REMOVED***

// check if the file exists
func fileExists(filePath string) bool ***REMOVED***
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) ***REMOVED***
		return false
	***REMOVED***
	return err == nil
***REMOVED***

func sendMesssage(input string) string ***REMOVED***
	messageBody := Message***REMOVED***
		MessageText: input,
	***REMOVED***
	url := fmt.Sprintf("https://api.nomi.ai/v1/nomis/%s/chat", nomiId)

	reqBody, err := json.Marshal(messageBody)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror marshaling the body:%s\033[31m", err)
		return ""
	***REMOVED***

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil ***REMOVED***
		fmt.Println("\033[31merror building post request:%s\033[31m", err)
		return ""
	***REMOVED***
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", apikey)
	stopSpinner := make(chan bool)
	go func() ***REMOVED***
		spinChars := `|/-\`
		i := 0
		for ***REMOVED***
			select ***REMOVED***
			case <-stopSpinner:
				return
			default:
				fmt.Printf("\033[95m\rWaiting for api response... %c\033[95m", spinChars[i%len(spinChars)])
				i++
				time.Sleep(100 * time.Millisecond)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	response, err := http.DefaultClient.Do(req)
	stopSpinner <- true
	fmt.Print("\r\033[K")
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror making post request:%s\033[31m", err)
		return ""
	***REMOVED***
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror reading body:%s\033[31m", err)
		return ""
	***REMOVED***
	var data ReplyText
	err = json.Unmarshal(b, &data)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror unmarshaling reply data:%s\033[31m\n", err)
		return ""
	***REMOVED***
	//ai response
	return fmt.Sprintf("\033[36m%s> %s\033[0m\n", nomiName, data.ReplyMessage.Text)
***REMOVED***

// help message
func showHelp() ***REMOVED***
	fmt.Println("Help: Enter your message to chat with the Nomi.")
	fmt.Println("Commands:")
	fmt.Println("  :h, :help - Show this help message")
	fmt.Println("  :p, :pchange - Change the default Nomi")
	fmt.Println("  :q, :quit - Quit the chat")
	fmt.Println("  :c, :change - Change current Nomi\n")
***REMOVED***

func changeDefaultNomi(defaultOrCurrent bool) ***REMOVED***
	for ***REMOVED***
		//fmt.Println("Changing the default Nomi...")
		if defaultOrCurrent ***REMOVED***
			fmt.Println("who do you want your default nomi to be?")
		***REMOVED*** else ***REMOVED***
			fmt.Println("which nomi do you wanna talk to now?")
		***REMOVED***
		listNomi()
		fmt.Scan(&nomiName)
		exist := validateNomi(nomiName, &currentData)

		if exist ***REMOVED***
			if !defaultOrCurrent ***REMOVED***
				updateId(nomiName)
				startChatting()
			***REMOVED*** else ***REMOVED***
				generateConfig()
				fmt.Printf("\033[32mChanged default chat to %s\033[0m\n", nomiName)
				fmt.Println("do you want to start chatting with the nomi", nomiName+"? (y/N) default No")
				var userInput string

				for ***REMOVED***
					fmt.Scan(&userInput)
					switch strings.ToLower(userInput) ***REMOVED***
					case "y", "yes":
						updateId(userInput)
						startChatting()
						break
					case "n", "no":
						startChatting()
						break
					default:
						fmt.Println("please enter y (yes) or n (no)")
					***REMOVED***
				***REMOVED***
			***REMOVED***

		***REMOVED*** else ***REMOVED***
			fmt.Printf("\033[31mNo Nomi with the name '%s' was found, please pick a valid name (case sensitive)\033[0m\n", nomiName)
		***REMOVED***
	***REMOVED***
***REMOVED***

func quitChat() ***REMOVED***
	fmt.Println("Exiting chat...")
	os.Exit(0)
***REMOVED***

func listNomi() ***REMOVED***
	for _, nomi := range currentData.Nomi ***REMOVED***
		fmt.Printf("Name: %s, Gender: %s\n", nomi.Name, nomi.Gender)
	***REMOVED***
***REMOVED***
func validateNomi(name string, currentData *Config) bool ***REMOVED***
	//	fmt.Println("currentData is :", currentData)
	//	fmt.Println("name uwuwuiwuwu:", name)
	trimmedName := strings.Trim(name, " ")

	//	fmt.Println("name uwuwuiwuwu:", trimmedName)
	i := 1
	for _, nomi := range currentData.Nomi ***REMOVED***

		fmt.Println("nomi name number " + string(i) + "name:" + nomi.Name)
		if trimmedName == nomi.Name ***REMOVED***
			nomiName = nomi.Name
			nomiId = nomi.Id
			return true
		***REMOVED***
		i++
	***REMOVED***
	return false
***REMOVED***

func listNomiViaApi() string ***REMOVED***
	fmt.Println("apikey:", apikey)
	req, err := http.NewRequest(http.MethodGet, "https://api.nomi.ai/v1/nomis", nil)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror making the request,%s\033[31m", err)
		return ""
	***REMOVED***
	req.Header.Add("Authorization", apikey)

	response, err := http.DefaultClient.Do(req)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror making the request:%s\033[31m", err)
		return ""
	***REMOVED***
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror reading the response:%s\033[31m", err)
		return ""
	***REMOVED***
	//fmt.Println(string(b))
	err = json.Unmarshal(b, &res)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror unmarshalling the response:%s\033[31m", err)
		return ""
	***REMOVED***
	for _, nomi := range res.Nomis ***REMOVED***
		fmt.Printf("Name: %s, Gender: %s\n", nomi.Name, nomi.Gender)
	***REMOVED***
	var nomiSlice strings.Builder
	for i, nomi := range res.Nomis ***REMOVED***
		nomiSlice.WriteString(fmt.Sprintf(`***REMOVED***
		"name": "%s",
		"id": "%s",
		"relationshipType": "%s",
		"gender": "%s"
	***REMOVED***`, nomi.Name, nomi.UUID, nomi.RelationshipType, nomi.Gender))
		if i < len(res.Nomis)-1 ***REMOVED***
			nomiSlice.WriteString(",\n")
		***REMOVED***
	***REMOVED***
	finalString := nomiSlice.String()
	//fmt.Println("finalString uwu:", finalString)
	return finalString
***REMOVED***

func listNomiViaApi2AndValidate(input string) bool ***REMOVED***
	fmt.Println("apikey:", apikey)
	req, err := http.NewRequest(http.MethodGet, "https://api.nomi.ai/v1/nomis", nil)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror making the request,%s\033[31m", err)
		return false
	***REMOVED***
	req.Header.Add("Authorization", apikey)

	response, err := http.DefaultClient.Do(req)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror making the request:%s\033[31m", err)
		return false
	***REMOVED***
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror reading the response:%s\033[31m", err)
		return false
	***REMOVED***
	//fmt.Println(string(b))
	err = json.Unmarshal(b, &res)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror unmarshalling the response:%s\033[31m", err)
		return false
	***REMOVED***
	for _, nomi := range res.Nomis ***REMOVED***
		//fmt.Printf("Name: %s, Gender: %s\n", nomi.Name, nomi.Gender)
		if input == nomi.Name ***REMOVED***
			nomiName = nomi.Name
			nomiId = nomi.UUID
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func updateId(name string) bool ***REMOVED***
	for _, nomi := range currentData.Nomi ***REMOVED***
		//fmt.Printf("Name: %s, Gender: %s\n", nomi.Name, nomi.Gender)
		if name == nomi.Name ***REMOVED***
			nomiName = nomi.Name
			nomiId = nomi.Id
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
