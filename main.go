package main

import (
	"io/ioutil"
	"log"
	"net"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jjexpat00/spoop/mopidy"
)

func main() {
	key, err := ioutil.ReadFile("key.conf")
	if err != nil {
		log.Fatal(err)
	}

	mop := mopidy.New("localhost:6680")
	panda := mopidy.Track{
		Model:   "Track",
		Name:    "panda",
		Uri:     "spotify:track:275a9yzwGB6ncAW4SxY7q3",
		Length:  248,
		TrackNo: 0,
	}

	tracks := make([]mopidy.Track, 1)
	tracks[0] = panda

	err = mop.AddTracks(tracks)
	if err != nil {
		log.Fatal(err)
	}

	err = mop.Play()
	if err != nil {
		log.Fatal(err)
	}

	err = mop.Tracks()
	if err != nil {
		log.Fatal(err)
	}

	err = connectToDiscord(string(key))
	if err != nil {
		log.Fatal(err)
	}

	for {
	}
}

func connectToDiscord(key string) error {
	discord, err := discordgo.New("Bot " + strings.Trim(key, "\n"))
	if err != nil {
		return err
	}

	err = discord.Open()
	if err != nil {
		return err
	}

	_, err = discord.ChannelMessageSendTTS("361963917179617304", "suh dude")
	if err != nil {
		return err
	}

	vc, err := discord.ChannelVoiceJoin("361963917179617301", "361963917183811584", false, false)
	if err != nil {
		return err
	}

	sayShit(vc)

	return nil
}

func sayShit(vc *discordgo.VoiceConnection) {
	mopidyStream, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.ListenUDP("udp", mopidyStream)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	buffer := make([]byte, 1024)

	for true {
		bytesRead, err := connection.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("bytes read %v", bytesRead)
		vc.OpusSend <- buffer[:bytesRead]
	}
}
