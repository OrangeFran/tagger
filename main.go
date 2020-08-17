package main

import (
    "fmt"
    "log"
    "os"
    "path"
    "errors"
    "strings"
    "io/ioutil"

    id3 "github.com/mikkyang/id3-go"
    cli "github.com/urfave/cli"
)

var (
    dry_run bool                // only show what the code would do
    format string = ""          // how to extract/use information in the name
    directory string = "./"     // look for files in this directory
)


// analyzes the name of a ".mp3" file
// and extracts information specified in the format string
func add(file os.FileInfo) error {
    //  - remove extension
    //  - split at " - "
    name := strings.ReplaceAll(file.Name(), ".mp3", "")

    fm := Formatter {}
    err := fm.Extract(name, format)
    if err != nil {
        return err
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
    } else {
        // just print out some information
        fmt.Printf("\n🎵 Tagging %s\n\n", file.Name())
        for key, val := range fm.Status() {
            fmt.Printf("\t%s: %s\n", key, val)
        }
    }

    return nil
}

func main() {
    app := &cli.App{
        Name: "tagger",
        Usage: "tag mp3 files from the cmdline",
        Flags: []cli.Flag{
            &cli.BoolFlag {
                Name: "dry-run",
                Usage: "only show what would be done",
                Destination: &dry_run,
            },
            &cli.StringFlag{
                Name: "directory",
                Aliases: []string{"d"},
                Usage: "tag all files from `DIRECTORY`",
                Destination: &directory,
            },
            &cli.StringFlag{
                Name: "format",
                Aliases: []string{"f"},
                Usage: "extract `FORMAT` out of the file names",
                Destination: &format,
            },
        },
        Action: func(c *cli.Context) error {
            if format == "" {
                return errors.New("🚫 No format specified")
            }
            return run()
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}

func run() error {
    if dry_run { fmt.Println("🏃 Running in dry-run mode") }

    // get all files from the directory
    files, err := ioutil.ReadDir(directory)
    if err != nil {
        return err
    }

    fmt.Printf("⏳ Analyzing %s ...\n", directory)

    for _, file := range files {
        // check if it's an .mp3 file
        if !strings.Contains(file.Name(), ".mp3") {
            fmt.Printf("🚫 Skipping  %s\n", file.Name())
            continue
        }
        // extrace information
        err := add(file)
        if err != nil {
            return err
        }
    }

    fmt.Println("\n✅ Finished!")
    return nil
}
