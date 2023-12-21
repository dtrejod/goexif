<p align=right>
![Build Status](https://github.com/dtrejod/goexif/actions/workflows/go.yml/badge.svg)
</p>

# goexif

The following is a Go based EXIF CLI tool.

## Commands

Commands available in the CLI.

### sort

Sort media into directories based on the date and time the media was taken by
running the `sort` sub-command. The sort command will reference the EXIF
metadata looking for common `Date` tags in order to determine when the media
was taken. When a date is found, the resulting media will be placed into a new
folder that uses the following `YYY/MM/DD` convention:

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
