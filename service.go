package main

import "time"

// Service is a Translator user.
type Service struct {
	translator Translator
}

func NewService() *Service {
	t := newRandomTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.9,
	)

	w := newTranslatorWrapper(
		t,
		4,
	)

	return &Service{
		translator: w,
	}
}
