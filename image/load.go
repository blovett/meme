package image

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/data"
	"github.com/nomad-software/meme/image/stream"
	"github.com/nomad-software/meme/output"
)

var (
	imageMap = make(map[string]string)
)

// Initialise the package. The embedded image assets are baked into the
// binary at build time, so a failure to read them here indicates a broken
// build rather than something a caller can recover from at runtime.
func init() {
	images, err := data.Files.ReadDir(data.ImagePath)
	if err != nil {
		panic(fmt.Errorf("could not read embedded images: %w", err))
	}

	for _, image := range images {
		id := strings.TrimSuffix(filepath.Base(image.Name()), data.ImageExtension)
		imageMap[id] = path.Join(data.ImagePath, image.Name())
	}
}

// Load an image from the passed string or stdin.
// The string will be a embedded asset id, an image URL or a local file.
func Load(opt cli.Options) (stream.Stream, error) {
	var s io.Reader
	var err error

	switch {
	case isURL(opt.Image):
		s, err = downloadURL(opt.Image)

	case isStdin(opt.Image):
		s, err = readStdin()

	case isAsset(opt.Image):
		s, err = loadAsset(opt.Image)

	default:
		var local bool
		local, err = isLocalFile(opt.Image)
		if err == nil {
			if local {
				s, err = readFile(opt.Image)
			} else {
				err = output.Errorf("image not recognised")
			}
		}
	}

	if err != nil {
		return stream.Stream{}, err
	}

	return stream.NewStream(s)
}

// Return true if the passed string is an embedded asset id, false if not.
func isAsset(id string) bool {
	_, ok := imageMap[id]
	return ok
}

// Load and return an embedded asset (image) by id.
// The id is assumed to exist.
func loadAsset(id string) (io.Reader, error) {
	image := imageMap[id]

	st, err := data.Files.ReadFile(image)
	if err != nil {
		return nil, output.WrapError(err, "could not read embedded image")
	}

	return bytes.NewReader(st), nil
}

// Return true if the passed string is an image URL, false if not.
func isURL(url string) bool {
	return strings.HasPrefix(url, "http")
}

// Download the image located at the passed image URL, decode and return it.
func downloadURL(url string) (io.Reader, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, output.WrapError(err, "request error")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, output.Errorf("could not access URL")
	}

	st, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, output.WrapError(err, "could not read response body")
	}

	return bytes.NewReader(st), nil
}

// Return true if the passed string is a file that exists on the local
// filesystem, false if not.
func isLocalFile(path string) (bool, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return false, output.WrapError(err, "could not expand path")
	}

	_, err = os.Stat(path)
	return err == nil, nil
}

// Read and return a file on the local filesystem.
// The file is assumed to exist.
func readFile(path string) (io.Reader, error) {
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, output.WrapError(err, "could not expand path")
	}

	st, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, output.WrapError(err, "could not read local file")
	}

	return bytes.NewReader(st), nil
}

// return true if the passed string is '-' meaning we should read the image
// from stdin.
func isStdin(path string) bool {
	return path == "-"
}

// Read the image from stdin.
func readStdin() (io.Reader, error) {
	st, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, output.WrapError(err, "could not read stdin")
	}

	return bytes.NewReader(st), nil
}

// Decal returns the named decal as a stream.
func Decal(name string) (stream.Stream, error) {
	st, err := data.Files.ReadFile(name)
	if err != nil {
		return stream.Stream{}, output.WrapError(err, "could not read embedded decal")
	}

	return stream.NewStream(bytes.NewReader(st))
}
