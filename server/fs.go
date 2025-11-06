// Add this custom filesystem type
package server

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
)

type restrictedFS struct {
	root string
}

func (fs *restrictedFS) Filelist(r *sftp.Request) (sftp.ListerAt, error) {
    // Resolve the requested path within our root
    path := filepath.Join(fs.root, r.Filepath)
    
    fmt.Printf("🔍 Filelist request: %s -> %s\n", r.Filepath, path)
    
    // Security check: make sure the resolved path is still within root
    cleanPath := filepath.Clean(path)
    cleanRoot := filepath.Clean(fs.root)
    if !strings.HasPrefix(cleanPath, cleanRoot) {
        fmt.Printf("❌ Permission denied: path outside root\n")
        return nil, os.ErrPermission
    }
    
    // Check if it's a file or directory
    info, err := os.Stat(path)
    if err != nil {
        fmt.Printf("❌ Stat error: %v\n", err)
        return nil, err
    }
    
    // If it's a file, return just that file's info
    if !info.IsDir() {
        fmt.Printf("✅ Returning single file info\n")
        return listerat([]os.FileInfo{info}), nil
    }
    
    // It's a directory, list its contents
    f, err := os.Open(path)
    if err != nil {
        fmt.Printf("❌ Directory open error: %v\n", err)
        return nil, err
    }
    
    list, err := f.Readdir(0)
    if err != nil {
        f.Close()
        fmt.Printf("❌ Readdir error: %v\n", err)
        return nil, err
    }
    f.Close()
    
    fmt.Printf("✅ Listed %d files\n", len(list))
    return listerat(list), nil
}

func (fs *restrictedFS) Fileread(r *sftp.Request) (io.ReaderAt, error) {
    path := filepath.Join(fs.root, r.Filepath)
    
    fmt.Printf("🔍 Fileread request: %s -> %s\n", r.Filepath, path)
    
    // Security check
    cleanPath := filepath.Clean(path)
    cleanRoot := filepath.Clean(fs.root)
    if !strings.HasPrefix(cleanPath, cleanRoot) {
        fmt.Printf("❌ Permission denied: path outside root\n")
        return nil, os.ErrPermission
    }
    
    // Check if file exists
    if _, err := os.Stat(path); err != nil {
        fmt.Printf("❌ File stat error: %v\n", err)
        return nil, err
    }
    
    file, err := os.Open(path)
    if err != nil {
        fmt.Printf("❌ File open error: %v\n", err)
        return nil, err
    }
    
    fmt.Printf("✅ File opened successfully: %s\n", path)
    return file, nil
}

func (fs *restrictedFS) Filewrite(r *sftp.Request) (io.WriterAt, error) {
	return nil, os.ErrPermission // Read-only
}

func (fs *restrictedFS) Filecmd(r *sftp.Request) error {
	return os.ErrPermission // Read-only
}

// Helper type for directory listings
type listerat []os.FileInfo

func (f listerat) ListAt(ls []os.FileInfo, offset int64) (int, error) {
	if offset >= int64(len(f)) {
		return 0, io.EOF
	}
	n := copy(ls, f[offset:])
	if n < len(ls) {
		return n, io.EOF
	}
	return n, nil
}
