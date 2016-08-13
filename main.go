package main

import (
	"crypto/sha256"
	"io"
	"log"
	"os"
	"bytes"
	"time"
	"syscall"
	"github.com/kardianos/osext"
)

func main() {
	binpath, err := osext.Executable()
	if err != nil {
		log.Fatal(err)
	}

	lasthash, err := hashfile(binpath)
	if err != nil {
		log.Fatal(err)
	}

	// Wait until the hash changes
	for {
		hash, err := hashfile(binpath)
		if err != nil {
			log.Fatal(err)
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
			log.Fatal(err)
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
		log.Fatal(err)
	}
}

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
