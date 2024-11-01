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

func GetApiKey() ***REMOVED***

***REMOVED***
func GetAudio(XI_API_KEY string, VOICE_ID string, OUTPUT_PATH string, TEXT_TO_SPEAK string) ***REMOVED***
	//fmt.Println(XI_API_KEY)
	ttsURL := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s/stream", VOICE_ID)

	headers := map[string]string***REMOVED***
		"Accept":     "application/json",
		"xi-api-key": XI_API_KEY,
	***REMOVED***

	// Set up the data payload for the API request
	data := map[string]interface***REMOVED******REMOVED******REMOVED***
		"text":     TEXT_TO_SPEAK,
		"model_id": "eleven_multilingual_v2",
		"voice_settings": map[string]interface***REMOVED******REMOVED******REMOVED***
			"stability":         0.5,
			"similarity_boost":  0.8,
			"style":             0.0,
			"use_speaker_boost": true,
		***REMOVED***,
	***REMOVED***

	jsonData, err := json.Marshal(data)
	if err != nil ***REMOVED***
		fmt.Println("Error marshalling JSON:", err)
		return
	***REMOVED***

	req, err := http.NewRequest("POST", ttsURL, bytes.NewBuffer(jsonData))
	if err != nil ***REMOVED***
		fmt.Println("Error creating request:", err)
		return
	***REMOVED***

	// Set headers
	for key, value := range headers ***REMOVED***
		req.Header.Set(key, value)
	***REMOVED***
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client***REMOVED******REMOVED***

	resp, err := client.Do(req)
	if err != nil ***REMOVED***
		fmt.Println("Error sending request:", err)
		return
	***REMOVED***
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode == http.StatusOK ***REMOVED***
		// Open the output file in write-binary mode
		file, err := os.Create(OUTPUT_PATH)
		if err != nil ***REMOVED***
			fmt.Println("Error creating output file:", err)
			return
		***REMOVED***
		defer file.Close()

		// Read the response in chunks and write to the file
		buf := make([]byte, CHUNK_SIZE)
		for ***REMOVED***
			n, err := resp.Body.Read(buf)
			if err != nil && err != io.EOF ***REMOVED***
				fmt.Println("Error reading response body:", err)
				return
			***REMOVED***
			if n == 0 ***REMOVED***
				break
			***REMOVED***
			if _, err := file.Write(buf[:n]); err != nil ***REMOVED***
				fmt.Println("Error writing to file:", err)
				return
			***REMOVED***
		***REMOVED***
		// Inform the user of success
		//fmt.Println("Audio stream saved successfully.")
	***REMOVED*** else ***REMOVED***
		// Print the error message if the request was not successful
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Error:", resp.Status, string(body))
	***REMOVED***
***REMOVED***
