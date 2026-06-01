// Copyright 2026 Zmicer Pasternak. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"iter"
	"log"
	"strings"
)

const invalidFileFormatMsg = "invalid file extension; expected %q"

type frameRecord struct {
	frame string // e. g. "TIT2"
	descr string // e. g. "Title"
}

func (r frameRecord) String() string {
	return fmt.Sprintf("%v (%v)", r.frame, r.descr) // -> "TIT2 (Title)"
}

func scanInputLine(s *bufio.Scanner) (string, error) {
	if hasTokens := s.Scan(); !hasTokens {
		// Note: `Err()` returns `nil` on `io.EOF`.
		return "", s.Err()
	}
	if input := s.Text(); len(input) > 0 {
		return strings.TrimSpace(input), nil
	}
	return "", errors.New("input is empty")
}

func confirmEdit(s *bufio.Scanner, frame string) bool {
	for {
		fmt.Printf("%v: change? [Y/n]: ", frame)
		input, err := scanInputLine(s)
		if err != nil {
			log.Println(err)
		}

		if input == "n" {
			return false
		}
		if input == "Y" {
			return true
		}
	}
}

// Returns a slice filter iterator.
//
// Example:
//
//	for i, v := range filter([]int{1,2,3,4,5}, isEven) …
func filter[T any](slice []T, fn func(i int, v T) bool) iter.Seq2[int, T] {
	return func(yield func(i int, v T) bool) {
		for i, v := range slice {
			if fn(i, v) && !yield(i, v) {
				return
			}
		}
	}
}
