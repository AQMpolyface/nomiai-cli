this is simple cli tool to talk to nomiai api.

Must have: a nomi api key. you can get one easely even with the free tier here: https://beta.nomi.ai/profile/integrations

Optional: an elevenlab api key. this will just allow the ai to have a voice. You can get one here for free: https://elevenlabs.io/app/speech-synthesis/text-to-speech

Features: an integration with elevenlab, and a shell to send and receive messages.

The nomi api is very limited, so there is no support for proactive messages, image or video for now. I will include them when the nomi team add that to their api

You can downloads the binary directly here: https://github.com/AQMpolyface/nomiai-cli/releases

To run it, just install the go programing language, and git then:


````
git clone https://github.com/AQMpolyface/nomiai-cli.git
cd nomiai-cli
go run nomi.go
````



the script will create a config.json in your current directory.
You will get prompted for your api key.

The nomi go sdk (https://github.com/vhalmd/nomi-go-sdk) helped me a lot to create some part of my code.

I would appreciate if you found a bug, or had an idea to make the tool better, please let me know via github. thanks!

You dont need to modify the code In any way, you will be prompted for your api key, and you will have to pick your Nomi by their names. (case sensitive)

You can downloads the binary directly here: https://github.com/AQMpolyface/nomiai-cli/releases
