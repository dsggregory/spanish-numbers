package langpractice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// LPResponse HTTP response from langpractice.com
type LPResponse struct {
	// AudioData speech for the spanish number in MPEG ADTS, layer III, v2, 160 kbps, 24 kHz, Monaural.
	//		Before JSonUnmarshal() it is base64-encoded.
	AudioData []byte `json:"audio_data"` // json decoder removes base64 b/c of type []byte
	// Number the random number
	Number int `json:"n"`
	// Target the spelling of the number
	Target struct {
		Phonetic string `json:"phonetic"`
		Written  string `json:"written"`
	} `json:"target"`
}

func (c *LangPractice) parseResponse(body io.Reader) (*LPResponse, error) {
	// decode the JSON response
	dec := json.NewDecoder(body)
	res := LPResponse{}
	if err := dec.Decode(&res); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &res, nil
}

// PlayResponse read the response from the LP API and play the audio
func (c *LangPractice) PlayResponse(res *LPResponse) error {
	if err := c.Play(bytes.NewReader(res.AudioData)); err != nil {
		return fmt.Errorf("playing audio: %w", err)
	}

	return nil
}
