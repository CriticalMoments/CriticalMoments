package appcore

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func TestCheckInvalidCacheDir(t *testing.T) {
	cache, err := newCacheWithBaseDir("/not/a/real/path/dude")
	if err == nil || cache != nil {
		t.Fatal()
	}
}

func TestCheckValidCacheDir(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal()
	}
	cache, err := newCacheWithBaseDir(dir)
	if err != nil || cache == nil {
		t.Fatal()
	}
}

func TestCheckFindExistingConfigNoEtag(t *testing.T) {
	// primary.config without etag example
	base := "/tmp/criticalmoments/testprimarycachenoetag"
	full := base + "/primary.config"
	os.MkdirAll(base, os.ModePerm)
	os.OpenFile(full, os.O_RDONLY|os.O_CREATE, 0666)
	cache, _ := newCacheWithBaseDir(base)
	path, etag := cache.existingCacheFileOfName("primary")
	if path != full || etag != "" {
		t.Fatal()
	}
}

func TestCheckFindExistingConfigWithEtag(t *testing.T) {
	// primary.config wit etag example
	base := "/tmp/criticalmoments/testprimarycacheetag"
	full := base + "/primary--etag--123456.config"
	os.MkdirAll(base, os.ModePerm)
	os.OpenFile(full, os.O_RDONLY|os.O_CREATE, 0666)
	cache, _ := newCacheWithBaseDir(base)
	path, etag := cache.existingCacheFileOfName("primary")
	if path != full || etag != "123456" {
		t.Fatal()
	}
}

func TestCheckFindInvalidConfigWithEtag(t *testing.T) {
	// primary.config without etag example
	base := "/tmp/criticalmoments/testprimarycacheinvalidname"
	full := base + "/primary-butmore.config"
	os.MkdirAll(base, os.ModePerm)
	os.OpenFile(full, os.O_RDONLY|os.O_CREATE, 0666)
	cache, _ := newCacheWithBaseDir(base)
	path, etag := cache.existingCacheFileOfName("primary")
	if path != "" || etag != "" {
		t.Fatal("found invalid cached")
	}
}

func TestDeletePriorCache(t *testing.T) {
	url := "https://storage.googleapis.com/critical-moments-test-cases/hello.config"
	configName := "primary"
	expectedEtag := "d73b04b0e696b0945283defa3eee4538"

	base := fmt.Sprintf("/tmp/criticalmoments/testdeletepriorconfig-%v", rand.Int())
	full := base + "/primary.config"
	os.MkdirAll(base, os.ModePerm)
	os.OpenFile(full, os.O_RDONLY|os.O_CREATE, 0666)
	cache, _ := newCacheWithBaseDir(base)
	path, etag := cache.existingCacheFileOfName("primary")
	if path == "" || etag != "" {
		t.Fatal()
	}
	if _, err := os.Stat(full); err != nil {
		t.Fatal("issue with manually created cache file, test invalid")
	}

	expectedPath := filepath.Join(base, fmt.Sprintf("primary--etag--%v.config", expectedEtag))
	filePath, err := cache.verifyOrFetchRemoteConfigFile(url, configName)
	if filePath != expectedPath || err != nil {
		t.Fatal("verify or fetch didn't fetch file")
	}

	if _, err := os.Stat(full); os.IsNotExist(err) {
		t.Fatal("Didn't delete prior cache file when we got a new one")
	}
}

func TestCheckCleanEtag(t *testing.T) {
	c := cleanEtag("W/\"weaktag\"")
	if c != "" {
		t.Fatal("weak etag passed clean")
	}
	c = cleanEtag("\"123456\"")
	if c != "123456" {
		t.Fatal("Valid etag failed clean")
	}
	c = cleanEtag("123456")
	if c != "" {
		t.Fatal("RFC requires quotes")
	}
	// Only character iOS doesn't allow in filename is forward slash. If we add other OSs we should add cases here
	c = cleanEtag("\"asdf/asdf\"")
	if c != "" {
		t.Fatal("Allowed invalid filename charater, colon. May fail on non MacOS systems, needs checking")
	}

}

