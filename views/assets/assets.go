package assets

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"io/fs"
	"sync"
)

//go:embed css/* js/*
var Assets embed.FS

var (
	versionOnce sync.Once
	version     string
)

// URL resolves an embedded asset relative to the assets dir (e.g.
// "css/output.css") to a public path with a content-hash query, so CDN and
// browser caches are busted whenever the bytes change across deploys.
func URL(rel string) string {
	return "/assets/" + rel + "?v=" + Version()
}

// Version is a short content hash of the mutable asset bundles (CSS + JS).
func Version() string {
	versionOnce.Do(func() {
		v, err := computeVersion()
		if err != nil {
			panic("assets: could not hash version: " + err.Error())
		}
		version = v
	})
	return version
}

// computeVersion hashes the CSS and JS bundles in walk order (embed.FS reads
// are sorted, so it is deterministic).
func computeVersion() (string, error) {
	h := sha256.New()
	for _, dir := range []string{"css", "js"} {
		err := fs.WalkDir(Assets, dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			data, err := Assets.ReadFile(path)
			if err != nil {
				return err
			}
			h.Write([]byte(path))
			h.Write(data)
			return nil
		})
		if err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(h.Sum(nil))[:12], nil
}
