// Copyright 2026 Zmicer Pasternak. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v1

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	id3v1NoTrackNumber byte = 1
	id3v1HeaderSize         = 128
	id3v1Marker             = "TAG"
)

var id3v1EmptyString = string(make([]byte, 30))

// ID3v1 data format has fixed size of 128 bytes.

type Id3v1 struct {
	rest     []byte // the rest of file data (the music part)
	title    string // 30 bytes
	artist   string // 30 bytes
	album    string // 30 bytes
	year     int    // 4 bytes; a four-digit year
	comment  string // 28 bytes (or 30 bytes taking up the `zeroByte` & `track` bytes)
	zeroByte byte   // 1 byte; if the track number is stored, this byte contains a binary 0
	track    byte   // 1 byte; invalid, if the previous byte is not a binary 0
	genre    Genre  // 1 byte; index of a genre, or 255
}

func (id3v1 *Id3v1) String() string {
	b := &strings.Builder{}

	fmt.Fprintf(b, "Title: %v\n", id3v1.title)
	fmt.Fprintf(b, "Artist: %v\n", id3v1.artist)
	fmt.Fprintf(b, "Album: %v\n", id3v1.album)
	fmt.Fprintf(b, "Year: %v\n", id3v1.year)
	fmt.Fprintf(b, "Comment: %v\n", id3v1.comment)

	if id3v1.zeroByte == 0 {
		fmt.Fprintf(b, "TrackNumber: %v\n", id3v1.track)
	}

	fmt.Fprintf(b, "Genre: %v\n", id3v1.genre)

	return b.String()
}

func (id3v1 *Id3v1) Reset() {
	id3v1.title = id3v1EmptyString
	id3v1.artist = id3v1EmptyString
	id3v1.album = id3v1EmptyString
	id3v1.year = 0
	id3v1.comment = id3v1EmptyString
	id3v1.zeroByte = id3v1NoTrackNumber
	id3v1.track = 0
	id3v1.genre = Undefined
}

func (id3v1 *Id3v1) SaveTo(path string) error {
	var tempFile = path + ".tmp"

	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("could not create file %q: %w", tempFile, err)
	}
	defer file.Close()

	if _, err := file.Write(id3v1.rest); err != nil {
		return fmt.Errorf("could not write `rest`: %w", err)
	}

	if _, err := file.Write([]byte(id3v1Marker)); err != nil {
		return fmt.Errorf("could not write `id3v1Marker`: %w", err)
	}

	if err := writeString(file, id3v1.title, 30); err != nil {
		return fmt.Errorf("could not write `id3v1.title`: %w", err)
	}

	if err := writeString(file, id3v1.artist, 30); err != nil {
		return fmt.Errorf("could not write `id3v1.artist`: %w", err)
	}

	if err := writeString(file, id3v1.album, 30); err != nil {
		return fmt.Errorf("could not write `id3v1.album`: %w", err)
	}

	if err := writeString(file, strconv.Itoa(id3v1.year), 4); err != nil {
		return fmt.Errorf("could not write `id3v1.year`: %w", err)
	}

	if id3v1.zeroByte != 0 {
		if len(id3v1.comment) > 28 {
			if err := writeString(file, id3v1.comment, 30); err != nil {
				return fmt.Errorf("could not write `id3v1.comment` (30 bytes): %w", err)
			}
		} else {
			if err := writeString(file, id3v1.comment, 28); err != nil {
				return fmt.Errorf("could not write `id3v1.comment` (28 bytes): %w", err)
			}
			if _, err := file.Write([]byte{1, 0}); err != nil {
				return fmt.Errorf("could not write `id3v1.zeroByte(1)` & `id3v1.track(0)`: %w", err)
			}
		}
	} else {
		if err := writeString(file, id3v1.comment, 28); err != nil {
			return fmt.Errorf("could not write `id3v1.comment` (28 bytes): %w", err)
		}
		if _, err := file.Write([]byte{0, id3v1.track}); err != nil {
			return fmt.Errorf("could not write `id3v1.track`: %w", err)
		}
	}

	if _, err := file.Write([]byte{byte(id3v1.genre)}); err != nil {
		return fmt.Errorf("could not write `id3v1.genre`: %w", err)
	}

	if err := os.Rename(tempFile, path); err != nil {
		// If rename fails, clean up the temporary file.
		os.Remove(tempFile)
		return fmt.Errorf("error renaming file %q: %w", tempFile, err)
	}

	return nil
}

