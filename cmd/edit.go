// Copyright 2026 Zmicer Pasternak. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"iter"
	"log"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/amagyzener/mp3tag/id3v2"
	"github.com/spf13/cobra"
)

const editFrameMsg = "%v: type any text to change, \"_\" to delete, or empty string to leave unchanged\n\tCurrent: %q\n\tNew: "

func init() {
	var frames = []frameRecord{
		{"TIT2", "Title"},
		{"TPE1", "Artist"},
		{"TPE2", "Band/Orchestra/Accompaniment"},
		{"TALB", "Album/Movie/Show"},
		{"TYER", "Year"},
		{"TRCK", "Track number/Position in set"},
		{"TPOS", "Part of a set"},
		{"TCON", "Content type"},
		{"TBPM", "BPM"},
		{"TCOP", "Copyright message"},
		{"TPUB", "Publisher"},
		{"TCOM", "Composer"},
		{"TEXT", "Lyricist/Text writer"},
		{"APIC", "Attached picture"},
		{"USLT", "Unsynchronised lyrics/text transcription"},
		{"COMM", "Comments"},
	}

	var editCmd = &cobra.Command{
		Use:     "edit PATH",
		Short:   "Edit ID3v2.3 tag",
		Long:    `Edit ID3v2.3 tag`,
		Args:    cobra.ExactArgs(1),
		Example: "mp3tag edit --tit2 --talb path/to/file.mp3",
		Run: func(cmd *cobra.Command, args []string) {
			// Check correct parts & extension.
			if ext, expect := filepath.Ext(args[0]), ".mp3"; ext != expect {
				log.Fatalf(invalidFileFormatMsg, expect)
			}

			tag, err := id3v2.Open(args[0], id3v2.Options{Parse: true})
			if err != nil {
				log.Fatalf("ID3v2: %v", err)
			}
			defer tag.Close()

			tag.SetVersion(3)
			tag.SetDefaultEncoding(id3v2.EncodingUTF16)

			flags := cmd.Flags()
			isAll, err := flags.GetBool("all")
			if err != nil {
				log.Fatal(err)
			}

			var editFrames iter.Seq2[int, frameRecord]

			if isAll {
				// Edit all frames except the specified ones.
				editFrames = filter(frames, func(_ int, r frameRecord) bool {
					flagState, _ := flags.GetBool(strings.ToLower(r.frame))
					return !flagState
				})
			} else {
				// Edit the specified frames OR edit the first 8 otherwise.
				if hasFlags := flags.NFlag() > 0; hasFlags {
					editFrames = filter(frames, func(_ int, r frameRecord) bool {
						flagState, _ := flags.GetBool(strings.ToLower(r.frame))
						return flagState
					})
				} else {
					editFrames = slices.All(frames[0:8]) // TIT2..=TCON
				}
			}

			scanner := bufio.NewScanner(os.Stdin)
			for _, v := range editFrames {
				switch v.frame {
				case "APIC":
					if yes := confirmEdit(scanner, v.String()); !yes {
						continue
					}

					for {
						fmt.Print("\tType path/to/file.{png|jpg|jpeg}: ")
						imagePath, err := scanInputLine(scanner)
						if err != nil {
							log.Println(err)
							continue
						}

						mimeTypes := map[string]string{
							// NOTE: there is no MIME type "image/jpg"!
							".png":  "image/png",
							".jpg":  "image/jpeg",
							".jpeg": "image/jpeg",
						}
						ext := filepath.Ext(imagePath)
						mimeType := mimeTypes[ext]

						if mimeType == "" {
							log.Printf(invalidFileFormatMsg, slices.Collect(maps.Keys(mimeTypes)))
							continue
						}

						picBytes, err := os.ReadFile(imagePath)
						if err != nil {
							log.Println(err)
							continue
						}

						tag.DeleteFrames(v.frame) // remove the old pic
						tag.AddAttachedPicture(id3v2.PictureFrame{
							Encoding:    id3v2.EncodingISO,
							MimeType:    mimeType,
							PictureType: id3v2.PTFrontCover,
							Description: "Front cover",
							Picture:     picBytes,
						})

						break
					}
				case "USLT":
					if yes := confirmEdit(scanner, v.String()); !yes {
						continue
					}

					var iso6392code string
					for {
						fmt.Print(
							"\tType lyrics language according to ISO-639-2, e. g. \"eng\", \"hun\", \"rus\", \"ukr\"\n" +
								"\tSee https://en.wikipedia.org/wiki/List_of_ISO_639-2_codes\n" +
								"\tCode: ",
						)
						if iso6392code, err = scanInputLine(scanner); err != nil {
							log.Println(err)
							continue
						}
						if len(iso6392code) == 3 {
							break
						}
					}

					for {
						fmt.Print("\tType path/to/file.txt (in UTF-8): ")
						txtPath, err := scanInputLine(scanner)
						if err != nil {
							log.Println(err)
							continue
						}

						if ext, expect := filepath.Ext(txtPath), ".txt"; ext != expect {
							log.Printf(invalidFileFormatMsg, expect)
							continue
						}

						lyrics, err := os.ReadFile(txtPath) // UTF-8 txt only
						if err != nil {
							log.Println(err)
							continue
						}

						tag.DeleteFrames(v.frame) // remove the old lyrics
						tag.AddUnsynchronisedLyricsFrame(id3v2.UnsynchronisedLyricsFrame{
							Encoding: tag.DefaultEncoding(),
							Language: iso6392code,
							Lyrics:   string(bytes.TrimSpace(lyrics)),
						})

						break
					}
				case "COMM":
					var frameText string
					if frames := tag.GetFrames(v.frame); len(frames) > 0 {
						idx := slices.IndexFunc(frames, func(e id3v2.Framer) bool {
							return len(e.(id3v2.CommentFrame).Description) == 0
						})
						if idx != -1 {
							frameText = frames[idx].(id3v2.CommentFrame).Text
						}
					}

					fmt.Printf(editFrameMsg, v.String(), frameText)
					if commentary, err := scanInputLine(scanner); err == nil {
						tag.DeleteFrames(v.frame) // remove the old comment
						tag.AddCommentFrame(id3v2.CommentFrame{
							Encoding:    tag.DefaultEncoding(),
							Description: "",    // always empty for commentary
							Language:    "eng", // hardcoded
							Text:        commentary,
						})
					}
				default:
					fmt.Printf(editFrameMsg, v.String(), tag.GetTextFrame(v.frame).Text)
					if input, err := scanInputLine(scanner); err == nil {
						if input == "_" {
							tag.DeleteFrames(v.frame)
						} else {
							tag.AddTextFrame(v.frame, tag.DefaultEncoding(), input)
						}
					}
				}
			}

			fmt.Println(strings.Repeat("-", 50))
			fmt.Println(tag)

			if err := tag.Save(); err != nil {
				log.Fatal(err)
			}
			log.Printf("Saved successfully to: %q", args[0])
		},
	}

	flagSet := editCmd.Flags()
	flagSet.BoolP("all", "a", false, "edit all frames except the ones specified")
	for _, opts := range frames {
		flagSet.Bool(
			strings.ToLower(opts.frame),
			false,
			fmt.Sprintf("edit %v frame", opts.String()),
		)
	}

	rootCmd.AddCommand(editCmd)
}
