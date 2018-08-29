package main

import (
	"testing"
)

func TestFileManager(t *testing.T) {
	testIO := []struct {
		name string
	}{
		{
			name: "a test",
		},
	}

	for _, test := range testIO {
		t.Run(test.name, func(t *testing.T) {
			manager := NewFileManager()
			manager.Manage()
		})
	}
}
