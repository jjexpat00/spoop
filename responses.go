package mopidy

type BasicResponse struct {
	Error  `json:"error"`
	Result interface{} `json:"result"`
}

type TrackResponse struct {
	Error  `json:"error"`
	Result *Track `json:"result"`
}

type TracksResponse struct {
	Error  `json:"error"`
	Result []Track `json:"result"`
}
