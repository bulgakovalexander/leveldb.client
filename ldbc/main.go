package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"gopkg.in/urfave/cli.v1"
	"leveldb.client/ldbc/command"
	"log"
	"os"
)

const version = "0.0.0.1"

func main() {
	app := New()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type LevelDbApp struct {
	*cli.App
	dbPath string
	db     *leveldb.DB
}

func New() *LevelDbApp {
	levelDbCmdLine := &LevelDbApp{App: cli.NewApp(), dbPath: "", db: nil}
	levelDbCmdLine.Version = version
	levelDbCmdLine.Usage = "LevelDB command line client"
	//levelDbCmdLine.UsageText = levelDbCmdLine.HelpName + " dbPath command [command options]"
	levelDbCmdLine.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "db",
			Usage:       "path to database",
			Value:       "./",
			Destination: &levelDbCmdLine.dbPath,
		},
	}
	levelDbCmdLine.Commands = []cli.Command{newFind(&levelDbCmdLine.dbPath)}
	return levelDbCmdLine
}

func (app *LevelDbApp) Run(args []string) error {
	err := app.App.Run(args)
	return err
}

func newFind(dbPath *string) cli.Command {

	cmd := &command.Find{}
	return cli.Command{
		Name:  "find",
		Usage: "find values in the database",
		Action: func(c *cli.Context) error {
			db, e := leveldb.OpenFile(*dbPath, &opt.Options{
				ErrorIfMissing: true,
				ReadOnly:       true,
			})
			if e != nil {
				return e
			}
			defer db.Close()
			return cmd.Run(db)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "out",
				Usage:       "output destination filename or stdout",
				Value:       "stdout",
				Destination: &cmd.Out,
			},
			cli.StringFlag{
				Name:        "file-format",
				Usage:       "a format of the output",
				Value:       "json",
				Destination: &cmd.Format,
			},
			cli.StringFlag{
				Name:        "key-format",
				Usage:       "raw or base64",
				Value:       "raw",
				Destination: &cmd.KeyFormat,
			},
			cli.StringFlag{
				Name:        "value-format",
				Usage:       "raw or base64",
				Value:       "base64",
				Destination: &cmd.ValueFormat,
			},
			cli.StringFlag{
				Name:        "key-prefix",
				Usage:       "filtering found values by a key prefix",
				Destination: &cmd.KeyPrefix,
			},
			cli.StringFlag{
				Name:        "key-from",
				Destination: &cmd.KeyFrom,
			},
			cli.StringFlag{
				Name:        "key-to",
				Destination: &cmd.KeyTo,
			},
		},
	}
}
