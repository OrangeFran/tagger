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
func Execute(target string, function func(string)error) error {
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
            fmt.Printf("[*] Skipping '%s'\n", name)
            return nil
        }

        // open the file as an mp3 one
        id3File, err := id3.Open(file)
        if err != nil {
            return errors.New(fmt.Sprintf("[-] Aborting ...\n[-] Failed to open '%s'", name))
        }
        defer id3File.Close()

        // set the tag to an empty
        // string if it should be removed
        fmt.Printf("[+] Removing from '%s'\n", name)
        if artist {
            id3File.SetArtist("")
            if verbose { fmt.Println("    |- removed artist") }
        }
        if title {
            id3File.SetTitle("")
            if verbose { fmt.Println("    |- removed title") }
        }
        if album {
            id3File.SetAlbum("")
            if verbose { fmt.Println("    |- removed album") }
        }
        if year {
            id3File.SetYear("")
            if verbose { fmt.Println("    |- removed year") }
        }
        if genre {
            id3File.SetGenre("")
            if verbose { fmt.Println("    |- removed genre") }
        }
        fmt.Println()

        return nil
    }

    return Execute(target, function)
}

// used to query tags
func Query(target string) error {
    function := func(file string) error {
        // the actual name without the path
        name := path.Base(file)

        if !strings.Contains(name, ".mp3") {
            fmt.Printf("[*] Skipping '%s'\n", name)
            return nil
        }

        // open the file as an mp3 one
        id3File, err := id3.Open(file)
        if err != nil {
            return errors.New(fmt.Sprintf("[-] Aborting ...\n[-] Failed to open '%s'", name))
        }
        defer id3File.Close()
        // extrace information
        fm := parser.Formatter {}
        err = fm.Query(id3File)
        if err != nil {
            return err
        }
        // print out lots of information
        status := fm.Status()
        if len(status) == 0 {
            fmt.Printf("[+] Nothing set for '%s'\n", name)
            return nil
        }

        fmt.Printf("[+] Querying '%s'\n", name)
        for key, val := range fm.Status() {
            fmt.Printf("    |- %s: '%s'\n", key, val)
        }
        fmt.Println()

        return nil
    }

    return Execute(target, function)
}

// used to tag files
// based on target and format
func Tag(target, format string, verbose, dry_run bool) error {
    if dry_run { fmt.Println("[*] Running in dry-run mode\n") }

    function := func(file string) error {
        // the actual name without the path
        name := path.Base(file)

        if !strings.Contains(name, ".mp3") {
            fmt.Printf("[*] Skipping '%s'\n", name)
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
            // print out lots of information
            fmt.Printf("[+] Tagging '%s'\n", name)
            for key, val := range fm.Status() {
                fmt.Printf("    |- %s: '%s'\n", key, val)
            }
            fmt.Println()
        } else {
            fmt.Printf("[+] Tagging '%s'\n", name)
        }

        // acutally tag the file
        // only if dry-run is false
        if !dry_run {
            // actually open the file as an mp3 one
            id3File, err := id3.Open(file)
            if err != nil {
                return errors.New(fmt.Sprintf("[-] Aborting ...\n[-] Failed to open '%s'", name))
            }

            fm.Apply(*id3File)
            id3File.Close()
        }

        return nil
    }

    return Execute(target, function)
}
