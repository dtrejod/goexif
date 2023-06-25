# goexif

https://github.com/dtrejod/goexif/actions/workflows/go.yml/badge.svg?branch=main


The following is a Go based EXIF CLI tool.

## Commands

Commands available in the CLI.

### sort

**WARNING: Command still in development. Use at your own risk**

Sort media into directories based on the date and time the media was taken by
running the `sort` sub-command. The sort command will reference the EXIF
metadata looking for common `Date` tags in order to determine when the media
was taken. When a date is found, the resulting media will be placed into a new folder that uses the following convention:

```
# Sort folder convention
YYYY/MM/DD/<image>
```

Reference the help text for the `sort` [command](./cmd/sort.go) for available options.
