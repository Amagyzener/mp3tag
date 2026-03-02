// Copyright 2026 Zmicer Pasternak. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/amagyzener/mp3tag/id3v1"
	"github.com/amagyzener/mp3tag/id3v2"
	"github.com/spf13/cobra"
)

func init() {
	var resetCmd = &cobra.Command{
		Use:     "reset PATH",
		Short:   "Reset tags (ID3v1 & ID3v2)",
		Long:    `Reset tags (ID3v1 & ID3v2)`,
		Args:    cobra.ExactArgs(1),
		Example: "mp3tag reset path/to/file.mp3",
		Run: func(cmd *cobra.Command, args []string) {
			// Check correct path & extension.
			if ext, expect := filepath.Ext(args[0]), ".mp3"; ext != expect {
				log.Fatalf(invalidFileFormatMsg, expect)
			}

			var (
				flags          = cmd.Flags()
				hasFlags       = flags.NFlag() > 0
				v1Flag, v2Flag = true, true
			)
			if hasFlags {
				v1Flag, _ = flags.GetBool("v1")
				v2Flag, _ = flags.GetBool("v2")
			}

			var pathError *os.PathError // @TODO: remove this line in Golang 1.26+

			// Reset v1.
			if v1Flag {
				var tagV1, err = id3v1.Open(args[0], id3v1.Options{Parse: false})
				// @TODO: use `errors.AsType` in Golang 1.26+
				if errors.As(err, &pathError) {
					log.Fatal(err)
				}

				if err := tagV1.SaveTo(args[0]); err != nil {
					log.Fatalf("ID3v1 save: %v", err)
				}
			}

			// Reset v2.
			if v2Flag {
				var tagV2, err = id3v2.Open(args[0], id3v2.Options{Parse: false})
				// @TODO: use `errors.AsType` in Golang 1.26+
				if errors.As(err, &pathError) {
					log.Fatal(err)
				}
				defer tagV2.Close()

				if err := tagV2.Save(); err != nil {
					log.Fatalf("ID3v2 save: %v", err)
				}
			}

			if v1Flag || v2Flag {
				log.Printf("Reset successfully to: %v", args[0])
			}
		},
	}

	var flagSet = resetCmd.Flags()
	flagSet.Bool("v1", false, "reset ID3v1 tag")
	flagSet.Bool("v2", false, "reset ID3v2 tag")

	rootCmd.AddCommand(resetCmd)
}
