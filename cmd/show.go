// Copyright 2026 Zmicer Pasternak. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/amagyzener/mp3tag/id3v1"
	"github.com/amagyzener/mp3tag/id3v2"
	"github.com/spf13/cobra"
)

func init() {
	var showCmd = &cobra.Command{
		Use:     "show PATH",
		Short:   "Show all tags (ID3v1 & ID3v2)",
		Long:    `Show all tags (ID3v1 & ID3v2)`,
		Args:    cobra.ExactArgs(1),
		Example: "mp3tag show path/to/file.mp3",
		Run: func(cmd *cobra.Command, args []string) {
			// Check correct path & extension.
			if ext, expect := filepath.Ext(args[0]), ".mp3"; ext != expect {
				log.Fatalf(invalidFileFormatMsg, expect)
			}

			// ID3v1 tag.
			tagV1, err := id3v1.Open(args[0], id3v1.Options{Parse: true})
			if err != nil {
				log.Fatalf("ID3v1: %v", err)
			}

			fmt.Println("[ID3v1]")
			fmt.Println(tagV1)

			// ID3v2.3 tag.
			tagV2, err := id3v2.Open(args[0], id3v2.Options{Parse: true})
			if err != nil {
				log.Fatalf("ID3v2: %v", err)
			}
			defer tagV2.Close()

			fmt.Printf("[ID3v2] VERSION %v\n", tagV2.Version())
			fmt.Println(tagV2)
		},
	}

	rootCmd.AddCommand(showCmd)
}
