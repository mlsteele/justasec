package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/kardianos/osext"
	"io"
	"os"
	"syscall"
	"time"
)

func main() {
	binpath, err := osext.Executable()
	if err != nil {
		die(err)
	}

	lasthash, err := hashfile(binpath)
	if err != nil {
		die(err)
	}

	// Wait until the hash changes
	for {
		hash, err := hashfile(binpath)
		if err != nil {
			die(err)
		}
		// log.Printf("hash of %v: %v", binpath, base64.StdEncoding.EncodeToString(hash))

		if !bytes.Equal(hash, lasthash) {
			lasthash = hash
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	// Wait until the hash stabilizes
	for {
		hash, err := hashfile(binpath)
		if err != nil {
			die(err)
		}

		if bytes.Equal(hash, lasthash) {
			break
		}
		lasthash = hash

		time.Sleep(10 * time.Millisecond)
	}

	// Exec the new binary
	err = syscall.Exec(binpath, os.Args, nil)
	if err != nil {
		die(err)
	}
}

// Get the sha256 hash of the file contents at path.
func hashfile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "justasec failed: %v\n", err)
	os.Exit(1)
}
