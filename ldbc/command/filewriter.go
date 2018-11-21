package command

import (
	"encoding/base64"
	"errors"
	"io"
)

type JsonWriter struct {
	out            io.Writer
	count          int64
	keyConverter   converter
	valueConverter converter
}

type converter func(in []byte) string

func toPlainText(in []byte) string {
	return string(in)
}

func toBase64(in []byte) string {
	return base64.StdEncoding.EncodeToString(in)
}

func getConverter(name string) (converter, error) {
	switch name {
	case "raw":
		return toPlainText, nil
	case "base64":
		return toBase64, nil
	default:
		return nil, errors.New("unsupported converter " + name)
	}
}

func NewJsonWriter(out io.Writer, keyFormat string, valueFormat string) (*JsonWriter, error) {
	if out == nil {
		return nil, errors.New("out writer is nil")
	}
	writer := new(JsonWriter)
	writer.out = out
	c, e := getConverter(keyFormat)
	if e != nil {
		return nil, e
	} else {
		writer.keyConverter = c
	}
	c, e = getConverter(valueFormat)
	if e != nil {
		return nil, e
	} else {
		writer.valueConverter = c
	}
	return writer, nil
}

func (w *JsonWriter) Write(key []byte, value []byte) error {
	if w.count > 0 {
		w.writeStr(",\n")
	}
	w.writeStr("{\"key\":\"" + w.keyConverter(key) + "\",\"value\":\"" + w.valueConverter(value) + "\"}")
	w.count++
	return nil
}

func (w *JsonWriter) writeStr(s string) (int, error) {
	return w.out.Write([]byte(s))
}
func (w *JsonWriter) Start() {
	(*w).writeStr("[\n")
}
func (w *JsonWriter) End() {
	(*w).writeStr("\n]")
}
