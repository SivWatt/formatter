package formatter

import (
	"fmt"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	timeRegexp      = `\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}([\+\-]\d{2}:\d{2}|Z)`
	levelRegexp     = `\[(PANI|FATA|ERRO|WARN|INFO|DEBU|TRAC)\]`
	processIDRegexp = `\[[\d\w]+\]`
	funcRegexp      = `\[[\d\w.\-_()]+\]`
	fileRegexp      = `\[[\d\w._\-:()]+\]`
	msg             = "test log"
)

var (
	logRegexp       = fmt.Sprintf("%s %s %s *", timeRegexp, levelRegexp, processIDRegexp)
	logCallerRegexp = fmt.Sprintf("%s %s %s %s %s *", timeRegexp, levelRegexp, processIDRegexp, funcRegexp, fileRegexp)
)

func init() {
	log.SetFormatter(&AppFormatter{})
}

func TestFormat(t *testing.T) {
	log.SetReportCaller(true)
	log.SetLevel(log.TraceLevel)
	log.Info(msg)
}

func TestFormat_WithTraceLevelandReportCallerTrue(t *testing.T) {
	b := &strings.Builder{}
	log.SetReportCaller(true)
	log.SetLevel(log.TraceLevel)
	log.SetOutput(b)

	log.Errorf(msg)
	assert.Contains(t, b.String(), msg)
	assert.Regexp(t, logCallerRegexp, b.String())
}

func TestFormat_WithTraceLevelandReportCallerFalse(t *testing.T) {
	b := &strings.Builder{}
	log.SetReportCaller(false)
	log.SetLevel(log.TraceLevel)
	log.SetOutput(b)

	log.Errorf(msg)
	assert.Contains(t, b.String(), msg)
	assert.Regexp(t, logRegexp, b.String())
}

func TestFormat_WithFieldStringValue(t *testing.T) {
	b := &strings.Builder{}
	key := "testkey"
	value := "testvalue"
	expectField := fmt.Sprintf("[%s:%s]", key, value)
	log.SetReportCaller(false)
	log.SetLevel(log.TraceLevel)
	log.SetOutput(b)

	log.WithField(key, value).Debug(msg)
	assert.Regexp(t, logRegexp, b.String())
	assert.Contains(t, b.String(), msg)
	assert.Contains(t, b.String(), expectField)
}

func TestFormat_WithFieldIntValue(t *testing.T) {
	b := &strings.Builder{}
	key := "testkey"
	value := 1024
	expectField := fmt.Sprintf("[%s:%d]", key, value)
	log.SetReportCaller(false)
	log.SetLevel(log.TraceLevel)
	log.SetOutput(b)

	log.WithField(key, value).Debug(msg)
	assert.Regexp(t, logRegexp, b.String())
	assert.Contains(t, b.String(), msg)
	assert.Contains(t, b.String(), expectField)
}

func TestFormat_WithFieldEmptyKey(t *testing.T) {
	b := &strings.Builder{}
	key := ""
	value := "testvalue"
	expectField := fmt.Sprintf("[%s:%s]", key, value)
	log.SetReportCaller(false)
	log.SetLevel(log.TraceLevel)
	log.SetOutput(b)

	log.WithField(key, value).Debug(msg)
	assert.Regexp(t, logRegexp, b.String())
	assert.Contains(t, b.String(), msg)
	assert.NotContains(t, b.String(), expectField)
}

func TestFormat_WithEmptyBuffer(t *testing.T) {
	f := &AppFormatter{}
	data, err := f.Format(&log.Entry{})

	assert.NoError(t, err)
	assert.Regexp(t, logRegexp, string(data))
}
