package elevenlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	CHUNK_SIZE = 1024
)

func GetApiKey() {

}
func GetAudio(XI_API_KEY string, VOICE_ID string, OUTPUT_PATH string, TEXT_TO_SPEAK string) {
	//fmt.Println(XI_API_KEY)
	ttsURL := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s/stream", VOICE_ID)

	headers := map[string]string{
		"Accept":     "application/json",
		"xi-api-key": XI_API_KEY,
	}

	// Set up the data payload for the API request
	data := map[string]interface{}{
		"text":     TEXT_TO_SPEAK,
		"model_id": "eleven_multilingual_v2",
		"voice_settings": map[string]interface{}{
			"stability":         0.5,
			"similarity_boost":  0.8,
			"style":             0.0,
			"use_speaker_boost": true,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", ttsURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode == http.StatusOK {
		// Open the output file in write-binary mode
		file, err := os.Create(OUTPUT_PATH)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			return
		}
		defer file.Close()

		// Read the response in chunks and write to the file
		buf := make([]byte, CHUNK_SIZE)
		for {
			n, err := resp.Body.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("Error reading response body:", err)
				return
			}
			if n == 0 {
				break
			}
			if _, err := file.Write(buf[:n]); err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
		// Inform the user of success
		//fmt.Println("Audio stream saved successfully.")
	} else {
		// Print the error message if the request was not successful
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Error:", resp.Status, string(body))
	}
}
