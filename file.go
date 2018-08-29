package main

import (
	"log"
	"os"
	"sync"
)

type FileOpener struct {
	results chan *string
}

func NewFileOpener() *FileOpener {
	return &FileOpener{
		results: make(chan *string),
	}
}

func (f *FileOpener) FindFiles(fileChan chan *os.File) {
	walk(f.results)

	for {
		select {
		case path := <-f.results:
			actualFile, err := os.Open(*path)
			if err != nil {
				log.Println(err.Error())
			}
			fileChan <- actualFile
		default:
			fileChan <- nil
			return
		}

	}
}

type FileManager struct {
	files    map[*os.File]int64
	mu       sync.Mutex
	fileChan chan *os.File
	opener   *FileOpener
}

func NewFileManager() *FileManager {
	return &FileManager{
		files:    make(map[*os.File]int64),
		fileChan: make(chan *os.File),
		opener:   NewFileOpener(),
	}
}

func (m *FileManager) Manage() {
	go m.opener.FindFiles(m.fileChan)
	for {
		select {
		case newFile := <-m.fileChan:
			if newFile == nil {
				return
			}
			m.insert(newFile)
		}
	}
}

func (m *FileManager) insert(f *os.File) {
	m.mu.Lock()
	defer m.mu.Unlock()

	stat, err := f.Stat()
	if err != nil {
		log.Printf("unable to stat file %q. skipping.\n", f.Name())
		return
	}

	m.files[f] = stat.Size()

	if len(m.files) > 10 {
		var smallest *os.File
		for file, size := range m.files {
			if smallest == nil || size < m.files[smallest] {
				smallest = file
			}
		}

		defer smallest.Close()
		delete(m.files, smallest)
	}
}

func walk(chan *string) error {
	return nil
}
