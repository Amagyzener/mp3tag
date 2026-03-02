// Copyright 2026 Zmicer Pasternak. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "mp3tag",
	Short:   "Mp3Tag is a command line tool for ID3 tag editing",
	Long:    "Mp3Tag is a command line tool for ID3 tag editing",
	Example: "mp3tag edit path/to/file.mp3",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
