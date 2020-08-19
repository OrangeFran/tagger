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
    verbose bool                // specify how much the code should spit out
    dry_run bool                // only show what the code would do
    format string = ""          // how to extract/use information in the name
    target string = "./"        // use that file/directory
)

func main() {
    // create the cli-flags and more
    app := &cli.App{
        Name: "tagger",
        Usage: "tag mp3 files from the cmdline",
        Commands: []*cli.Command {
            {
                Name: "get",
                Aliases: []string{"g"},
                Usage: "query tags",
                Flags: []cli.Flag {
                    &cli.StringFlag {
                        Name: "target",
                        Aliases: []string{"t"},
                        Usage: "query this `TARGET`",
                        Destination: &target,
                        Required: true,
                    },
                },
                Action: func(c *cli.Context) error {
                    return get()
                },
            },
            {
                Name: "set",
                Aliases: []string{"s"},
                Usage: "tag files",
                Flags: []cli.Flag {
                    &cli.BoolFlag {
                        Name: "dry-run",
                        Usage: "only show what would be done",
                        Destination: &dry_run,
                    },
                    &cli.BoolFlag {
                        Name: "verbose",
                        Aliases: []string{"v"},
                        Usage: "specify the amount of output",
                        Destination: &verbose,
                    },
                    &cli.StringFlag {
                        Name: "target",
                        Aliases: []string{"t"},
                        Usage: "tag this `TARGET`",
                        Destination: &target,
                        Required: true,
                    },
                    &cli.StringFlag {
                        Name: "format",
                        Aliases: []string{"f"},
                        Usage: "specify `FORMAT` to tag files",
                        Destination: &format,
                        Required: true,
                    },
                },
                Action: func(c *cli.Context) error {
                    return tag()
                },
            },
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}

// used to query tags
func get() error {
    // check if target is a file
    fi, err := os.Stat(target)
    if err != nil {
        return err
    }
    // if it's a file, extract information and return
    if mode := fi.Mode(); mode.IsRegular() {
        if !strings.Contains(target, ".mp3") {
            fmt.Printf("\n   Skipping  %s", target)
        }
        // open the file as an mp3 one
        id3File, err := id3.Open(path.Join(target))
        if err != nil {
            return errors.New(fmt.Sprintf("Aborting ...\n⁉️  Failed to open %s", fi.Name()))
        }
        // extrace information
        fm := Formatter {}
        err = fm.Query(id3File)
        if err != nil {
            return err
        }
        id3File.Close()
        // print out lots of information
        fmt.Printf("\n🎵 Querying %s\n\n", fi.Name())
        for key, val := range fm.Status() {
            fmt.Printf("\t%s: %s\n", key, val)
        }

        return nil
    }

    // if we're here, we know that target is a directory
    // get all files from the directory
    files, err := ioutil.ReadDir(target)
    if err != nil {
        return err
    }

    for _, file := range files {
        // check if it's an .mp3 file
        if !strings.Contains(file.Name(), ".mp3") {
            fmt.Printf("\n   Skipping  %s\n", file.Name())
            continue
        }
        // open the file as an mp3 one
        id3File, err := id3.Open(path.Join(target, file.Name()))
        if err != nil {
            return errors.New(fmt.Sprintf("Aborting ...\n⁉️  Failed to open %s",  file.Name()))
        }
        // extrace information
        fm := Formatter {}
        err = fm.Query(id3File)
        if err != nil {
            return err
        }
        id3File.Close()
        // print out lots of information
        fmt.Printf("\n🎵 Querying %s\n\n", file.Name())
        for key, val := range fm.Status() {
            fmt.Printf("\t%s: %s\n", key, val)
        }
    }

    return nil
}

// used to tag files
// based on target and format
func tag() error {
    if dry_run { fmt.Println("\n🏃 Running in dry-run mode") }

    // check if target is a file
    fi, err := os.Stat(target)
    if err != nil {
        return err
    }
    // if it's a file, extract information and return
    if mode := fi.Mode(); mode.IsRegular() {
        if !strings.Contains(target, ".mp3") {
            fmt.Printf("\n   Skipping  %s", target)
        }
        // extrace information
        return add(fi.Name(), true)
    }

    // if we're here, we know that target is a directory
    // get all files from the directory
    files, err := ioutil.ReadDir(target)
    if err != nil {
        return err
    }

    for _, file := range files {
        // check if it's an .mp3 file
        if !strings.Contains(file.Name(), ".mp3") {
            fmt.Printf("\n   Skipping  %s\n", file.Name())
            continue
        }
        // extrace information
        err := add(file.Name(), false)
        if err != nil {
            return err
        }
    }

    return nil
}

// analyzes the name of a ".mp3" file
// and extracts information specified in the format string
func add(file string, isFile bool) error {
    //  - remove extension
    //  - split at " - "
    name := strings.ReplaceAll(file, ".mp3", "")

    fm := Formatter {}
    err := fm.Extract(name)
    if err != nil {
        return err
    }

    if verbose {
        // print out lots of information
        fmt.Printf("\n🎵 Tagging %s\n\n", file)
        for key, val := range fm.Status() {
            fmt.Printf("\t%s: %s\n", key, val)
        }
    } else {
        fmt.Printf("\n🎵 Tagging %s", file)
    }

    // acutally tag the file
    // only if dry-run is false
    if !dry_run {
        var err error
        var id3File *id3.File
        // actually open the file as an mp3 one
        if isFile {
            id3File, err = id3.Open(path.Join(target))
        } else {
            id3File, err = id3.Open(path.Join(target, file))
        }

        if err != nil {
            return errors.New(fmt.Sprintf("Aborting ...\n⁉️  Failed to open %s", file))
        }

        fm.Apply(*id3File)
        id3File.Close()
    }

    return nil
}
