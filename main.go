package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var (
	out string
	dir string
)

func init() {
	flag.StringVar(&out, "out", "index.json", "Output file name of index")
	flag.StringVar(&dir, "dir", ".", "Directory to create index for")
}

func main() {
	flag.Parse()

	dict := make(map[string]string)
	filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		err, hash := hash(path)
		if err != nil {
			log.Printf("Failed hashing %s: %s", path, err)
			return nil
		}

		log.Printf("%s %s", hash, path)
		dict[hash] = path
		return nil
	})

	json, err := json.MarshalIndent(dict, "", "  ")
	if err != nil {
		panic(err)
	}

	if os.WriteFile(out, json, 0644) != nil {
		panic(err)
	}
}

func hash(file string) (error, string) {
	handle, err := os.Open(file)
	if err != nil {
		return err, ""
	}

	defer handle.Close()

	h := md5.New()
	if _, err := io.Copy(h, handle); err != nil {
		return err, ""
	}

	output := h.Sum(nil)
	str := hex.EncodeToString(output)
	return nil, str
}
