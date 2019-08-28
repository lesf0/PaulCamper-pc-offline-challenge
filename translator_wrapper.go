package main

import (
	"context"
	"golang.org/x/text/language"
)

//	translatorWrapper implements Translator interface
//	and wraps another Translate instance to add some middleware functionality
type translatorWrapper struct {
	translator Translator
}

//	Just passing the data to wrapped instance
func (t translatorWrapper) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	return t.translator.Translate(ctx, from, to, data)
}

func newTranslatorWrapper(t Translator) *translatorWrapper {
	return &translatorWrapper{
		translator: t,
	}
}