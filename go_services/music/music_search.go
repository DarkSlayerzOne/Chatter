package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
	"github.com/nats-io/nats.go"
)

type Music struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Path      string    `json:"path"`
	Keyword   string    `json:"keyWord"`
	Lyrics    string    `json:"lyrics"`
	Genre     string    `json:"genre"`
	PlayGroup string    `json:"playGroup"`
	Artist    string    `json:"artist"`
	DateAdded time.Time `json:"dateAdded"`
}

const (
	Main_Topic = "send_message"

	Play_Music_Voice_Keyword        = "play music"
	Play_All_Voice_Keyword          = "play my list"
	Play_ALL_Randomly_Voice_Keyword = "play any music"
)

var musicDirectory = `Directory here`

func main() {
	wait := make(chan bool)

	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal(err)
	}

	nc.Subscribe(Main_Topic, func(m *nats.Msg) {
		log.Printf("Received a message: %s\n", string(m.Data))

		msgResult := string(m.Data)

		if strings.Contains(msgResult, Play_Music_Voice_Keyword) {
			log.Printf("System info: Play music %v", msgResult)

			playMusicTitle(msgResult)
		}

		if strings.Contains(msgResult, Play_All_Voice_Keyword) {
			log.Printf("System info: %v", msgResult)

			playAll()
		}

		if strings.Contains(msgResult, Play_ALL_Randomly_Voice_Keyword) {
			log.Printf("System info: %v", msgResult)

			playRandomMusic()
		}

	})

	log.Println("Subscribed to", Main_Topic)

	<-wait
}

func playMusicTitle(musicKeyword string) {

	if musicKeyword != "" {

		input := musicKeyword[11:]
		log.Printf("System info: Search for %v", input)
		musicList := tempDB()

		var musicToPlay string = ""
		var title string = ""

		for _, song := range musicList {

			if strings.Contains(song.Keyword, strings.ToLower(input)) {
				log.Printf(`System info: Music found "%s" song by "%s" genre "%s"`, song.Title, song.Artist, song.Genre)
				title = song.Title
				musicToPlay = song.Path
			}
		}

		if musicToPlay == "" {
			return
		}

		isPlaying := musicPlayer(musicToPlay)

		if !isPlaying {
			log.Printf("System info: Music %s is now finished playing", title)
		}

	}
}

func playAll() {

	musicList := tempDB()

	for _, music := range musicList {

		log.Printf(`System info: Now playing "%s" song by "%s"`, music.Title, music.Artist)

		isPlaying := musicPlayer(music.Path)

		if isPlaying {
			break
		}
	}
}

func playRandomMusic() {

	rand.Seed(time.Now().Unix())

	musicList := tempDB()

	for i := range musicList {
		j := rand.Intn(i + 1)
		musicList[i], musicList[j] = musicList[j], musicList[i]
	}

	for _, music := range musicList {

		if music.Title == "" {
			log.Printf("System info: No music was found")
			return
		}

		log.Printf(`System info: Now playing "%s" song by "%s"`, music.Title, music.Artist)

		isPlaying := musicPlayer(music.Path)

		if isPlaying {
			break
		}
	}
}

func musicPlayer(musicPlay string) bool {

	var isPlaying bool = true

	f, _ := os.Open(musicPlay)

	defer f.Close()

	d, _ := mp3.NewDecoder(f)

	c, ready, _ := oto.NewContext(d.SampleRate(), 2, 2)

	<-ready

	p := c.NewPlayer(d)
	defer p.Close()
	p.Play()

	log.Printf("System info : Length %d[bytes]\n", d.Length())
	for {
		time.Sleep(time.Second)
		if !p.IsPlaying() {
			isPlaying = false
			break
		}
	}

	return isPlaying
}

// TODO: use a cloud storage like s3 and a NOSQLDB
func tempDB() []Music {

	return []Music{
		{ID: 1, Title: "Cowbell Warrior", Path: fmt.Sprintf("%sSXMPRA - COWBELL WARRIOR!.mp3", musicDirectory), Keyword: "cowbell warrior", Lyrics: "", Genre: "Rap", PlayGroup: "", Artist: "SXMPRA", DateAdded: time.Now()},
		{ID: 2, Title: "Green Green Grass", Path: fmt.Sprintf("%sGeorge Ezra - Green Green Grass.mp3", musicDirectory), Keyword: "green green grass", Lyrics: "", Genre: "Pop", PlayGroup: "", Artist: "George Ezra", DateAdded: time.Now()},
	}
}
