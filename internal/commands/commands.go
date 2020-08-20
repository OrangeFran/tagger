package commands

import (
    "os"
    "fmt"
    "path"
    "errors"
    "strings"
    "io/ioutil"

    id3 "github.com/mikkyang/id3-go"

    // local imports
    parser "github.com/orangefran/tagger/internal/parser"
)

// takes a function as input
// after looking at the target, it figures out the type
//
// if it's a file, it executes the given function on this file
// if it's a directory, it executes it for every file in that directory
//
// the string that this function passes to the "function"
// is always the absolute path to the file
func ExecuteFunc(target string, function func(string)error) error {
    // check if target is a file
    fi, err := os.Stat(target)
    if err != nil {
        return err
    }
    // if it's a file, run the function and return
    if mode := fi.Mode(); mode.IsRegular() {
        return function(target)
    }

    // if we're here, we know that target is a directory
    // get all files from the directory
    files, err := ioutil.ReadDir(target)
    if err != nil {
        return err
    }

    // and run the function on them
    for _, file := range files {
        if err := function(path.Join(target, file.Name())); err != nil {
            return err
        }
    }

    return nil
}

// remove tags from files
func Remove(target string, verbose, artist, title, album, year, genre bool) error {
    function := func(file string) error {
        // the actual name without the path
        name := path.Base(file)

        // check if it's an .mp3 file
        if !strings.Contains(name, ".mp3") {
            return nil
        }

        // open the file as an mp3 one
        id3File, err := id3.Open(file)
        if err != nil {
            return errors.New(fmt.Sprintf("Failed to open '%s'", name))
        }
        defer id3File.Close()

        // set the tag to an empty
        // string if it should be removed
        if verbose { fmt.Printf("Removing from '%s'\n", name) }
        if artist {
            id3File.SetArtist("")
        }
        if title {
            id3File.SetTitle("")
        }
        if album {
            id3File.SetAlbum("")
        }
        if year {
            id3File.SetYear("")
        }
        if genre {
            id3File.SetGenre("")
        }

        return nil
    }

    return ExecuteFunc(target, function)
}

// used to query tags
func Query(target, format string, verbose bool) error {
    function := func(file string) error {
        // the actual name without the path
        name := path.Base(file)

        if !strings.Contains(name, ".mp3") {
            return nil
        }

        // open the file as an mp3 one
        id3File, err := id3.Open(file)
        if err != nil {
            return errors.New(fmt.Sprintf("Failed to open '%s'", name))
        }
        defer id3File.Close()
        // extrace information
        fm := parser.Formatter {}
        err = fm.Query(id3File)
        if err != nil {
            return err
        }
        // print out lots of information
        output, err := fm.Output(format)
        if err != nil {
            return err
        }

        fmt.Println(output)
        // if verbose {
        //     fmt.Printf("[+] Querying '%s'\n", name)
        //     for key, val := range fm.Status() {
        //         fmt.Printf("    |- %s: '%s'\n", key, val)
        //     }
        //     fmt.Println()
        // }

        return nil
    }

    return ExecuteFunc(target, function)
}

// used to tag files
// based on target and format
func Tag(target, format string, verbose, dry_run bool) error {
    function := func(file string) error {
        // the actual name without the path
        name := path.Base(file)

        if !strings.Contains(name, ".mp3") {
            return nil
        }

        // remove extension
        content := strings.ReplaceAll(name, ".mp3", "")

        fm := parser.Formatter {}
        err := fm.Extract(content, format)
        if err != nil {
            return err
        }

        if verbose {
            fmt.Printf("Tagging '%s'\n", name)
        }

        // acutally tag the file
        // only if dry-run is false
        if !dry_run {
            // open the file as an mp3 one
            id3File, err := id3.Open(file)
            if err != nil {
                return errors.New(fmt.Sprintf("Failed to open '%s'", name))
            }
            defer id3File.Close()

            fm.Apply(*id3File)
        }

        return nil
    }

    return ExecuteFunc(target, function)
}

// sets values manually
func Static(target string, verbose bool, fm parser.Formatter) error {
    function := func(file string) error {
        // the actual name without the path
        name := path.Base(file)
        if !strings.Contains(name, ".mp3") {
            return nil
        }

        if verbose {
            // print out lots of information
            fmt.Printf("Tagging '%s'\n", name)
        }

        // open the file as an mp3 one
        id3File, err := id3.Open(file)
        if err != nil {
            return errors.New(fmt.Sprintf("Failed to open '%s'", name))
        }
        defer id3File.Close()

        fm.Apply(*id3File)

        return nil
    }

    return ExecuteFunc(target, function)
}
