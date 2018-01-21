package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "os/user"
    "flag"
    "errors"
)

func downloadSegment(url string, fileName string) (bytes int64, err error) {

    outFile, err := os.Create(fileName)
    if err != nil {
        return 0, err
    }
    defer outFile.Close()

    response, err := http.Get(url)
    if err != nil {
        return 0, err
    }
    defer response.Body.Close()

    if response.StatusCode != 200 {
        err := errors.New(fmt.Sprintf("Expected HTTP 200, got HTTP %d.\n", response.StatusCode))
        return 0, err
    }

    bytesWritten, err := io.Copy(outFile, response.Body)
    if err != nil {
        return 0, err
    }

    return bytesWritten, nil
}

func downloadSegments(workDir string, numSegments int, downloadUrl string, streamSuffix string) (totalBytes int64, err error){
    err = os.Mkdir(workDir, 0775)
    if err != nil {
        return 0, err
    }

    segmentlist, err := os.Create(fmt.Sprintf("%s/segmentlist", workDir))
    if err != nil {
        return 0, err
    }
    defer segmentlist.Close()

    totalBytes = 0

    for segment := 1; segment <= numSegments; segment++ {
        url := fmt.Sprintf("%s/segment%d%s", downloadUrl, segment, streamSuffix)
        fileName := fmt.Sprintf("%s/segment%d%s", workDir, segment, streamSuffix)

        fmt.Printf("Downloading url: %s...", url)
        bytesWritten, err := downloadSegment(url, fileName)

        if err != nil  && segment < 2 {
            fmt.Println("Segment download failed, breaking", "-", err)
            return totalBytes, err
        }
        fmt.Fprintf(segmentlist, "file '%s'\n", fileName)
        fmt.Println("success.", bytesWritten, "bytes downloaded.")
        totalBytes += bytesWritten
    }

    return totalBytes, nil
}

func cleanUp(workDir string) {
    fmt.Println("Cleaning up.")
    err := os.RemoveAll(workDir)
    if err != nil {
        fmt.Println("Error cleaning up workdir.", "-", err)
        return
    }
}

func main() {

    usr, err := user.Current()
    if err != nil {
        fmt.Println(err)
        return
    }

    baseUrl := flag.String("baseurl", "http://abcradiomodhls.abc-cdn.net.au/i/triplej/audio", "ABC Radio CDN URL.")
    show := flag.String("show", "hip", "Defaults to Triple J HipHop Show.")
    showNum := flag.Int("shownum", 1, "Normally just one show per week. Set to 0 to ignore.")
    showDate := flag.String("showdate", "REQUIRED", "[REQUIRED] Date in format YYYY-MM-DD.")
    showFormat := flag.String("showformat", "m4a", "Format of stream stored on CDN.")
    streamSuffix := flag.String("streamsuffix", "_0_a.ts", "Suffix of stream segments stored on CDN.")
    downloadDir := flag.String("downloaddir", fmt.Sprintf("%s/Downloads", usr.HomeDir), "Directory to download mp3 to.")
    numSegments := flag.Int("numsegments", 1080, "Number of 10s segments in a show.")
    ffmpegPath := flag.String("ffmpegpath", "/usr/bin/ffmpeg", "ffmpeg binary path.")
    flag.Parse()

    if *showDate == "REQUIRED" {
        fmt.Println("Must specify -showdate flag. Exiting.")
        return
    }

    var fileName string
    if *showNum == 0 {
        fileName = fmt.Sprintf("%s-%s", *show, *showDate)
    } else {
        fileName = fmt.Sprintf("%s-%d-%s", *show, *showNum, *showDate)
    }

    workDir := fmt.Sprintf("%s/%s", *downloadDir, fileName)

    totalBytes, err := downloadSegments(workDir, *numSegments, fmt.Sprintf("%s/%s.%s", *baseUrl, fileName, *showFormat), *streamSuffix)
    if err != nil {
        fmt.Println("Error downloading segments", "-", err)
        cleanUp(workDir)
        return
    }

    outputTSPath := fmt.Sprintf("%s/output.ts", workDir)

    concatCmd := exec.Command(*ffmpegPath, "-safe", "0", "-f", "concat", "-i", fmt.Sprintf("%s/segmentlist", workDir), "-acodec", "copy", outputTSPath)
    concatOutput, err := concatCmd.CombinedOutput()
    fmt.Println(concatCmd.Args)
    if err != nil {
        os.Stderr.WriteString(fmt.Sprintf("Concat error: %s\n",err.Error()))
        cleanUp(workDir)
        return
    }
    if len(concatOutput) > 0 {
        fmt.Println(string(concatOutput))
    }

    encodeCmd := exec.Command(*ffmpegPath, "-i", outputTSPath, "-safe", "0", "-f", "mp3", "-acodec", "mp3", fmt.Sprintf("%s/%s.mp3", *downloadDir, fileName))
    encodeOutput, err := encodeCmd.CombinedOutput()
    if err != nil {
        os.Stderr.WriteString(fmt.Sprintf("Encode error: %s\n",err.Error()))
        cleanUp(workDir)
        return
    }
    if len(encodeOutput) > 0 {
        fmt.Println(string(encodeOutput))
    }

    cleanUp(workDir)
    fmt.Print("Exiting,", totalBytes, "downloaded.")
}
