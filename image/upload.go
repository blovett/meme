package image

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/image/stream"
	"github.com/nomad-software/meme/output"
)

const (
	uploadURL = "https://api.imgur.com/3/upload"
)

// Upload the image.
func Upload(opt cli.Options, st stream.Stream) (string, error) {
	base64 := base64.StdEncoding.EncodeToString(st.Bytes())
	return upload(opt, base64)
}

// Perform the request to the storage provider.
func upload(opt cli.Options, base64 string) (string, error) {
	req, err := http.NewRequest("POST", uploadURL, strings.NewReader(base64))
	if err != nil {
		return "", output.WrapError(err, "could not create upload request")
	}
	req.Header.Set("Authorization", "Client-ID "+opt.ClientID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", output.WrapError(err, "could not upload image")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", output.WrapError(err, "could not read response body")
		}

		var imgur imgurResponse
		if err := json.Unmarshal(body, &imgur); err != nil {
			return "", output.WrapError(err, "could not decode json response")
		}

		return imgur.Data.Link, nil
	}

	return "", output.Errorf("could not upload image")
}

type imgurResponse struct {
	Data imgurData `json:"data"`
}

type imgurData struct {
	Link string `json:"link"`
}
