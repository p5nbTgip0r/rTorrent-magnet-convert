package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var pattern = regexp.MustCompile("^magnet:\\?xt=urn:btih:([^&/]+).*")

func read(path string) {
	content, err := ioutil.ReadFile(path)
	if check(err) {
		return
	}

	hash, magnet, err := parse(string(content))
	if check(err) {
		return
	}

	directoryPath := filepath.Dir(path)

	err = writeTorrent(directoryPath, hash, magnet)
	if check(err) {
		return
	}

	err = os.Remove(path)
	if check(err) {
		return
	}

	log.Printf("Successfully converted magnet at \"%s\"", path)
}

// Returns torrent hash and full magnet, respectively.
// If there is no match found, then only an error is returned
func parse(data string) (hash string, magnet string, error error) {
	matches := pattern.FindStringSubmatch(data)
	if matches == nil || len(matches) != 2 {
		log.Printf("No matches")
		return "", "", errors.New("could not match with regex")
	}

	log.Printf("Matches: %s", matches)
	return matches[1], matches[0], nil
}

func writeTorrent(path string, hash string, magnet string) error {
	builder := strings.Builder{}
	// details about this can be found here:
	// https://github.com/rakshasa/rtorrent/issues/819
	builder.WriteString("d10:magnet-uri")
	builder.WriteString(strconv.Itoa(len(magnet)))
	builder.WriteString(":")
	builder.WriteString(magnet)
	builder.WriteString("e")

	data := []byte(builder.String())
	builder.Reset()
	builder.WriteString(path)
	builder.WriteRune(filepath.Separator)
	builder.WriteString(hash)
	builder.WriteString(".torrent")

	err := ioutil.WriteFile(builder.String(), data, permission)
	return err
}

func check(e error) bool {
	if e != nil {
		log.Println(e)
		return true
	}
	return false
}
