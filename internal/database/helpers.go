package database

import "os"

// removeFile deletes a file from the filesystem.
func removeFile(path string) {
	os.Remove(path)
}
