package main

import (
    "fmt"
    "errors"
    "text/scanner"
    "strings"
    id3 "github.com/mikkyang/id3-go"
)

// hold data which gets added
// to the .mp3 file
type Formatter struct {
    Artist string
    Title string
    Album string
    Year string
    Genre string
}

func (f Formatter) Apply(file id3.File) {
    // if a value is empty, do not set it
    empty := func(a string) bool {
        if len(a) == 0 { return true } else { return false }
    }

    if !empty(f.Artist) {
        file.SetArtist(f.Artist)
    }
    if !empty(f.Title) {
        file.SetTitle(f.Title)
    }
    if !empty(f.Album) {
        file.SetAlbum(f.Album)
    }
    if !empty(f.Year) {
        file.SetYear(f.Year)
    }
    if !empty(f.Genre) {
        file.SetGenre(f.Genre)
    }
}

func (fm *Formatter) Extract(content, format string) error {
    fmt.Printf("Extracting from %s with %s\n", content, format)
    // loop through each char in format
    // and match it with content
    //
    // if a % is found, look for the following string
    // and put he read information into the field it belongs to
    var c rune
    var cont scanner.Scanner
    cont.Init(strings.NewReader(content))

    var f rune
    var form scanner.Scanner
    form.Init(strings.NewReader(format))

    split := ""
    noscan := false
    until_end := false

    var field string
    var specifier rune
    for {
        c = cont.Next()
        if noscan {
            noscan = false
        } else {
            f = form.Next()
        }

        if c == scanner.EOF || f == scanner.EOF {
            break
        }

        if f == '%' {
            specifier = form.Next()
            // flush the split and field vars
            split = ""
            field = ""
            for {
                f = form.Next()
                if f == scanner.EOF {
                    until_end = true
                    break
                }

                if f == '%' {
                    noscan = true
                    break
                }
                split = split + string(f)
            }

            field = field + string(c)
            for {
                c = cont.Next()
                if until_end {
                    if c == scanner.EOF {
                        field = strings.ReplaceAll(field, split, "")
                        break
                    }
                    field = field + string(c)
                    continue
                }
                // if scanner.EOF was found
                if c == scanner.EOF {
                    return errors.New("Invalid format specified!")
                }
                field = field + string(c)
                if strings.Contains(field, split) {
                    // remove the string that was specified in the format
                    field = strings.ReplaceAll(field, split, "")
                    break
                }
            }
            // finally add the string to the formatter
            // and go on to the next one
            switch specifier {
            case 'a':
                fm.Artist = field
            case 't':
                fm.Title = field
            case 'l':
                fm.Album = field
            case 'y':
                fm.Year = field
            case 'g':
                fm.Genre = field
            }

            fmt.Printf("New value for %s: %s\n", string(specifier), field)

            continue
        }

        if c != f {
            return errors.New("⁉️  Invalid format specifier")
        }

    }

    return nil
}
