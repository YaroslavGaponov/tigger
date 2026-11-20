package watcher

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"
)

type Handler func(path string)

type FileWatcher struct {
	watcher  *fsnotify.Watcher
	onCreate []Handler
	onRemove []Handler
}

func NewFileWatcher() (*FileWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &FileWatcher{
		watcher: w,
	}, nil
}

func (f *FileWatcher) OnCreate(handler Handler) {
	f.onCreate = append(f.onCreate, handler)
}

func (f *FileWatcher) OnRemove(handler Handler) {
	f.onRemove = append(f.onRemove, handler)
}

func (f *FileWatcher) Watch(ctx context.Context, paths []string) error {
	for _, path := range paths {
		if err := f.watcher.Add(path); err != nil {
			return err
		}
	}

	go f.listen(ctx)
	return nil
}

func (f *FileWatcher) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case event, ok := <-f.watcher.Events:
			if !ok {
				return
			}

			switch {
			case event.Has(fsnotify.Create):
				for _, handler := range f.onCreate {
					handler(event.Name)
				}

			case event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename):
				for _, handler := range f.onRemove {
					handler(event.Name)
				}
			}

		case err, ok := <-f.watcher.Errors:
			if !ok {
				return
			}
			log.Println("watcher error:", err)
		}
	}
}

func (f *FileWatcher) Close() error {
	return f.watcher.Close()
}
