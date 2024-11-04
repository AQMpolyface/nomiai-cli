package main

import (
	"ai/packages/elevenlab"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	//"regexp"
	"strings"
	"time"

	"github.com/chzyer/readline"
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
	Apikey       string       `json:"apiKey"`
	DefaultName  string       `json:"default"`
	ElevenlabKey string       `json:"elevenlab"`
	EnableEleven string       `json:"activateVoice"`
	VoiceId      string       `json:"voiceid"`
	Nomi         []NomiConfig `json:"nomis"`
}
type NomiConfig struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	Gender string `json:"gender"`
}

const filePath string = "config.json"
const audioFile string = "output.mp3"

// var re = regexp.MustCompile(`\*(.*?)\*`)
var elevenlabKey string
var res GetNomisResponse
var currentData Config
var apikey string
var nomiId string
var nomiName string
var activateElevenlab bool
var voiceId2 string = "cgSgspJ2msm6clMCkdW9"

func startChatting() {

	fmt.Println("(type :h to see the list of options)")
	fmt.Println("Welcome to the chat with ", nomiName)

	//purplePrompt := "\033[35mYou> \033[0m" // \033[35m for purple text, \033[0m to reset color
	bluePrompt := "\033[34mYou> " //\033[0m"

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          bluePrompt,
		InterruptPrompt: "\n",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatalf("Failed to initialize readline: %s", err)
	}
	defer rl.Close()
	for {
		//	fmt.Print("\033[34mYou> "//\033[0m")

		input, err := rl.Readline()
		//	input = re.ReplaceAllString(input, "\033[3m$1\033[0m\033[34m")
		if err != nil {
			if err == readline.ErrInterrupt {
				fmt.Println("\033[31mPlease use :q or :quit if you want to exit the chat\033[0m")
				continue
			}
			fmt.Println("\033[31mError reading input\033[0m", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		switch input {
		case ":help", ":h":
			showHelp()
		case ":pchange", ":p":
			changeDefaultNomi(true)
		case ":q", ":quit":
			quitChat()
		case ":lv", ":listVoices":
			getVoices()
		case ":c", ":change":
			changeDefaultNomi(false)
		case ":r", ":restart":
			regenerateConfig()
			/* todo
			case ":t", ":theme":
				showThemelevenlab.ChangeTheme()
				theme, exists := elevenlab.GetTheme(selectedTheme)
				if !exists {
					log.Fatalf("Theme %s does not exist.", showTheme)
				}*/
		case ":de", ":deactivateeleven":
			deactivateEleven()
		case ":ea", ":actEleven":
			activateEleven()
		case ":ae", "addElevenKey":
			addElevenKey()
		case ":re", ":reload":
			main()
		default:
			responseMessage := sendMesssage(input)
			if responseMessage == "" {
				fmt.Printf("\033[31mthere was a fatal error, sorry :( \033[31m")
				return
			}

			//responseMessage = re.ReplaceAllString(responseMessage, "\033[3m$1\033[0m\033[34m\033[36m")
			//ansi code doesnt work everywhere so deactivated it. hope the resp of the ansi code work

			fmt.Println(responseMessage)
			//fmt.Println("activateelevenlab", activateElevenlab)
			if activateElevenlab {
				go func() {
					cmd := exec.Command("mpv", audioFile)
					err := cmd.Run()
					if err != nil {
						fmt.Println("error running mpv:", err)
					}
				}()
			}
		}
	}
}
func getVoices() {
	voices := elevenlab.GetVoices(elevenlabKey)

	fmt.Println("\033[32mAvailable Voices:\033[0m")
	for _, voice := range voices {
		fmt.Printf("ID: %s, Name: %s\n", voice.VoiceID, voice.Name)
	}
}
func addElevenKey() {
	fmt.Printf("enter your api key here:")
	fmt.Scan(&elevenlabKey)
	activateElevenlab = true

	generateConfig()

}
func activateEleven() {
	fmt.Println("activateleleven:", activateElevenlab)
	if activateElevenlab {
		fmt.Println("elevenlab is already activated")
		return
	} else {
		activateElevenlab = true
		generateConfig()
		fmt.Printf("\033[32msuccessfully activated elevenlab.\033[0m")
	}
}
func deactivateEleven() {
	fmt.Println("activateleleven:", activateElevenlab)
	if !activateElevenlab {
		fmt.Println("elevenlab is already deactivated")
		return
	} else {
		activateElevenlab = false
		generateConfig()
		fmt.Printf("\033[32msuccessfully activated elevenlab.\033[0m")

	}
}
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
		///fmt.Println(apikey)
		elevenlabKey = currentData.ElevenlabKey
		str := "true"
		activateElevenlab = str == currentData.EnableEleven
		//fmt.Println("activateelevenlab:", activateElevenlab)
		//if a name is set in the config.json, it will auitomatically pick it
		if currentData.DefaultName != "" {
			nomiName = currentData.DefaultName
			updateId(nomiName)
		} else {
			fmt.Println("\033[31mits most likely that the config is malformed. Please delete it and re-run the program\033[0m")
			return
		}

	} else {
		regenerateConfig()
	}

	startChatting()
}

