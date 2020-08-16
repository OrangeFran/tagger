package main

import (
    "testing"
)

func TestParseSimple(t *testing.T) {
    // create a simple formatter
    format := "%a - %t"
    content := "artist - title"

    fm := Formatter {}
    fm.Extract(content, format)

    // evaluate the results
    if fm.Artist != "artist" || fm.Title != "title" {
        t.Errorf("Simple test failed!\nWanted %s & %s, got %s & %s\n", "artist", "title", fm.Artist, fm.Title)
    }
}

func TestParseBeginning(t *testing.T) {
    // create a simple formatter
    format := "music_%a - %t (%y)"
    content := "music_artist - title (18/10/2020)"

    fm := Formatter {}
    fm.Extract(content, format)

    // evaluate the results
    if fm.Artist != "artist" || fm.Title != "title" || fm.Year != "18/10/2020" {
        t.Errorf("Beginning test failed!\nWanted %s & %s & %s, got %s & %s & %s\n", "artist", "title", "18/10/2020", fm.Artist, fm.Title, fm.Year)
    }
}

