package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/nats-io/nats.go"
)

func tempDB(title string) string {

	music := make(map[string]string)

	music["timber"] = "https://www.youtube.com/watch?v=hHUbLv4ThOo"
	music["sweet but psycho"] = "https://www.youtube.com/watch?v=WXBHCQYxwr0"
	music["somebody that i used to know"] = "https://www.youtube.com/watch?v=8UVNT4wvIGY"
	music["take on me"] = "https://www.youtube.com/watch?v=djV11Xbc914"

	return music[strings.ToLower(title)]
}

func main() {

	wait := make(chan bool)

	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}
	const topic = "send_message"

	nc.Subscribe(topic, func(m *nats.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))

		msgResult := string(m.Data)

		const keyWord = "play music"

		if strings.Contains(msgResult, keyWord) {

			input := msgResult[11:]
			url := tempDB(input)
			var err error

			switch runtime.GOOS {
			case "linux":
				err = exec.Command("xdg-open", url).Start()
			case "windows":
				err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
			case "darwin":
				err = exec.Command("open", url).Start()
			default:
				err = fmt.Errorf("unsupported platform")
			}
			if err != nil {
				log.Fatal(err)
			}
		}

	})

	log.Println("Subscribed to", topic)

	<-wait

}