// function generateConfig generates the config file
func regenerateConfig() {
	fmt.Println("please paste here your api key to generate your config:")
	fmt.Scan(&apikey)
	//fmt.Println(apikey)
	var activateElevenlabString string
	fmt.Println("do you wanna add an elevenlab api key? (so your nomi can have a voice) (y/n)")
	fmt.Println("\033[31mYou \033[1;4mNEED\033[0m \033[31mto install mpv if you want this to work\033[0m")

	for {
		fmt.Scan(&activateElevenlabString)
		switch strings.ToLower(activateElevenlabString) {
		case "y", "yes":
			isIt := isMpvInstalled()

			if !isIt {
				fmt.Println("mpv isnt detected. please ensure that it is on your PATH")
				regenerateConfig()
			}
			activateElevenlab = true
			fmt.Println("paste your api key here:")
			fmt.Scan(&elevenlabKey)
			fmt.Println("do you to use the default voice?(the default one is cgSgspJ2msm6clMCkdW9) (y/n) ")
			var voiceId string
			for {
				fmt.Scan(&voiceId)
				switch strings.ToLower(voiceId) {
				case "y", "yes":
					fmt.Println("the default voice id is set")
					voiceId2 = "cgSgspJ2msm6clMCkdW9"
					break
				case "n", "no", "":

					fmt.Println("enter the voice id here:")
					fmt.Scan(&voiceId2)
					break
				default:
					fmt.Println("please input y (yes) or n (no)")
					continue
				}
				break
			}
			break
		case "n", "no", "":
			activateElevenlab = false
			break
		default:
			fmt.Println("please input y (yes) or n (no)")
			continue
		}
		break
	}

	for {
		fmt.Println("please pick a default nomi (by their name) to chat with (you can change this later, type :h to know more):")
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
func generateConfig() {
	_, finalString := listAndValidate("", 2)
	//fmt.Println(finalString)

	var towrite string
	if activateElevenlab {
		towrite = fmt.Sprintf(`{
	"apiKey": "%s",
   "default": "%s",
   "elevenlab": "%s",
   "activateVoice": "%t",
   "voiceid": "%s",
	"nomis": [
    %s
    ]
}`, apikey, nomiName, elevenlabKey, activateElevenlab, voiceId2, finalString)
	} else {
		towrite = fmt.Sprintf(`{
        "apiKey": "%s",
		"default": "%s",
		"elevenlab": "none",
        "activateVoice": "%t",
        "voiceid": "%s",
		"nomis": [
    %s
    ]
}`, apikey, nomiName, activateElevenlab, voiceId2, finalString)
	}
	os.WriteFile(filePath, []byte(towrite), 0644)
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
	//fmt.Println(string(b))
	var data ReplyText
	err = json.Unmarshal(b, &data)
	if err != nil {
		fmt.Printf("\033[31merror unmarshaling reply data:%s\033[31m\n", err)
		return ""
	}
	//fmt.Println("activate elevenlab:", activateElevenlab)
	if activateElevenlab {

		getAudioUwu(data.ReplyMessage.Text)
		exec.Command("mpv", audioFile)
	}
	//ai response
	return fmt.Sprintf("\033[36m%s> %s\033[0m", nomiName, data.ReplyMessage.Text)
}

// help message

func showHelp() {
	fmt.Println("Available Commands:")
	fmt.Println("  :h, :help         - Show this help message")
	fmt.Println("  :p, :pchange      - Change the default chat Nomi")
	fmt.Println("  :q, :quit         - Quit the chat session")
	fmt.Println("  :c, :change       - Change the current chat Nomi")
	fmt.Println("  :r, :restart      - Regenerate the configuration from scratch")
	//fmt.Println("  :t, :theme        - Change the color scheme")
	fmt.Println("  :de, :deactivateeleven - Deactivate the Eleven Labs integration")
	fmt.Println("  :ea, :actEleven   - Activate the Eleven Labs integration")
	fmt.Println("  :ae, :addElevenKey - Add your Eleven Labs API key")
	fmt.Println("  :re, :reload      - Reload the application (in case you manualy changed the config.json)")
	fmt.Println("  :lv, :listvoices  -  List aviable elevenlab voices (require api key)")
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
		fmt.Println("this line will never execute. if it does, i dont know wtf you did honestly (uwu)")
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
func getAudioUwu(input string) {
	if fileExists(audioFile) {
		err := os.Remove(audioFile)
		if err != nil {
			fmt.Println("Error clearing audio file", err)
		}
	}
	stopSpinner := make(chan bool)
	go func() {
		spinChars := `|/-\`
		i := 0
		for {
			select {
			case <-stopSpinner:
				return
			default:
				fmt.Printf("\033[95m\rfetching elevenlab audio... %c\033[0m", spinChars[i%len(spinChars)])
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	elevenlab.GetAudio(elevenlabKey, voiceId2, audioFile, input)
	stopSpinner <- true
	fmt.Print("\r\033[K")
}
func isMpvInstalled() bool {
	cmd := exec.Command("mpv", "--version")
	err := cmd.Run()
	if err != nil {
		fmt.Println("error running mpv:", err)
	}
	return err == nil
}
