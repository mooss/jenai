package workspace

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// readContent returns an io.Reader for the given path, which can be a local file or a URL.
// It also returns a cleanup function to be called when the reader is no longer needed.
func readContent(path string) (io.ReadCloser, error) {
	if isURL(path) {
		return downloadContent(path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// isURL checks if the given string is a valid URL.
func isURL(path string) bool {
	u, err := url.Parse(path)
	return err == nil && u.Host != "" && (u.Scheme == "http" || u.Scheme == "https")
}

// downloadContent downloads content from the given URL.
// The caller must close the returned reader.
func downloadContent(fileURL string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, fileURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download URL %s: %w", fileURL, err)
	}

	client := *http.DefaultClient
	client.Timeout = 2 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download URL %s: %w", fileURL, err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		msg := "failed to download URL %s, status code: %d"
		return nil, fmt.Errorf(msg, fileURL, resp.StatusCode)
	}

	return resp.Body, nil
}

// collectDirFiles returns every regular file under root.
func collectDirFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