func TestCheckNetworkFetch(t *testing.T) {
	base := fmt.Sprintf("/tmp/criticalmoments/testprimarycachenetwork-%v", rand.Int())
	url := "https://storage.googleapis.com/critical-moments-test-cases/hello.config"
	configName := "primary"
	expectedEtag := "d73b04b0e696b0945283defa3eee4538"
	os.MkdirAll(base, os.ModePerm)

	etag := fetchEtag(url)
	if etag != expectedEtag {
		t.Fatalf("ETag of helloworld doesn't match expected. Could be network issue, so test may not be fatal. Etag: %v", etag)
	}

	cache, err := newCacheWithBaseDir(base)
	if err != nil {
		t.Fatal(err)
	}
	filePath, etag := cache.existingCacheFileOfName(configName)
	if filePath != "" || etag != "" {
		t.Fatal("Cached before fetch")
	}

	filePath, err = cache.fetchAndCache(url, configName)
	if err != nil {
		t.Fatal(err)
	}
	expectedPath := filepath.Join(base, fmt.Sprintf("primary--etag--%v.config", expectedEtag))
	if filePath != expectedPath {
		t.Fatal("File not cached to expected path")
	}
	filePath, etag = cache.existingCacheFileOfName(configName)
	if filePath != expectedPath || etag != expectedEtag {
		t.Fatal("Cache failed")
	}

	d, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(d) != "helloworld\n" {
		t.Fatal("Cache content incorrect")
	}

	preVerifyOrFetchFileInfo, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}

	filePath, err = cache.verifyOrFetchRemoteConfigFile(url, configName)
	if filePath != expectedPath || err != nil {
		t.Fatal("verify or fetch didn't find cached file")
	}

	postVerifyOrFetchFileInfo, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}

	if preVerifyOrFetchFileInfo.ModTime() != postVerifyOrFetchFileInfo.ModTime() {
		t.Fatal("Fetched when cache was available")
	}
}

func TestCheckNetworkFetchMainPath(t *testing.T) {
	base := fmt.Sprintf("/tmp/criticalmoments/testprimarycachenetworkmain-%v", rand.Int())
	url := "https://storage.googleapis.com/critical-moments-test-cases/hello.config"
	configName := "primary"
	expectedEtag := "d73b04b0e696b0945283defa3eee4538"
	os.MkdirAll(base, os.ModePerm)

	cache, err := newCacheWithBaseDir(base)
	if err != nil {
		t.Fatal(err)
	}

	filePath, etag := cache.existingCacheFileOfName(configName)
	if filePath != "" || etag != "" {
		t.Fatal("Cached before fetch")
	}

	expectedPath := filepath.Join(base, fmt.Sprintf("primary--etag--%v.config", expectedEtag))
	filePath, err = cache.verifyOrFetchRemoteConfigFile(url, configName)
	if filePath != expectedPath || err != nil {
		t.Fatal("verify or fetch didn't find cached file")
	}

	filePath, etag = cache.existingCacheFileOfName(configName)
	if filePath != expectedPath || etag != expectedEtag {
		t.Fatal("Cache failed")
	}

	preVerifyOrFetchFileInfo, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}

	filePath, err = cache.verifyOrFetchRemoteConfigFile(url, configName)
	if filePath != expectedPath || err != nil {
		t.Fatal("verify or fetch didn't find cached file")
	}

	postVerifyOrFetchFileInfo, err := os.Stat(filePath)
	if err != nil {
		t.Fatal(err)
	}

	if preVerifyOrFetchFileInfo.ModTime() != postVerifyOrFetchFileInfo.ModTime() {
		t.Fatal("Fetched when cache was available")
	}
}
