package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"github.com/nats-io/nats.go"
)

var musicDirectory = ``

func tempDB(title string) string {

	music := make(map[string]string)

	music["cowbell warrior"] = fmt.Sprintf("%sSXMPRA - COWBELL WARRIOR!.mp3", musicDirectory)

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
			log.Printf("System info: Play music %v", msgResult)

			if msgResult != "" {

				input := msgResult[11:]
				log.Printf("System info: Play %v", input)
				url := tempDB(input)
				f, _ := os.Open(url)

				defer f.Close()

				d, _ := mp3.NewDecoder(f)

				c, ready, _ := oto.NewContext(d.SampleRate(), 2, 2)

				<-ready

				p := c.NewPlayer(d)
				defer p.Close()
				p.Play()

				fmt.Printf("Length: %d[bytes]\n", d.Length())
				for {
					time.Sleep(time.Second)
					if !p.IsPlaying() {
						break
					}
				}
			}
		}

	})

	log.Println("Subscribed to", topic)

	<-wait
}
