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
	Apikey      string       `json:"apiKey"`
	DefaultName string       `json:"default"`
	Nomi        []NomiConfig `json:"nomis"`
}
type NomiConfig struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	Gender string `json:"gender"`
}

const filePath string = "config.json"

var res GetNomisResponse
var currentData Config
var apikey string
var nomiId string
var nomiName string

func main() {
	if fileExists(filePath) {
		jsonData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("\033[31merror reading %s:%s\033[31m", filePath, err)
			return
		}

		err = json.Unmarshal(jsonData, &currentData)
		if err != nil {
			fmt.Printf("\033[31merror unmarshaling data:%s\033[31m", err)
			return
		}
		apikey = currentData.Apikey
		//if a name is set in the config.json, it will auitomatically pick it
		if currentData.DefaultName != "" {
			nomiName = currentData.DefaultName
		} else {
			fmt.Println("\033[31mits most likely that the config is malformed. Please delete it and re-run the program\033[0m")
			return
		}

	} else {
		fmt.Println("No config file found, please paste here your api key to generate config:")
		fmt.Scan(&apikey)
		//fmt.Println(apikey)

		for {
			fmt.Println("please pick a default nomi (by their name) to chat with (you can change this later, type :h to know more):")
			fmt.Println("listapi:")
			_, nomis := listAndValidate("", 3)
			fmt.Println(nomis)
			var userInput string
			fmt.Scan(&userInput)
			//fmt.Println(userInput)

			exists, _ := listAndValidate(userInput, 1)

			//fmt.Println("exists:", exists)
			if !exists {
				fmt.Printf("\033[31mNo Nomi with the name '%s' was found, please pick a valid name (case sensitive)\033[0m\n", userInput)
			} else {
				generateConfig()
				break
			}
		}
	}

	startChatting()
}

// function generateConfig generates the config file
func generateConfig() {
	_, finalString := listAndValidate("", 2)
	//fmt.Println(finalString)

	towrite := fmt.Sprintf(`{
	"apiKey": "%s",
    "default": "%s",
	"nomis": [
    %s
    ]
}`, apikey, nomiName, finalString)
	os.WriteFile(filePath, []byte(towrite), 0644)
}

func startChatting() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("(type :h to see the list of options)")
	fmt.Println("Welcome to Nomi chat with ", nomiName)

	for {
		fmt.Print("\033[34mYou> \033[0m")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
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
			if responseMessage == "" {
				fmt.Printf("\033[31mthere was a fatal error, sorry :( \033[31m")
				return
			}
			fmt.Println(responseMessage)
		}
	}
}

