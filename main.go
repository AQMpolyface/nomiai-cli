package main

import (
	"ai/packages/elevenlab"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
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
	Apikey       string       `json:"apiKey"`
	DefaultName  string       `json:"default"`
	ElevenlabKey string       `json:"elevenlab"`
	EnableEleven string       `json:"activateVoice"`
	Nomi         []NomiConfig `json:"nomis"`
***REMOVED***
type NomiConfig struct ***REMOVED***
	Name   string `json:"name"`
	Id     string `json:"id"`
	Gender string `json:"gender"`
***REMOVED***

const filePath string = "config.json"
const audioFile string = "output.mp3"

var elevenlabKey string
var res GetNomisResponse
var currentData Config
var apikey string
var nomiId string
var nomiName string
var activateElevenlab bool
var voiceId2 string = "cgSgspJ2msm6clMCkdW9"

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
		///fmt.Println(apikey)
		elevenlabKey = currentData.ElevenlabKey
		str := "true"
		activateElevenlab = str == currentData.EnableEleven
		//fmt.Println("activateelevenlab:", activateElevenlab)
		//if a name is set in the config.json, it will auitomatically pick it
		if currentData.DefaultName != "" ***REMOVED***
			nomiName = currentData.DefaultName
			updateId(nomiName)
		***REMOVED*** else ***REMOVED***
			fmt.Println("\033[31mits most likely that the config is malformed. Please delete it and re-run the program\033[0m")
			return
		***REMOVED***

	***REMOVED*** else ***REMOVED***
		regenerateConfig()
	***REMOVED***

	startChatting()
***REMOVED***

// function generateConfig generates the config file
func regenerateConfig() ***REMOVED***
	fmt.Println("please paste here your api key to generate your config:")
	fmt.Scan(&apikey)
	//fmt.Println(apikey)
	var activateElevenlabString string
	fmt.Println("do you wanna add an elevenlab api key? (so your nomi can have a voice) (y/n)")
	fmt.Println("\033[31mYou NEED to install mpv if you want this to work\033[0m")

	for ***REMOVED***
		fmt.Scan(&activateElevenlabString)
		switch strings.ToLower(activateElevenlabString) ***REMOVED***
		case "y", "yes":
			isIt := isMpvInstalled()

			if !isIt ***REMOVED***
				fmt.Println("mpv isnt detected. please ensure that it is on your PATH")
				regenerateConfig()
			***REMOVED***
			activateElevenlab = true
			fmt.Println("paste your api key here:")
			fmt.Scan(&elevenlabKey)
			fmt.Println("do you to use a different voice?(the default one is cgSgspJ2msm6clMCkdW9) (y/n) ")
			var voiceId string
			for ***REMOVED***
				fmt.Scan(&voiceId)
				switch strings.ToLower(voiceId) ***REMOVED***
				case "y", "yes":
					fmt.Println("enter the voice id here:")
					fmt.Scan(&voiceId2)
					break
				case "n", "no", "":
					fmt.Println("the default voice id is set")
					voiceId2 = "cgSgspJ2msm6clMCkdW9"
					break
				default:
					fmt.Println("please input y (yes) ot n (no)")
					continue
				***REMOVED***
				break
			***REMOVED***
			break
		case "n", "no", "":
			activateElevenlab = false
			break
		default:
			fmt.Println("please input y (yes) or n (no)")
			continue
		***REMOVED***
		break
	***REMOVED***

	for ***REMOVED***
		fmt.Println("please pick a default nomi (by their name) to chat with (you can change this later, type :h to know more):")
		_, nomis := listAndValidate("", 3)
		fmt.Println(nomis)
		var userInput string
		fmt.Scan(&userInput)
		//fmt.Println(userInput)

		exists, _ := listAndValidate(userInput, 1)

		//fmt.Println("exists:", exists)
		if !exists ***REMOVED***
			fmt.Printf("\033[31mNo Nomi with the name '%s' was found, please pick a valid name (case sensitive)\033[0m\n", userInput)
		***REMOVED*** else ***REMOVED***
			generateConfig()
			break
		***REMOVED***
	***REMOVED***
***REMOVED***
func generateConfig() ***REMOVED***
	_, finalString := listAndValidate("", 2)
	//fmt.Println(finalString)

	var towrite string
	if activateElevenlab ***REMOVED***
		towrite = fmt.Sprintf(`***REMOVED***
	"apiKey": "%s",
   "default": "%s",
   "elevenlab": "%s",
   "activateVoice": "%t",
***REMOVED***
    %s
***REMOVED***
***REMOVED***`, apikey, nomiName, elevenlabKey, activateElevenlab, finalString)
	***REMOVED*** else ***REMOVED***
		towrite = fmt.Sprintf(`***REMOVED***
        "apiKey": "%s",
		"default": "%s",
		"elevenlab": "none",
        "enableEleven": "%t",
	***REMOVED***
    %s
***REMOVED***
***REMOVED***`, apikey, nomiName, activateElevenlab, finalString)
	***REMOVED***
	os.WriteFile(filePath, []byte(towrite), 0644)
***REMOVED***

func startChatting() ***REMOVED***
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("(type :h to see the list of options)")
	fmt.Println("Welcome to the chat with ", nomiName)
	for ***REMOVED***

		fmt.Print("\033[34mYou> \033[0m")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" ***REMOVED***
			continue
		***REMOVED***

		switch input ***REMOVED***
		case ":help", ":h":
			showHelp()
		case ":pchange", ":p":
			changeDefaultNomi(true)
		case ":q", ":quit":
			quitChat()
		case ":c", ":change":
			changeDefaultNomi(false)
		case ":r", ":restart":
			regenerateConfig()
			// To do
			/* 	case ":ed", ":deactivateeleven":
				deactivateEleven()
			case ":ae", "addElevenKey":
				addElevenKey()
			case ":re", ":reload":
				reloadConf() */
		default:
			responseMessage := sendMesssage(input)
			if responseMessage == "" ***REMOVED***
				fmt.Printf("\033[31mthere was a fatal error, sorry :( \033[31m")
				return
			***REMOVED***
			fmt.Println(responseMessage)
			//fmt.Println("activateelevenlab", activateElevenlab)
			if activateElevenlab ***REMOVED***
				go func() ***REMOVED***
					cmd := exec.Command("mpv", audioFile)
					err := cmd.Run()
					if err != nil ***REMOVED***
						fmt.Println("error running mpv:", err)
					***REMOVED***
				***REMOVED***()
			***REMOVED***
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
	//fmt.Println(string(b))
	var data ReplyText
	err = json.Unmarshal(b, &data)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror unmarshaling reply data:%s\033[31m\n", err)
		return ""
	***REMOVED***
	//fmt.Println("activate elevenlab:", activateElevenlab)
	if activateElevenlab ***REMOVED***

		getAudioUwu(data.ReplyMessage.Text)
		exec.Command("mpv", audioFile)
	***REMOVED***
	//ai response
	return fmt.Sprintf("\033[36m%s> %s\033[0m", nomiName, data.ReplyMessage.Text)
***REMOVED***

// help message
func showHelp() ***REMOVED***
	fmt.Println("Help: Enter your message to chat with the Nomi.")
	fmt.Println("Commands:")
	fmt.Println("  :h, :help - Show this help message")
	fmt.Println("  :p, :pchange - Change the default chat Nomi")
	fmt.Println("  :q, :quit - Quit the chat")
	fmt.Println("  :c, :change - Change current chat Nomi")
	fmt.Println("  :r, :restart - restart from scratch the config.json")
***REMOVED***

func changeDefaultNomi(defaultOrCurrent bool) ***REMOVED***
	for ***REMOVED***
		//fmt.Println("Changing the default Nomi...")
		if defaultOrCurrent ***REMOVED***
			fmt.Println("who do you want your default nomi to be?")
		***REMOVED*** else ***REMOVED***
			fmt.Println("which nomi do you wanna talk to now?")
		***REMOVED***
		_, nomis := listAndValidate("", 3)
		fmt.Println(nomis)
		fmt.Scan(&nomiName)
		exist, _ := listAndValidate(nomiName, 1)

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
func listAndValidate(input string, number int) (bool, string) ***REMOVED***
	req, err := http.NewRequest(http.MethodGet, "https://api.nomi.ai/v1/nomis", nil)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror making the request,%s\033[31m", err)
		return false, ""
	***REMOVED***
	req.Header.Add("Authorization", apikey)

	response, err := http.DefaultClient.Do(req)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror making the request:%s\033[31m", err)
		return false, ""
	***REMOVED***
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK ***REMOVED***
		fmt.Printf("\033[31merror: received status code %d\033[31m", response.StatusCode)
		return false, ""
	***REMOVED***

	b, err := io.ReadAll(response.Body)
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror reading the response:%s\033[31m", err)
		return false, ""
	***REMOVED***

	err = json.Unmarshal(b, &res) // Ensure `res` is defined properly
	if err != nil ***REMOVED***
		fmt.Printf("\033[31merror unmarshalling the response:%s\033[31m", err)
		return false, ""
	***REMOVED***

	switch number ***REMOVED***
	case 1:
		for _, nomi := range res.Nomis ***REMOVED***
			if input == nomi.Name ***REMOVED***
				nomiName = nomi.Name
				nomiId = nomi.UUID
				return true, ""
			***REMOVED***
		***REMOVED***
		return false, ""

	case 2:
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
		return true, finalString

	case 3:
		var nomiSlice strings.Builder
		for _, nomi := range res.Nomis ***REMOVED***
			uwu := fmt.Sprintf("Name: %s, Gender: %s\n", nomi.Name, nomi.Gender)
			nomiSlice.WriteString(uwu)
		***REMOVED***
		finalString := nomiSlice.String()
		return true, finalString

	default:
		fmt.Println("Invalid case number provided.")
		return false, ""
	***REMOVED***
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
func getAudioUwu(input string) ***REMOVED***
	if fileExists(audioFile) ***REMOVED***
		err := os.Remove(audioFile)
		if err != nil ***REMOVED***
			fmt.Println("Error clearing audio file", err)
		***REMOVED***
	***REMOVED***
	stopSpinner := make(chan bool)
	go func() ***REMOVED***
		spinChars := `|/-\`
		i := 0
		for ***REMOVED***
			select ***REMOVED***
			case <-stopSpinner:
				return
			default:
				fmt.Printf("\033[95m\rfetching elevenlab audio... %c\033[0m", spinChars[i%len(spinChars)])
				i++
				time.Sleep(100 * time.Millisecond)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	elevenlab.GetAudio(elevenlabKey, voiceId2, audioFile, input)
	stopSpinner <- true
	fmt.Print("\r\033[K")
***REMOVED***
func isMpvInstalled() bool ***REMOVED***
	cmd := exec.Command("mpv", "--version")
	err := cmd.Run()
	if err != nil ***REMOVED***
		fmt.Println("error running mpv:", err)
	***REMOVED***
	return err == nil
***REMOVED***
