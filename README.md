# Postcard

A **work in progress** CLI & library for interacting with `.postcard` files, and the current reference implementation of the structure of this filetype.

A `.postcard` file represents a physical postcard digitally; a double-sided image with metadata like: URLs referencing to the sender(s) and receiver(s), the date it was sent on, the lat/long of where it was sent from, transcriptions of any written text on the front/back, and image descriptions of the front/back.

The contained CLI tool, `postcards`, is able to compile front & back images, and a metadata file into a `.postcard` file:

```bash
go install github.com/dotpostcard/postcards-go/cmd/postcards@latest
postcards compile fixtures/hello-front.jpg
```

## Implementation notes

The package at the root of this repo handles the 'steady state' interaction with postcard files — eg. readinging & writing — and have extremely limited dependencies. The packages one level deeper (eg. `compile` and `validate`) hold functionality that needs complex dependencies (eg. a capable image processing library), except for `internal`, which holds packages that are common dependencies.

## File format notes

A `.postcard` file is a 4 part file with the following parts. The first is fixed length, the others start with a uint32 describing their length. All integers are big endian.

1. The string `postcard`, followed by three uint8 bytes representing the `X.Y.Z` of the version number of the library that created the postcard (in that order).
2. (a uint32 defining the length of this section and) JSON for the postcard metadata (see `types.go` for spec)
3. (a uint32 defining the length of this section and) a WebP image file representing the front of the postcard.
  - Ideally with a transparent background
    - Hopefully the `postcarder` tool will one day help with auto-removal of the background of scanned postcards
4. (a uint32 defining the length of this section and) a WebP image for the back of the postcard (identical in format to 3)
  - The physical dimensions should be within 1% of physical dimensions of the front of the postcard. This allows for different resolutions on front & back
  - Ideally co-registered, so flipping about the verical or horizontal axis (for homoriented postcards) or about one of the diagonal axes (for heteroriented postcards) have the same or extremely similar outlines

The ordering of these chunks is set so version and metadata can be assessed early as larger postcard files are being streamed.
