// Copyright 2017 Albert Nigmatzianov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package id3v2

type Options struct {
	// Parse the ID3v2 tag from a file (true) or create an empty tag (false).
	Parse bool

	// Define specific frames to parse (unless parse is disabled).
	//
	// For example, `[]string{"Artist", "Title"}`
	// will only parse artist and title frames.
	// IDs (e. g., "TPE1", "TIT2") can be specified as well.
	//
	// If the slice is blank or nil, all frames in the tag will be parsed.
	ParseFrames []string
}
