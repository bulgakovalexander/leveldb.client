package command

import (
	"errors"
	"flag"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"io"
	"os"
)

const stdout = "stdout"
const json = "json"
const param_keyPrefix = "key-prefix"
const param_keyFrom = "key-from"

type Find struct {
	Out         string
	Format      string
	KeyFormat   string
	ValueFormat string
	KeyPrefix   string
	KeyFrom     string
	KeyTo       string
}

func (r *Find) Run(db *leveldb.DB) error {
	if db == nil {
		return errors.New("db is not initialized")
	}
	if r.Out == "" {
		return errors.New("out is missed")
	}
	if r.Format != json {
		return errors.New("unsupported file format " + r.Format)
	} else {
		var writer io.Writer
		if r.Out == stdout {
			writer = os.Stdout
		} else {
			file, e := os.Create(r.Out)
			if e != nil {
				return e
			} else {
				writer = file
				defer file.Close()
			}
		}

		outWriter, e := NewJsonWriter(writer, r.KeyFormat, r.ValueFormat)
		if e != nil {
			return e
		}
		outWriter.Start()

		filter, e := r.getFilter()
		if e != nil {
			return e
		}
		iterator := db.NewIterator(filter, nil)
		for iterator.Next() {
			key := iterator.Key()
			value := iterator.Value()
			outWriter.Write(key, value)
		}
		iterator.Release()
		e = iterator.Error()
		if e != nil {
			return e
		}

		outWriter.End()
		return nil
	}
}

func (r *Find) getFilter() (*util.Range, error) {
	var filter *util.Range = nil
	if r.KeyPrefix != "" {
		filter = util.BytesPrefix([]byte(r.KeyPrefix))
	}
	if r.KeyFrom != "" {
		if filter != nil {
			return nil, errors.New("don'n use " + param_keyPrefix + " and " + param_keyFrom + " together")
		}
		start := []byte(r.KeyFrom)
		var limit []byte = nil
		if r.KeyTo != "" {
			limit = []byte(r.KeyTo)
		}
		filter = &util.Range{Start: start, Limit: limit}
	}
	return filter, nil
}

func (r *Find) Parse(args []string) error {
	sliced := args[3:]
	options := flag.NewFlagSet("", flag.ExitOnError)
	options.StringVar(&r.Out, "out", stdout, "output file")
	options.StringVar(&r.Format, "file-format", json, "output format")
	options.StringVar(&r.KeyFormat, "key-format", "raw", "key format")
	options.StringVar(&r.ValueFormat, "value-format", "base64", "value format")
	options.StringVar(&r.KeyPrefix, param_keyPrefix, "", "filter by key")
	options.StringVar(&r.KeyFrom, param_keyFrom, "", "strict result by first included key")
	options.StringVar(&r.KeyTo, "key-to", "", "strict result by last excluded key")

	e := options.Parse(sliced)
	return e
}