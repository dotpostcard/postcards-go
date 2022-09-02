# Postcarder

An attempt to build a CLI tool that will take a front & back scanned image of a postcard and try to produce a `.postcard` file. This will be an IPFS CARv2 file with a single root, holding the following IPLD references:

- Images for the front and back of the postcard
  - Physical dimensions should be correctly set (ie. DPI)
  - Ideally co-registered, so flipping about the verical axis (for homoriented postcards) or about one of the diagonal axis (for heteroriented postcards) have the same or extremely similar outlines
  - Ideally with a transparent background, either with an additional mask image, or using a transparency-capanble image format (webp?)
- A metadata file containing information about the postcard. Possibly including:
  - Date sent
  - Date received
  - Sender URI(s)
  - Intended recipient URI(s)
  - Transcription of back (or front!) text
  - Alt-text of front/back images

## Progress

The complex parts of this code are:
- Automatic co-registration
- Automatic removal of background
  - This is where I've been focusing my efforts
