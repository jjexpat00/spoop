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

var dialer = websocket.Dialer{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	url  string
	rpc  string
	conn *websocket.Conn
}

type Track struct {
	Model   string `json:"__model__"`
	Name    string `json:"name"`
	Uri     string `json:"uri"`
	Length  int    `json:"length"`
	TrackNo int    `json:"track_no"`
}

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
