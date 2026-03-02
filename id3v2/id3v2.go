// Copyright 2016 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package id3v2 is the ID3 parsing and writing Go library.
package id3v2

import (
	"io"
	"os"
)

// Available picture types for the picture frame.
const (
	PTOther = iota
	PTFileIcon
	PTOtherFileIcon
	PTFrontCover
	PTBackCover
	PTLeafletPage
	PTMedia
	PTLeadArtistSoloist
	PTArtistPerformer
	PTConductor
	PTBandOrchestra
	PTComposer
	PTLyricistTextWriter
	PTRecordingLocation
	PTDuringRecording
	PTDuringPerformance
	PTMovieScreenCapture
	PTBrightColouredFish
	PTIllustration
	PTBandArtistLogotype
	PTPublisherStudioLogotype
)

// Opens a file and parses its tag (if needed).
// If there is no tag in the file, creates a new ID3v2.4.
func Open(filePath string, opts Options) (*Tag, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return ParseReader(file, opts)
}

// Parses and finds the tag in the specified io.Reader considering the options.
// If there is no tag, creates a new ID3v2.4.
func ParseReader(rd io.Reader, opts Options) (*Tag, error) {
	tag := NewEmptyTag()
	err := tag.parse(rd, opts)
	return tag, err
}

// Returns an empty ID3v2.4 tag without any frames and an io.Reader.
func NewEmptyTag() *Tag {
	tag := new(Tag)
	tag.init(nil, 0, 4)
	return tag
}
