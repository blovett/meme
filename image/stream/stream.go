package stream

import (
	"bytes"
	"image"
	"image/gif"
	"image/png"
	"io"
	"io/ioutil"

	"github.com/nomad-software/meme/output"
)

// Stream contains information about a loaded image.
type Stream struct {
	io.Reader
	bytes []byte
	index int
	typ   string
}

// Bytes returns the stream's bytes.
func (st *Stream) Bytes() []byte {
	return st.bytes
}

// Read implements the io.Reader interface for the Stream.
func (st *Stream) Read(b []byte) (n int, err error) {
	if st.index >= len(st.bytes) {
		return 0, io.EOF
	}
	n = copy(b, st.bytes[st.index:])
	st.index += n
	return
}

// IsGif returns true if the loaded image is a gif.
func (st *Stream) IsGif() bool {
	return st.typ == "gif"
}

// IsJpg returns true if the loaded image is a Jpeg.
func (st *Stream) IsJpg() bool {
	return st.typ == "jpeg"
}

// IsPng returns true if the loaded image is a Png.
func (st *Stream) IsPng() bool {
	return st.typ == "png"
}

// FileExt returns the file extension of the image.
func (st *Stream) FileExt() string {
	if st.IsGif() {
		return "gif"
	}
	if st.IsJpg() {
		return "jpg"
	}
	if st.IsPng() {
		return "png"
	}
	panic("File extension not recognised")
}

// NewStream creates a new stream. It returns an error if stream is not
// readable or does not contain a recognisable image, which can legitimately
// happen for untrusted input (a corrupt upload, a URL that doesn't point at
// an image, etc), so callers must handle it rather than the process exiting.
func NewStream(stream io.Reader) (Stream, error) {
	a, err := ioutil.ReadAll(stream)
	if err != nil {
		return Stream{}, output.WrapError(err, "could not read image bytes")
	}

	b := make([]byte, len(a))
	copy(b, a)

	_, typ, err := image.DecodeConfig(bytes.NewReader(a))
	if err != nil {
		return Stream{}, output.WrapError(err, "could not decode image config")
	}

	return Stream{
		bytes: b,
		typ:   typ,
	}, nil
}

// EncodeImage encodes an image into a stream.
func EncodeImage(img image.Image) (Stream, error) {
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		return Stream{}, output.WrapError(err, "could not encode image")
	}
	return NewStream(&buffer)
}

// DecodeImage decodes the byte stream and returns an image.
func (st *Stream) DecodeImage() (image.Image, error) {
	img, _, err := image.Decode(st)
	if err != nil {
		return nil, output.WrapError(err, "could not decode image")
	}
	return img, nil
}

// EncodeGif encodes a gif into a stream.
func EncodeGif(img *gif.GIF) (Stream, error) {
	var buffer bytes.Buffer
	if err := gif.EncodeAll(&buffer, img); err != nil {
		return Stream{}, output.WrapError(err, "could not encode gif")
	}
	return NewStream(&buffer)
}

// DecodeGif decodes the byte stream and returns a gif.
func (st *Stream) DecodeGif() (*gif.GIF, error) {
	if !st.IsGif() {
		return nil, output.Errorf("can't decode stream to gif")
	}
	g, err := gif.DecodeAll(st)
	if err != nil {
		return nil, output.WrapError(err, "could not decode gif")
	}
	return g, nil
}
