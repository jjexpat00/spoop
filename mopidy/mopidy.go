package mopidy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	rpc "github.com/gorilla/rpc/v2/json2"
	"github.com/gorilla/websocket"
)


// Create dialer.
var dialer = websocket.Dialer{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


// Create Client structure.
type Client struct {
	url  string
	rpc  string
	conn *websocket.Conn
}


// Create Track structure.
type Track struct {
	Model   string `json:"__model__"`
	Name    string `json:"name"`
	Uri     string `json:"uri"`
	Length  int    `json:"length"`
	TrackNo int    `json:"track_no"`
}


// Feeds mopidy into Client.
func New(host string) *Client {
	return &Client{
		url: fmt.Sprintf("ws://%s/mopidy/ws", host),
		rpc: fmt.Sprintf("http://%s/mopidy/rpc", host),
	}
}

func (mopidy *Client) Connect() error {
	ws, _, err := dialer.Dial(mopidy.url, nil)
	if err != nil {
		return err
	}

	mopidy.conn = ws
	return nil
}

func (m *Client) Call(command string, params interface{}) error {
	if params == nil {
		params = map[string]string{}
	}

	buff, err := rpc.EncodeClientRequest(command, params)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(buff)
	resp, err := http.Post(m.rpc, "application/json", reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Println(string(bytes))

	return nil
}


// Add Client functionality. For now, everything works
// as a preset playing one song from Spotify.
func (m *Client) AddTracks(tracks []Track) error {
	params := map[string][]Track{"tracks": tracks}
	err := m.Call("core.tracklist.add", params)
	return err
}

func (m *Client) Play() error {
	err := m.Call("core.playback.play", nil)
	return err
}

func (m *Client) Tracks() error {
	err := m.Call("core.playback.get_current_track", nil)
	return err
}


// Will work on these in main, probably with if statements...
func (m *Client) Pause() error {
	err := m.Call("core.playback.pause", nil)
	return err
}

func (m *Client) Resume() error {
	err := m.Call("core.playback.resume", nil)
	return err
}
