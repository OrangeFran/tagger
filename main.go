package main

import (
    "fmt"
    "log"
    "os"
    "path"
    "flag"
    "errors"
    "strings"
    "io/ioutil"
    id3 "github.com/mikkyang/id3-go"
)


// analyzes the name of a ".mp3" file
// and extracts information specified in the format string
func add(file os.FileInfo, format, directory string) error {
    //  - remove extension
    //  - split at " - "
    name := strings.ReplaceAll(file.Name(), ".mp3", "")

    fm := Formatter {}
    err := fm.Extract(name, format)
    if err != nil {
        return err
    }
    // var artist, title string
    // if len(pieces) == 1 {
    //     artist = "Unknown Artist"
    //     title = pieces[0]
    // } else {
    //     artist = pieces[0]
    //     title = pieces[1]
    // }

    // status message
    // fmt.Printf("🎵 Tagging %s\n\t-> %s (artist)\n\t-> %s (title)\n", file.Name(), artist, title)

    // acutally tag the file
    id3File, err := id3.Open(path.Join(directory, file.Name()))
    if err != nil {
        return errors.New(fmt.Sprintf("⁉️  Failed to open %s", file.Name()))
    }
    defer id3File.Close()

    return nil
}

func main() {
    // look for files in this directory
    var directory string
    // how to extract/use information in the name
    var format string

    // some flag values
    flag.StringVar(&directory, "d", "./", "specify the directory")
    flag.StringVar(&format, "f", "%a - %t", "specify the format of the file names")
    flag.Parse()

    // get all files from the directory
    files, err := ioutil.ReadDir(directory)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("⏳ Analyzing %s ...\n\n", directory)

    for _, file := range files {
        // check if it's an .mp3 file
        if !strings.Contains(file.Name(), ".mp3") {
            fmt.Printf("➡️  Skipping  %s\n", file.Name())
            continue
        }
        // extrace information
        err := add(file, format, directory)
        if err != nil {
            log.Fatal(err)
        }
    }

    fmt.Println("\n✅ Finished!")
}
