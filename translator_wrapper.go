package main

import (
	"context"
	"golang.org/x/text/language"

	"time"
	"math"
	"errors"
	"fmt"
)

type translatorWrapperCacheEntry struct {
	res string
	created_at time.Time
}

//	translatorWrapper implements Translator interface
//	and wraps another Translate instance to add some middleware functionality
type translatorWrapper struct {
	translator Translator
	max_attempts int
	cache map[string]translatorWrapperCacheEntry
	cache_lifetime time.Duration
}

//	Attempting `max_attempts` retries before giving up.
//	Delay between retries is 10 times bigger every time.
//	Results are cached for `cache_lifetime`.
func (t translatorWrapper) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	//	TODO : better key
	cache_key := fmt.Sprintf("%v#%v#%s", from, to, data);

	if val, ok := t.cache[cache_key]; ok {
		if (time.Since(val.created_at) < t.cache_lifetime) {
			return val.res, nil
		} else {
			delete(t.cache, cache_key)
		}
	}

	var res string
	var err error
	for i := 1; i <= t.max_attempts; i++ {
		res, err = t.translator.Translate(ctx, from, to, data)

		if (err == nil) {
			t.cache[cache_key] = translatorWrapperCacheEntry{res, time.Now()}

			return res, nil
		} else if (i < t.max_attempts) {
			time.Sleep(time.Duration(int64(math.Pow10(i))*int64(time.Millisecond)))
		}
	}

	return "", errors.New(fmt.Sprintf("%d attempts failed with error \"%s\"", t.max_attempts, err.Error()))
}

func newTranslatorWrapper(t Translator, max_attempts int, lifetime time.Duration) *translatorWrapper {
	return &translatorWrapper{
		translator: t,
		max_attempts: max_attempts,
		cache: map[string]translatorWrapperCacheEntry{},
		cache_lifetime: lifetime,
	}
}