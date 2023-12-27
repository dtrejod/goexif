# goexif

![Build Status](https://github.com/dtrejod/goexif/actions/workflows/go.yml/badge.svg)


The following is a Go based EXIF CLI tool.

## Commands

Commands available in the CLI.

### sort

Sort media into directories based on the date and time the media was taken by
running the `sort` sub-command. The sort command will reference the media
metadata (e.g. EXIF) looking for common `Date` tags in order to determine when
the media was taken. When a date is found, the resulting media will be placed
into a new folder that uses the following `YYY/MM/DD` convention:

Example:
```
# Input
./image1.jpg
./image2.heic

# goexif sort
$ ./goexif sort --src-dir . --dry-run=false

# Output - YYYY/MM/DD
./2023/01/01/image1.jpg
./2023/02/01/image2.heic
```

Reference the help text for the `sort` [command](./cmd/sort.go) for available options.

```
./goexif sort --help
```

### date

Date prints the discovered date metadata from the media

Example:
```
# goexif date
$ ./goexif date --src-file IMG.jpg
{"level":"info","ts":1703708105.232586,"caller":"mediadate/date.go:26","msg":"Found date metadata for media.","humanTimestamp":"2020-06-26 23:19:26 +0000 UTC","unixTimestamp":1593213566}
```
