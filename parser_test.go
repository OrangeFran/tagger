package main

import (
    "fmt"
    "errors"
    "testing"
)

// a helper function for evaluating test results
func EvaluateResults(results map[string]string) error {
    passed := true
    failed := ""

    for correct, got := range results {
        if got != correct {
            passed = false
            failed = failed + fmt.Sprintf("[-] Wanted %s, got %s\n", correct, got)
        }
    }

    if !passed {
        return errors.New("[-] Test failed!\n" + failed)
    }

    return nil
}

func TestParseSimple(t *testing.T) {
    // create a simple formatter
    format = "%a - %t"
    content := "artist - title"

    fm := Formatter {}
    fm.Extract(content)

    // evaluate the results
    results := map[string]string{
        "artist": fm.Artist,
        "title": fm.Title,
    }

    if err := EvaluateResults(results); err != nil {
        t.Errorf(err.Error())
    }
    // if fm.Artist != "artist" || fm.Title != "title" {
    //     t.Errorf("Simple test failed!\nWanted %s & %s, got %s & %s\n", "artist", "title", fm.Artist, fm.Title)
    // }
}

func TestParseBeginning(t *testing.T) {
    // create a simple formatter
    format = "music_%a - %t (%y)"
    content := "music_artist - title (18/10/2020)"

    fm := Formatter {}
    fm.Extract(content)

    // evaluate the results
    results := map[string]string{
        "artist": fm.Artist,
        "title": fm.Title,
        "18/10/2020": fm.Year,
    }

    if err := EvaluateResults(results); err != nil {
        t.Errorf(err.Error())
    }
}

func TestParseError(t *testing.T) {
    // create a simple formatter
    format = "music_%a - %t & (%y)"
    content := "music_artist - title (18/10/2020)"

    fm := Formatter {}
    err := fm.Extract(content)

    // evaluate the results
    if err == nil {
        t.Errorf("The function should have returned an error!")
    }
}

func TestParseErrorCutOff(t *testing.T) {
    // create a simple formatter
    format = "music_%a - %"
    content := "music_artist - title (18/10/2020)"

    fm := Formatter {}
    err := fm.Extract(content)

    // evaluate the results
    if err == nil {
        t.Errorf("The function should have returned an error!")
    }
}
