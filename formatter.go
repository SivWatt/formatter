package formatter

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	defaultTimestampFormat = time.RFC3339
	FieldKeyFile           = "file"
	FieldKeyFunc           = "func"
)

var _ log.Formatter = &myFormatter{}

type myFormatter struct {
	DisableTimestamp bool
	DisablePID       bool
}

// Format implements interface
func (f *myFormatter) Format(entry *log.Entry) ([]byte, error) {
	// initial buffer if necessary
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// output fixed fields
	// timestamp
	if !f.DisableTimestamp {
		f.appendWithoutKey(b, entry.Time.Local().Format(defaultTimestampFormat))
	}

	// log level
	f.appendFieldLevel(b, entry.Level.String())

	// process id
	if !f.DisablePID {
		f.appendProcessID(b, os.Getpid())
	}

	if entry.HasCaller() {
		// function field
		f.appendBracketsValue(b, path.Base(entry.Caller.Function)+"()")
		// file field
		f.appendBracketsValue(b, path.Base(entry.Caller.File)+":"+strconv.Itoa(entry.Caller.Line))
	}

	// message
	trimmed := strings.TrimSpace(entry.Message)
	if trimmed != "" {
		f.appendWithoutKey(b, trimmed)
	}

	// sorting additional fields
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		// remove the key-value pair which key is an empty or space string.
		// to avoid from the weird log text: " = xxx value" which has no key string
		t := strings.TrimSpace(k)
		if t == "" {
			continue
		}

		keys = append(keys, k)
	}
	sort.Strings(keys)

	// output additional fields
	for _, k := range keys {
		f.appendKeyValue(b, k, entry.Data[k])
	}

	// given line break at the end
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *myFormatter) appendWithoutKey(b *bytes.Buffer, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}

	f.appendValue(b, value)
}

func (f *myFormatter) appendFieldLevel(b *bytes.Buffer, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}

	f.appendBracketsValue(b, fmt.Sprintf("%s", strings.ToUpper(value.(string))[:4]))
}

func (f *myFormatter) appendProcessID(b *bytes.Buffer, value int) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}

	f.appendBracketsValue(b, fmt.Sprintf("%.4x", value))
}

func (f *myFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}

	b.WriteByte('[')
	b.WriteString(key)
	b.WriteByte(':')
	f.appendValue(b, value)
	b.WriteByte(']')
}

func (f *myFormatter) appendBracketsValue(b *bytes.Buffer, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}

	b.WriteByte('[')
	f.appendValue(b, value)
	b.WriteByte(']')
}

func (f *myFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	b.WriteString(stringVal)
}
