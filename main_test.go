package main

import (
	"context"
	"math/rand"
	"time"

	"golang.org/x/text/language"

	"testing"

	"sync"
)

func TestSingle(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	s := NewService()
	_, err := s.translator.Translate(ctx, language.English, language.Japanese, "test")

	if (err != nil) {
		t.Errorf("Got error: \"%v\"", err)
	}
}

func TestSequence(t *testing.T) {
	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	s := NewService()
	res, _ := s.translator.Translate(ctx, language.English, language.Japanese, "test")
	res2, _ := s.translator.Translate(ctx, language.English, language.Japanese, "test")

	if (res != res2) {
		t.Errorf("Results does not match: got \"%s\" and \"%s\"", res, res2)
	}
}

func TestConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	var res [5]string

	ctx := context.Background()
	rand.Seed(time.Now().UTC().UnixNano())
	s := NewService()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(s *Service, ctx context.Context, str *string) {
			defer wg.Done()

			res, err := s.translator.Translate(ctx, language.English, language.Japanese, "test")

			if (err == nil) {
				*str = res
			}
		}(s, ctx, &res[i])
	}

	wg.Wait()

	for i := 1; i < 5; i++ {
		if (res[0] != res[i]) {
			t.Errorf("Results does not match: %v", res)
			break
		}
	}
}