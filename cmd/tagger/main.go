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
    format string = ""          // how to extract/use information in the name
    verbose bool                // specify how much the code should spit out
    target string = ""          // use that file/directory
)

func main() {
    // create the cli-flags and more
    app := &cli.App{
        Name: "tagger",
        Usage: "tag mp3 files from the cmdline",
        Commands: []*cli.Command {
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
                },
                Action: func(c *cli.Context) error {
                    return commands.Query(target)
                },
            },
            {
                Name: "tag",
                Aliases: []string{"t"},
                Usage: "Tags files",
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
                Name: "manually",
                Aliases: []string{"m"},
                Usage: "Tag with manual values",
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
                        Usage: "Specifies artist manually",
                    },
                    &cli.StringFlag {
                        Name: "title",
                        Usage: "Specifies title manually",
                    },
                    &cli.StringFlag {
                        Name: "album",
                        Usage: "Specifies album manually",
                    },
                    &cli.StringFlag {
                        Name: "year",
                        Usage: "Specifies year manually",
                    },
                    &cli.StringFlag {
                        Name: "genre",
                        Usage: "Specifies genre manually",
                    },
                },
                Action: func(c *cli.Context) error {
                    fm := parser.Formatter {
                        c.String("artist"),
                        c.String("title"),
                        c.String("album"),
                        c.String("year"),
                        c.String("genre"),
                    }

                    return commands.Manually(target, verbose, fm)
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
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}