// check if the file exists
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func sendMesssage(input string) string {
	messageBody := Message{
		MessageText: input,
	}
	url := fmt.Sprintf("https://api.nomi.ai/v1/nomis/%s/chat", nomiId)

	reqBody, err := json.Marshal(messageBody)
	if err != nil {
		fmt.Printf("\033[31merror marshaling the body:%s\033[31m", err)
		return ""
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("\033[31merror building post request:%s\033[31m", err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", apikey)
	stopSpinner := make(chan bool)
	go func() {
		spinChars := `|/-\`
		i := 0
		for {
			select {
			case <-stopSpinner:
				return
			default:
				fmt.Printf("\033[95m\rWaiting for api response... %c\033[95m", spinChars[i%len(spinChars)])
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	response, err := http.DefaultClient.Do(req)
	stopSpinner <- true
	fmt.Print("\r\033[K")
	if err != nil {
		fmt.Printf("\033[31merror making post request:%s\033[31m", err)
		return ""
	}
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("\033[31merror reading body:%s\033[31m", err)
		return ""
	}
	var data ReplyText
	err = json.Unmarshal(b, &data)
	if err != nil {
		fmt.Printf("\033[31merror unmarshaling reply data:%s\033[31m\n", err)
		return ""
	}
	//ai response
	return fmt.Sprintf("\033[36m%s> %s\033[0m", nomiName, data.ReplyMessage.Text)
}

// help message
func showHelp() {
	fmt.Println("Help: Enter your message to chat with the Nomi.")
	fmt.Println("Commands:")
	fmt.Println("  :h, :help - Show this help message")
	fmt.Println("  :p, :pchange - Change the default Nomi")
	fmt.Println("  :q, :quit - Quit the chat")
	fmt.Println("  :c, :change - Change current Nomi\n")
}

func changeDefaultNomi(defaultOrCurrent bool) {
	for {
		//fmt.Println("Changing the default Nomi...")
		if defaultOrCurrent {
			fmt.Println("who do you want your default nomi to be?")
		} else {
			fmt.Println("which nomi do you wanna talk to now?")
		}
		_, nomis := listAndValidate("", 3)
		fmt.Println(nomis)
		fmt.Scan(&nomiName)
		exist, _ := listAndValidate(nomiName, 1)

		if exist {
			if !defaultOrCurrent {
				updateId(nomiName)
				startChatting()
			} else {
				generateConfig()
				fmt.Printf("\033[32mChanged default chat to %s\033[0m\n", nomiName)
				fmt.Println("do you want to start chatting with the nomi", nomiName+"? (y/N) default No")
				var userInput string

				for {
					fmt.Scan(&userInput)
					switch strings.ToLower(userInput) {
					case "y", "yes":
						updateId(userInput)
						startChatting()
						break
					case "n", "no":
						startChatting()
						break
					default:
						fmt.Println("please enter y (yes) or n (no)")
					}
				}
			}

		} else {
			fmt.Printf("\033[31mNo Nomi with the name '%s' was found, please pick a valid name (case sensitive)\033[0m\n", nomiName)
		}
	}
}

func quitChat() {
	fmt.Println("Exiting chat...")
	os.Exit(0)
}

func listAndValidate(input string, number int) (bool, string) {
	req, err := http.NewRequest(http.MethodGet, "https://api.nomi.ai/v1/nomis", nil)
	if err != nil {
		fmt.Printf("\033[31merror making the request,%s\033[31m", err)
		return false, ""
	}
	req.Header.Add("Authorization", apikey)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("\033[31merror making the request:%s\033[31m", err)
		return false, ""
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("\033[31merror: received status code %d\033[31m", response.StatusCode)
		return false, ""
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("\033[31merror reading the response:%s\033[31m", err)
		return false, ""
	}

	err = json.Unmarshal(b, &res) // Ensure `res` is defined properly
	if err != nil {
		fmt.Printf("\033[31merror unmarshalling the response:%s\033[31m", err)
		return false, ""
	}

	switch number {
	case 1:
		for _, nomi := range res.Nomis {
			if input == nomi.Name {
				nomiName = nomi.Name
				nomiId = nomi.UUID
				return true, ""
			}
		}
		return false, ""

	case 2:
		var nomiSlice strings.Builder
		for i, nomi := range res.Nomis {
			nomiSlice.WriteString(fmt.Sprintf(`{
                "name": "%s",
                "id": "%s",
                "relationshipType": "%s",
                "gender": "%s"
            }`, nomi.Name, nomi.UUID, nomi.RelationshipType, nomi.Gender))
			if i < len(res.Nomis)-1 {
				nomiSlice.WriteString(",\n")
			}
		}
		finalString := nomiSlice.String()
		return true, finalString

	case 3:
		var nomiSlice strings.Builder
		for _, nomi := range res.Nomis {
			uwu := fmt.Sprintf("Name: %s, Gender: %s\n", nomi.Name, nomi.Gender)
			nomiSlice.WriteString(uwu)
		}
		finalString := nomiSlice.String()
		return true, finalString

	default:
		fmt.Println("Invalid case number provided.")
		return false, ""
	}
}

func updateId(name string) bool {
	for _, nomi := range currentData.Nomi {
		//fmt.Printf("Name: %s, Gender: %s\n", nomi.Name, nomi.Gender)
		if name == nomi.Name {
			nomiName = nomi.Name
			nomiId = nomi.Id
			return true
		}
	}
	return false
}
