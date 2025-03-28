package workers

import (
	"fmt"
	reflect "reflect"
	"strings"
	"sync"
)

const subsKeyDelimiter = "."

type workerBase struct {
	subscriptions sync.Map // symbol -> SubscriptionData
}

type SubscriptionData struct {
	Service      Unsubscriber
	ErrorHandler func(error)
}

func (w *workerBase) Stop() {
	w.UnsubscribeAll()
}

func getSubsKey(args ...string) string {
	if len(args) == 0 {
		return "empty"
	}

	return strings.Join(args, subsKeyDelimiter)
}

func (w *workerBase) IsSubscriptionExists(args ...string) bool {
	key := getSubsKey(args...)

	_, isExists := w.subscriptions.Load(key)
	return isExists
}

func (w *workerBase) Save(
	unsubscriber Unsubscriber,
	errorHandler func(error),
	args ...string,
) {
	key := getSubsKey(args...)

	w.subscriptions.Store(key, SubscriptionData{
		Service:      unsubscriber,
		ErrorHandler: errorHandler,
	})
}

func (w *workerBase) Unsubscribe(args ...string) {
	key := getSubsKey(args...)

	iSub, isExists := w.subscriptions.Load(key)
	if !isExists {
		return
	}

	subsData, isConvertable := iSub.(SubscriptionData)
	if !isConvertable {
		fmt.Printf(
			"unsubscribe: get subs data: unknown format: %s\n",
			reflect.ValueOf(iSub).String(),
		)
		return
	}

	// stop service
	if subsData.Service != nil {
		if err := subsData.Service.Unsubscribe(); err != nil && subsData.ErrorHandler != nil {
			subsData.ErrorHandler(fmt.Errorf(
				"unsubscribe %q: %w",
				key, err,
			))
		}
	}

	// remove subscription data
	w.subscriptions.Delete(key)
}

func (w *workerBase) getSubArgs() [][]string {
	var args [][]string
	w.subscriptions.Range(func(ikey, _ any) bool {
		key, isConvertable := ikey.(string)
		if !isConvertable {
			return true
		}

		args = append(args, strings.Split(key, subsKeyDelimiter))
		return true
	})
	return args
}

func (w *workerBase) UnsubscribeAll() {
	args := w.getSubArgs()
	for _, subArgs := range args {
		w.Unsubscribe(subArgs...)
	}
}
