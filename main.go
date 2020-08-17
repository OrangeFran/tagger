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
func add(file os.FileInfo, format, directory string, dry_run bool) error {
    //  - remove extension
    //  - split at " - "
    name := strings.ReplaceAll(file.Name(), ".mp3", "")

    fm := Formatter {}
    err := fm.Extract(name, format)
    if err != nil {
        return err
    }

    // status message
    fmt.Printf("🎵 Tagging %s\n\n", file.Name())
    for key, val := range fm.Status() {
        fmt.Printf("\t%s: %s", key, val)
    }

    // acutally tag the file
    // only if dry-run is false
    if !dry_run {
        id3File, err := id3.Open(path.Join(directory, file.Name()))
        if err != nil {
            return errors.New(fmt.Sprintf("⁉️  Failed to open %s", file.Name()))
        }

        fm.Apply(*id3File)
        id3File.Close()
    }

    return nil
}

func main() {
    // look for files in this directory
    var directory string
    // how to extract/use information in the name
    var format string
    // only show what the code would do
    var dry_run bool

    // some flag values
    flag.BoolVar(&dry_run, "dry-run", false, "only show what would happen")
    flag.StringVar(&directory, "d", "./", "specify the directory")
    flag.StringVar(&format, "f", "%a - %t", "specify the format of the file names")
    flag.Parse()

    if dry_run { fmt.Println("🏃 Running in dry-run mode") }

    // get all files from the directory
    files, err := ioutil.ReadDir(directory)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("⏳ Analyzing %s ...\n\n", directory)

    for _, file := range files {
        // check if it's an .mp3 file
        if !strings.Contains(file.Name(), ".mp3") {
            fmt.Printf("🚫 Skipping  %s\n", file.Name())
            continue
        }
        // extrace information
        err := add(file, format, directory, dry_run)
        if err != nil {
            log.Fatal(err)
        }
    }

    fmt.Println("\n✅ Finished!")
}
