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
    format string               // how to extract/use information in the name
    target string               // use that file/directory
)

func main() {
    // create the cli app
    app := &cli.App{
        Name: "tagger",
        Usage: "tag mp3 files from the cmdline",
        Commands: []*cli.Command {
            {
                Name: "tag",
                Aliases: []string{"t"},
                Usage: "Tags files with dynamic values",
                Flags: []cli.Flag {
                    &cli.BoolFlag {
                        Name: "dry-run",
                        Usage: "Only show what would be done",
                        Destination: &dry_run,
                    },
                    &cli.BoolFlag {
                        Name: "verbose",
                        Aliases: []string{"v"},
                        Usage: "Adds more output",
                        Destination: &verbose,
                    },
                    &cli.StringFlag {
                        Name: "target",
                        Aliases: []string{"t"},
                        Usage: "Tags `TARGET`",
                        Destination: &target,
                        Required: true,
                    },
                    &cli.StringFlag {
                        Name: "format",
                        Aliases: []string{"f"},
                        Usage: "Specifies `FORMAT`; used to tag files",
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
                        Usage: "Adds more output",
                        Destination: &verbose,
                    },
                    &cli.StringFlag {
                        Name: "target",
                        Aliases: []string{"t"},
                        Usage: "Tags `TARGET`",
                        Destination: &target,
                        Required: true,
                    },
                    &cli.StringFlag {
                        Name: "artist",
                        Usage: "Sets artist to `ARTIST`",
                    },
                    &cli.StringFlag {
                        Name: "title",
                        Usage: "Sets title to `TITLE`",
                    },
                    &cli.StringFlag {
                        Name: "album",
                        Usage: "Sets album to `ALBUM`",
                    },
                    &cli.StringFlag {
                        Name: "year",
                        Usage: "Sets year to `YEAR`",
                    },
                    &cli.StringFlag {
                        Name: "genre",
                        Usage: "Sets genre to `GENRE`",
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
                        Usage: "Removes tags from `TARGET`",
                        Destination: &target,
                        Required: true,
                    },
                    &cli.BoolFlag {
                        Name: "verbose",
                        Aliases: []string{"v"},
                        Usage: "Adds more output",
                        Destination: &verbose,
                    },
                    &cli.BoolFlag {
                        Name: "artist",
                        Usage: "Removes the artist tag",
                    },
                    &cli.BoolFlag {
                        Name: "title",
                        Usage: "Removes the title tag",
                    },
                    &cli.BoolFlag {
                        Name: "album",
                        Usage: "Removes the album tag",
                    },
                    &cli.BoolFlag {
                        Name: "year",
                        Usage: "Removes the year tag",
                    },
                    &cli.BoolFlag {
                        Name: "genre",
                        Usage: "Removes the genre tag",
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
                Usage: "Query tags",
                Flags: []cli.Flag {
                    &cli.StringFlag {
                        Name: "target",
                        Aliases: []string{"t"},
                        Usage: "Queries `TARGET`",
                        Destination: &target,
                        Required: true,
                    },
                    &cli.StringFlag {
                        Name: "format",
                        Aliases: []string{"f"},
                        Usage: "Specifies `FORMAT` to output information",
                        Destination: &format,
                        Required: true,
                    },
                    &cli.BoolFlag {
                        Name: "verbose",
                        Aliases: []string{"v"},
                        Usage: "Adds more output",
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
        log.Fatal(err)
    }
}
