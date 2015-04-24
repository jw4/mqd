package dispatcher

import "os"

type folderQueue struct {
	folder string
	items  map[string]itemInfo
}

type itemInfo struct {
	fi        *os.FileInfo
	processed bool
	success   bool
}

func NewPickupFolderQueue(path string) (MailQueueDispatcher, error) {
	reader := &folderQueue{folder: path, items: make(map[string]itemInfo)}
	return reader, nil
}
