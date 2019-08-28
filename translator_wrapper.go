package main

import (
	"context"
	"golang.org/x/text/language"

	"time"
	"math"
	"errors"
	"fmt"
)

//	translatorWrapper implements Translator interface
//	and wraps another Translate instance to add some middleware functionality
type translatorWrapper struct {
	translator Translator
	max_attempts int
}

//	Attempting `max_attempts` retries before giving up.
//	Delay between retries is 10 times bigger every time.
func (t translatorWrapper) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	var res string
	var err error
	for i := 1; i <= t.max_attempts; i++ {
		res, err = t.translator.Translate(ctx, from, to, data)

		if (err == nil) {
			return res, nil
		} else if (i < t.max_attempts) {
			time.Sleep(time.Duration(int64(math.Pow10(i))*int64(time.Millisecond)))
		}
	}

	return "", errors.New(fmt.Sprintf("%d attempts failed with error \"%s\"", t.max_attempts, err.Error()))
}

func newTranslatorWrapper(t Translator, max_attempts int) *translatorWrapper {
	return &translatorWrapper{
		translator: t,
		max_attempts: max_attempts,
	}
}