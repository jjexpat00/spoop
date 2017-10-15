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


	// Here we read the Discord bot's token.
	key, err := ioutil.ReadFile("key.conf")
	if err != nil {
		log.Fatal(err)
	}


	// Mopi is set as the device Mopidy runs through. We also
	// create SpotifyTestSong, where the URI is the most important
	// song identifier at the moment.
	Mopi := mopidy.New("localhost:6680")
	SpotifyTestSong := mopidy.Track{
		Model:   "Track",
		Name:    "First Love (from the film)",
		Uri:     "spotify:track:30Zy2WFNjl4pUIkEocjxdN",
		Length:  195,
		TrackNo: 0,
	}


	// tracks forms the queue in which Track are loaded into, and
	// in this case, SpotifyTestSong is loaded as the first Track.
	tracks := make([]mopidy.Track, 1)
	tracks[0] = SpotifyTestSong

	err = Mopi.AddTracks(tracks)
	if err != nil {
		log.Fatal(err)
	}

	err = Mopi.Play()
	if err != nil {
		log.Fatal(err)
	}

	err = Mopi.Tracks()
	if err != nil {
		log.Fatal(err)
	}

	err = connectToDiscord(string(key))
	if err != nil {
		log.Fatal(err)
	}


	// This keeps main persistent?
	for {
	}
}


// For now, spoop automatically writes out a message and joins
// the General channel (associated with Channel ID).
func connectToDiscord(key string) error {

	// For some reason, we need to trim the first line of the key file.
	discord, err := discordgo.New("Bot " + strings.Trim(key, "\n"))
	if err != nil {
		return err
	}

	err = discord.Open()
	if err != nil {
		return err
	}

	// Sends a message on joining.
	_, err = discord.ChannelMessageSendTTS("361963917179617304", "suh dude")
	if err != nil {
		return err
	}

	// Joins a preset channel and its respective server. Syntax
	// for ChannelVoiceJoin is (gID, cID, mute, deaf).
	vc, err := discord.ChannelVoiceJoin("361963917179617301", "361963917183811584", false, false)
	if err != nil {
		return err
	}


	// Upon joining the channel, spoop begins playing the preset
	// URI in Track. Not sure if it's proper, but I believe some
	// if statements are warranted here. I'm also not sure if 
	// commands (eg. !pause, !stop, etc) should be added here.
	musicStream(vc)

	return nil
}


// Working on responses later...

//func reponseObjectForCommand(command string) interface{} {
//	mapping := map[string]interface{}{
//		"core.playback.get_current_track": &TrackResponse{},
//		"core.tracklist.get_tracks":       &TracksResponse{},
//	}
//
//	obj := mapping[command]
//	if obj == nil {
//		obj = &BasicResponse{}
//	}
//
//	return obj
//}


// Here is the core of the UDP sink. We feed musicStream into
// the voice channel spoop joins. Please excuse the port number.
func musicStream(vc *discordgo.VoiceConnection) {
	mopidyStream, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.ListenUDP("udp", mopidyStream)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()


	// Still not sure how to optimize buffer, will work on
	// that later to see if the stuttering kink can be fixed.
	buffer := make([]byte, 160)

	for true {
		bytesRead, err := connection.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}


		// This just displays activity on the server.
		log.Printf("Bytes read %v", bytesRead)


		// Here we send the ingested buffer via OpusSend.
		vc.OpusSend <- buffer[:bytesRead]
	}
}
