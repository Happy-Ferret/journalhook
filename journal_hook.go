package journal_hook

import (
	"fmt"
	logrus "github.com/Sirupsen/logrus"
	"github.com/coreos/go-systemd/journal"
	"io/ioutil"
)

type JournaldHook struct{}

var (
	severityMap = map[logrus.Level]journal.Priority{
		logrus.DebugLevel: journal.PriDebug,
		logrus.InfoLevel:  journal.PriInfo,
		logrus.WarnLevel:  journal.PriWarning,
		logrus.ErrorLevel: journal.PriErr,
		logrus.FatalLevel: journal.PriCrit,
		logrus.PanicLevel: journal.PriEmerg,
	}
)

// Journal wants strings but logrus takes anything.
func stringifyEntries(data map[string]interface{}) map[string]string {
	entries := make(map[string]string)
	for k, v := range data {
		// Journal wants uppercase strings.
		entries[k] = fmt.Sprintf("%v", v)
	}
	return entries
}

func (hook *JournaldHook) Fire(entry *logrus.Entry) error {
	return journal.Send(entry.Message, severityMap[entry.Level], stringifyEntries(entry.Data))
}

// `Levels()` returns a slice of `Levels` the hook is fired for.
func (hook *JournaldHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

// Adds the Journal hook if journal is enabled
// Sets log output to ioutil.Discard so stdout isn't captured.
func EnableJournal() {
	if !journal.Enabled() {
		logrus.Warning("Journal not available but user requests we log to it. Ignoring")
	} else {
		logrus.AddHook(&JournaldHook{})
		logrus.SetOutput(ioutil.Discard)
	}
}
