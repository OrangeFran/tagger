package parser

import (
    "errors"
    "text/scanner"
    "strings"

    id3 "github.com/mikkyang/id3-go"
)

// if a value is empty, do not set it
func empty(a string) bool {
    if len(a) == 0 { return true } else { return false }
}

// hold data which gets added
// to the .mp3 file
type Formatter struct {
    Artist string
    Title string
    Album string
    Year string
    Genre string
}

func (f Formatter) Status() map[string]string {
    info := make(map[string]string)

    if !empty(f.Artist) {
        info["artist"] = f.Artist
    }
    if !empty(f.Title) {
        info["title"] = f.Title
    }
    if !empty(f.Album) {
        info["album"] = f.Album
    }
    if !empty(f.Year) {
        info["year"] = f.Year
    }
    if !empty(f.Genre) {
        info["genre"] = f.Genre
    }

    return info
}

func (fm Formatter) Output(format string) (string, error) {
    output := ""
    // create a scanner to loop through each character
    var f rune
    var form scanner.Scanner
    form.Init(strings.NewReader(format))

    for {
        switch f = form.Next(); f {
        case '\\':
            f = form.Next()
        case '%':
            // add the specified
            // information to the output string
            switch form.Next() {
            case 'a':
                output = output + fm.Artist
            case 't':
                output = output + fm.Title
            case 'l':
                output = output + fm.Album
            case 'y':
                output = output + fm.Year
            case 'g':
                output = output + fm.Genre
            default:
                return "", errors.New("[-] Invalid format")
            }
            continue
        case scanner.EOF:
            return output, nil
        }
        // if nothing matched, just add the char
        // to the output string
        output = output + string(f)
    }
}

func (f Formatter) Apply(file id3.File) {
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

// simply query the information
// and add it to the struct
func (fm *Formatter) Query(file *id3.File) error {
    artist := file.Artist()
    title := file.Title()
    album := file.Album()
    year := file.Year()
    genre := file.Genre()

    if !empty(artist) {
        fm.Artist = artist
    }
    if !empty(title) {
        fm.Title = title
    }
    if !empty(album) {
        fm.Album = album
    }
    if !empty(year) {
        fm.Year = year
    }
    if !empty(genre) {
        fm.Genre = genre
    }

    return nil
}

// extracts information out of content (typically the title of file)
// based on the format variable, which is using the following specifiers/identifiers
//
// %a   -> the artist
// %t   -> the title
// %l   -> the name of the album
// %y   -> the year
// %g   -> the genre
//
// one simple format I often use is: %a - %t
// this means I save my files like this:
//      "Justin Bieber - Baby.mp3"
// this is just an example, definetely not my taste
func (fm *Formatter) Extract(content, format string) error {
    // loop through each char in format
    // and match it with content
    //
    // if a % is found, look for the following string
    // and put he read information into the field it belongs to
    var c, f rune
    var form, cont scanner.Scanner
    cont.Init(strings.NewReader(content))
    form.Init(strings.NewReader(format))

    split := ""
    noscan, until_end := false, false
    var field, specifier string

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

        // allows escaping characters
        // with a backslash
        if f == '\\' {
            f = form.Next()
            if c == scanner.EOF || f == scanner.EOF {
                return errors.New("[-] Invalid format")
            }
            continue
        }

        if f == '%' {
            f = form.Next()
            if f == scanner.EOF {
                return errors.New("[-] Invalid format")
            }
            specifier = string(f)
            // flush the split and field vars
            split = ""
            field = ""
            for {
                f = form.Next()
                if f == scanner.EOF {
                    until_end = true
                    break
                }

                // allows escaping characters
                // with a backslash
                if f == '\\' {
                    f = form.Next()
                    if c == scanner.EOF || f == scanner.EOF {
                        return errors.New("[-] Invalid format")
                    }
                    continue
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
                    return errors.New("[-] Invalid format")
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
            case "a":
                fm.Artist = field
            case "t":
                fm.Title = field
            case "l":
                fm.Album = field
            case "y":
                fm.Year = field
            case "g":
                fm.Genre = field
            }

            continue
        }

        if c != f {
            return errors.New("[-] Invalid format")
        }

    }

    return nil
}
