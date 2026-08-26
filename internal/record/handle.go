package record

import "os"

type fileHandle struct {
	file *os.File
	path string
}

func openFileHandle(path string) (*fileHandle, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	return &fileHandle{file: f, path: path}, nil
}

func (h *fileHandle) Write(b []byte) (int, error) {
	if h.file == nil {
		return 0, os.ErrClosed
	}
	return h.file.Write(b)
}

func (h *fileHandle) Close() error {
	if h.file == nil {
		return nil
	}
	err := h.file.Close()
	h.file = nil
	return err
}

func (h *fileHandle) Path() string {
	return h.path
}
