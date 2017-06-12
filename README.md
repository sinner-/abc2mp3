# abc2mp3
Convert ABC radio streams to mp3

## Requirements
 * Tested on Fedora 24.
 * ffmpeg: `dnf install ffmpeg`.
 * Golang: `dnf install golang`.

# Building

 * Make sure `$GOPATH` is set and `$GOPATH/src` exists!
 * `git clone git@github.com:sinner-/abc2mp3.git $GOPATH/src/github.com/sinner-/abc2mp3`
 * `go install github.com/sinner-/abc2mp3`

## Usage

 * Make sure that `$PATH` includes `$GOBIN`!
 * I mostly just use this for the Triple J Hip Hop Show, e.g. `abc2mp3 -showdate 2017-01-12`.
 * Find your encoded mp3 in `$HOME/Downloads/hip-1-2017-01-12.mp3` on successful execution.

```
Usage of abc2mp3:
  -baseurl string
        ABC Radio CDN URL. (default "http://abcradiomodhls.abc-cdn.net.au/i/triplej/audio")
  -downloaddir string
        Directory to download mp3 to. (default "/home/sina/Downloads")
  -ffmpegpath string
        ffmpeg binary path. (default "/usr/bin/ffmpeg")
  -numsegments int
        Number of 10s segments in a show. (default 1081)
  -show string
        Defaults to Triple J HipHop Show. (default "hip")
  -showdate string
        [REQUIRED] Date in format YYYY-MM-DD. (default "REQUIRED")
  -showformat string
        Format of stream stored on CDN. (default "m4a")
  -shownum int
        Normally just one show per week. (default 1)
  -streamsuffix string
        Suffix of stream segments stored on CDN. (default "_0_a.ts")
```

## Download DoubleJ Something Different
  * `abc2mp3 -baseurl http://abcradiomodhls.abc-cdn.net.au/i/doublej/audio -show som -shownum 0 -numsegments 720 -showdate 2017-06-06`
