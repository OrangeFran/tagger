package main

import (
    "os"
    "log"

    cli "github.com/urfave/cli/v2"

    // local imports
    parser "github.com/orangefran/tagger/internal/parser"
    commands "github.com/orangefran/tagger/internal/commands"
)

var (
    dry_run bool                // only show what the code would do
    verbose bool                // specify how much the code should spit out
    debug   bool                // specify if error should be shown if failed
    format string               // how to extract/use information in the name
    target string               // use that file/directory
)

func main() {
    // create the cli app
    app := &cli.App{
        Name: "tagger",
        Usage: "tag mp3 files from the cmdline",
        Flags: []cli.Flag {
            &cli.BoolFlag {
                Name: "debug",
                Usage: "show debugging information if failed",
                Destination: &debug,
            },
        },
        Commands: []*cli.Command {
            {
                Name: "tag",
                Aliases: []string{"t"},
                Usage: "Tags files with dynamic values",
                Flags: []cli.Flag {
                    &cli.BoolFlag {
                        Name: "dry-run",
                        Usage: "only show what would be done",
                        Destination: &dry_run,
                    },
                    &cli.BoolFlag {
                        Name: "verbose",
                        Aliases: []string{"v"},
                        Usage: "add more output",
                        Destination: &verbose,
                    },
                    &cli.StringFlag {
                        Name: "target",
                        Aliases: []string{"t"},
                        Usage: "tag `TARGET`",
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
                    return commands.Tag(target, format, verbose, dry_run)
                },
            },
            {
                Name: "static",
                Aliases: []string{"s"},
                Usage: "Tags files with static values",
                Flags: []cli.Flag {
                    &cli.BoolFlag {
                        Name: "verbose",
                        Aliases: []string{"v"},
                        Usage: "add more output",
                        Destination: &verbose,
                    },
                    &cli.StringFlag {
                        Name: "target",
                        Aliases: []string{"t"},
                        Usage: "tag `TARGET`",
                        Destination: &target,
                        Required: true,
                    },
                    &cli.StringFlag {
                        Name: "artist",
                        Usage: "set artist to `ARTIST`",
                    },
                    &cli.StringFlag {
                        Name: "title",
                        Usage: "set title to `TITLE`",
                    },
                    &cli.StringFlag {
                        Name: "album",
                        Usage: "set album to `ALBUM`",
                    },
                    &cli.StringFlag {
                        Name: "year",
                        Usage: "set year to `YEAR`",
                    },
                    &cli.StringFlag {
                        Name: "genre",
                        Usage: "set genre to `GENRE`",
                    },
                },
                Action: func(c *cli.Context) error {
                    fm := parser.Formatter {
                        Artist: c.String("artist"),
                        Title: c.String("title"),
                        Album: c.String("album"),
                        Year: c.String("year"),
                        Genre: c.String("genre"),
                    }

                    return commands.Static(target, verbose, fm)
                },
            },
            {
                Name: "remove",
                Aliases: []string{"r"},
                Usage: "Removes tags",
                Flags: []cli.Flag {
                    &cli.StringFlag {
                        Name: "target",
                        Aliases: []string{"t"},
                        Usage: "remove from `TARGET`",
                        Destination: &target,
                        Required: true,
                    },
                    &cli.BoolFlag {
                        Name: "verbose",
                        Aliases: []string{"v"},
                        Usage: "add more output",
                        Destination: &verbose,
                    },
                    &cli.BoolFlag {
                        Name: "artist",
                        Usage: "remove the artist tag",
                    },
                    &cli.BoolFlag {
                        Name: "title",
                        Usage: "remove the title tag",
                    },
                    &cli.BoolFlag {
                        Name: "album",
                        Usage: "remove the album tag",
                    },
                    &cli.BoolFlag {
                        Name: "year",
                        Usage: "remove the year tag",
                    },
                    &cli.BoolFlag {
                        Name: "genre",
                        Usage: "remove the genre tag",
                    },
                },
                Action: func(c *cli.Context) error {
                    return commands.Remove(
                        target,
                        verbose,
                        c.Bool("artist"),
                        c.Bool("title"),
                        c.Bool("album"),
                        c.Bool("year"),
                        c.Bool("genre"),
                    )
                },
            },
            {
                Name: "query",
                Aliases: []string{"q"},
                Usage: "Queries tags",
                Flags: []cli.Flag {
                    &cli.StringFlag {
                        Name: "target",
                        Aliases: []string{"t"},
                        Usage: "query `TARGET`",
                        Destination: &target,
                        Required: true,
                    },
                    &cli.StringFlag {
                        Name: "format",
                        Aliases: []string{"f"},
                        Usage: "specifiy `FORMAT` to output information",
                        Destination: &format,
                        Required: true,
                    },
                    &cli.BoolFlag {
                        Name: "verbose",
                        Aliases: []string{"v"},
                        Usage: "add more output",
                        Destination: &verbose,
                    },
                },
                Action: func(c *cli.Context) error {
                    return commands.Query(target, format, verbose)
                },
            },
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        if debug {
            log.Fatal(err)
        } else {
            log.Fatal("[-] Application failed")
        }
    }
}
