package main

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

const yamlParseTimeout = 2 * time.Second

func ParseYAML(content string) (*ResumeData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), yamlParseTimeout)
	defer cancel()

	type result struct {
		data ResumeData
		err  error
	}
	ch := make(chan result, 1)

	go func() {
		var data ResumeData
		err := yaml.Unmarshal([]byte(content), &data)
		ch <- result{data, err}
	}()

	select {
	case r := <-ch:
		if r.err != nil {
			return nil, fmt.Errorf("yaml parse: %w", r.err)
		}
		return &r.data, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("yaml parse exceeded %s timeout", yamlParseTimeout)
	}
}
