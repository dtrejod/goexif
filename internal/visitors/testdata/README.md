# Testdata

## Image generator

Used linked 1x1 Pixel generator for generating testdata:
https://shoonia.github.io/1x1/#5829a5ff

## Mona Lisa

Variations of Mona Lisa testdata pulled from Wikipedia. All images are in Public Domain:

https://en.wikipedia.org/wiki/Mona_Lisa

## Using exiftool

exiftool is useful for setting EXIF metadata on exisiting media files

### Set EXIF data metadata
Set all date fields to a specific time.

```
exiftool "-AllDates=2000:01:01 00:00:00" <file>
```

where

```
-AllDates=YYYY:MM:DD HH:MM:SS
```

### Remove all EXIF metadata

```
exiftool -all= <file>
```

## Using touch

touch is useful for setting file metadata on exisiting files.


```
touch -a -m -t 200101010000 testdata/white.png
```

where

```
-a = accessed
-m = modified
-t = timestamp - use [[CC]YY]MMDDhhmm[.ss] time format
```
