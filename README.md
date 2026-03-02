# ID3 tag editor CLI

## Commands

### Show MP3 tags

```bash
mp3tag show path/to/file.mp3
```

### Edit MP3 tags (ID3v2.3 only)

Edit all supported frames listed below.

```bash
mp3tag edit -a path/to/file.mp3
```

Edit all frames except those specified.

```bash
mp3tag edit -a --tbpm --tcop --comm path/to/file.mp3
```

Edit the first 8 frames (from `TIT2` to `TCON` inclusive).

```bash
mp3tag edit path/to/file.mp3
```

Edit the specified frames.

```bash
mp3tag edit --tit2 --tpe1 path/to/file.mp3
```

#### Flags

- `--tit2` to edit `Title`.
- `--tpe1` to edit `Artist`.
- `--tpe2` to edit `Band/Orchestra/Accompaniment`.
- `--talb` to edit `Album/Movie/Show`.
- `--tyer` to edit `Year` (_dddd_).
- `--trck` to edit `Track number/Position in set` (_dd/dd_).
- `--tpos` to edit `Part of a set` (_dd/dd_).
- `--tcon` to edit `Content type` (a genre).
- `--tbpm` to edit `BPM`.
- `--tcop` to edit `Copyright message`.
- `--tpub` to edit `Publisher`.
- `--tcom` to edit `Composer`.
- `--text` to edit `Lyricist/Text writer`.
- `--apic` to edit `Attached picture` (front cover).
- `--uslt` to edit `Unsynchronised lyrics/text transcription`.
- `--comm` to edit `Comments`.

### Remove tags

Reset all tags, both ID3v1 & ID3v2.

```bash
mp3tag reset path/to/file.mp3
```

#### Flags

- `--v1` to reset ID3v1 tag only.
- `--v2` to reset ID3v2 tag only.

## TODO

- `[~]` Fix UTF-16LE with BOM to UTF-8 decoding. (\*)
- `[v]` Fix invalid "COMM" on ID3v1 tag reset.
- `[~]` Fix incorrect re-edit (broken lyrics; due to UTF-16LE with BOM to UTF-8 decoding?). (\*)
- `[?]` Add frame flags to `show` command (show the specified frames).

\* _Turned off LE/BE detection entirely as a workaround (the program assumes UTF-16 to be LE with BOM only)._
