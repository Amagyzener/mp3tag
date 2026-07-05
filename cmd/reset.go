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
	const (
		v1Flag = "v1"
		v2Flag = "v2"
	)

	var resetCmd = &cobra.Command{
		Use:     "reset PATH",
		Short:   "Reset ID3v1 & ID3v2 tags",
		Long:    `Reset ID3v1 & ID3v2 tags`,
		Args:    cobra.ExactArgs(1),
		Example: "mp3tag reset path/to/file.mp3",
		Run: func(cmd *cobra.Command, args []string) {
			// Check correct path & extension.
			if ext, expect := filepath.Ext(args[0]), ".mp3"; ext != expect {
				log.Fatalf(invalidFileFormatMsg, expect)
			}

			var (
				flags                = cmd.Flags()
				hasFlags             = flags.NFlag() > 0
				hasV1Flag, hasV2Flag = true, true
			)
			if hasFlags {
				hasV1Flag, _ = flags.GetBool(v1Flag)
				hasV2Flag, _ = flags.GetBool(v2Flag)
			}

			// Reset v1.
			if hasV1Flag {
				var tagV1, err = id3v1.Open(args[0], id3v1.Options{Parse: false})
				if err, ok := errors.AsType[*os.PathError](err); !ok {
					log.Fatal(err)
				}

				if err := tagV1.SaveTo(args[0]); err != nil {
					log.Fatalln("ID3v1 save error:", err)
				}

				log.Println("Reset ID3v1 successfully to:", args[0])
			}

			// Reset v2.
			if hasV2Flag {
				var tagV2, err = id3v2.Open(args[0], id3v2.Options{Parse: false})
				if err, ok := errors.AsType[*os.PathError](err); !ok {
					log.Fatal(err)
				}
				defer tagV2.Close()

				if err := tagV2.Save(); err != nil {
					log.Fatalln("ID3v2 save error:", err)
				}

				log.Println("Reset ID3v2 successfully to:", args[0])
			}
		},
	}

	var flagSet = resetCmd.Flags()
	flagSet.Bool(v1Flag, false, "reset ID3v1 tag")
	flagSet.Bool(v2Flag, false, "reset ID3v2 tag")

	rootCmd.AddCommand(resetCmd)
}