// func (id3v1 *Id3v1) SetTitle(title string) {}
// func (id3v1 *Id3v1) SetArtist(artist string) {}
// func (id3v1 *Id3v1) SetAlbum(album string) {}
// func (id3v1 *Id3v1) SetYear(year int) {}
// func (id3v1 *Id3v1) SetComment(comment string) {}
// func (id3v1 *Id3v1) SetTrackNumber(track byte) {}
// func (id3v1 *Id3v1) SetGenre(genre Genre) {}

func (id3v1 *Id3v1) parse(input io.ReadSeeker, opts Options) error {
	if !opts.Parse {
		id3v1.Reset()
		return id3v1.writeMusicPart(input)
	}

	byteSeq, err := seekAndRead(input, -id3v1HeaderSize, io.SeekEnd, id3v1HeaderSize)
	if err != nil {
		return fmt.Errorf("could not read file header: %w", err)
	}

	id3v1.title = stringBeforeZero(byteSeq[3:33])
	id3v1.artist = stringBeforeZero(byteSeq[33:63])
	id3v1.album = stringBeforeZero(byteSeq[63:93])
	id3v1.year, _ = strconv.Atoi(string(byteSeq[93:97]))

	// The track number is stored in the last two bytes of the comment field.
	// If the comment is 29 or 30 characters long, no track number can be stored.
	if byteSeq[125] == 0 {
		id3v1.comment = stringBeforeZero(byteSeq[97:125])
		id3v1.zeroByte = 0
		id3v1.track = byteSeq[126]
	} else {
		id3v1.comment = stringBeforeZero(byteSeq[97:127])
		id3v1.zeroByte = byteSeq[125]
		id3v1.track = 0
	}

	// Index in a list of genres, or 255.
	id3v1.genre = Genre(byteSeq[127])

	// Read another file data.
	if _, err := input.Seek(0, io.SeekStart); err != nil {
		return err
	}

	return id3v1.writeMusicPart(input)
}

func (id3v1 *Id3v1) writeMusicPart(input io.ReadSeeker) error {
	data, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("could not read music data: %w", err)
	}
	id3v1.rest = data[:len(data)-id3v1HeaderSize]

	return nil
}

// Opens a file and returns its ID3v1 tag.
func Open(filePath string, opts Options) (*Id3v1, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Parse(file, opts)
}

// Parses `input` and returns its ID3v1 tag.
func Parse(input io.ReadSeeker, opts Options) (*Id3v1, error) {
	tag := &Id3v1{}
	err := tag.parse(input, opts)

	return tag, err
}

func seekAndRead(input io.ReadSeeker, offset int64, whence int, read int) ([]byte, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	if _, err := input.Seek(offset, whence); err != nil {
		return nil, err
	}

	data := make([]byte, read)
	bytesRead, err := input.Read(data)
	if err != nil {
		return nil, err
	}
	if bytesRead != read {
		return nil, fmt.Errorf("read is incomplete: %v/%v bytes", bytesRead, read)
	}

	return data, nil
}

func stringBeforeZero(data []byte) string {
	before, _, ok := bytes.Cut(data, []byte{0})
	if !ok {
		return string(data)
	}
	return string(before)
}

func writeString(input io.Writer, data string, size int) error {
	if dataLen := len(data); dataLen > size {
		return fmt.Errorf("data to be written is too long: %v", dataLen)
	}

	byteStr := make([]byte, size)
	for i, val := range data {
		byteStr[i] = byte(val)
	}

	n, err := input.Write(byteStr)
	if err != nil {
		return err
	}
	if n != size {
		return fmt.Errorf("write is incomplete: %v/%v bytes", n, size)
	}

	return nil
}
