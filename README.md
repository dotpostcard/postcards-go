# Postcard

A **work in progress** library for interacting with `.postcard` files, and the current reference implementation of the structure of this filetype.

A `.postcard` file represents a physical postcard digitally; a double-sided image with metadata like: URLs referencing to the sender(s) and receiver(s), the date it was sent on, the lat/long of where it was sent from, transcriptions of any written text on the front/back, and image descriptions of the front/back.

The contained CLI tool, `postcarder`, is able to compile front & back images, and a metadata file into a `.postcard` file:

```bash
go install github.com/jphastings/postcards-go/cmd/postcarder@latest
postcarder compile fixtures/hello-front.jpg
```

## Implementation notes

The package at the root of this repo handles the 'steady state' interaction with postcard files — eg. readinging & writing — and have extremely limited dependencies. The packages one level deeper (eg. `compile` and `validate`) hold functionality that needs complex dependencies (eg. a capable image processing library), except for `internal`, which holds packages that are common dependencies.

## File format notes

A `.postcard` file is a tarball with 4 files (in the following order):

1. A `postcard-vX.Y.Z` empty file, where `X.Y.Z` in the filename is the semantic version of the library that created it.
  - This means that the first 8 bytes of a `.postcard` file are always `postcard` (ie. `70 6f 73 74 63 61 72 64`)
2. A JSON metadata file (see `types.go` for spec)
3. A WebP image file representing the front of the postcard.
  - Physical dimensions should be correctly set (ie. DPI)
  - Ideally with a transparent background
    - Hopefully the `postcarder` tool will one day help with auto-removal of the background of scanned postcards
4. An image for the back of the postcard (identical in format to 3)
  - The physical dimensions should be within 1% of physical dimensions of the front of the postcard. This allows for different resolutions on front & back
  - Ideally co-registered, so flipping about the verical or horizontal axis (for homoriented postcards) or about one of the diagonal axes (for heteroriented postcards) have the same or extremely similar outlines

The ordering of these files in the tarball matters so version and metadata can be assessed early as larger postcard files are being streamed.
