package appcore

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type cache struct {
	baseDirectory string
}

const (
	configFileSuffix = ".config"
	etagDelim        = "--etag--"
)

func newCacheWithBaseDir(cacheDirPath string) (*cache, error) {
	// validate cache dir exists
	if _, err := os.Stat(cacheDirPath); os.IsNotExist(err) {
		return nil, errors.New("CriticalMoments: Cache directory does not exist")
	}

	cache := cache{
		baseDirectory: cacheDirPath,
	}

	return &cache, nil
}

func (c *cache) verifyOrFetchRemoteConfigFile(url string, configFileName string) (filepath string, err error) {
	// filename: primary--etag--[ETAG].config if etag, if not primary.config
	// only one "primary*.config" at a time.

	// find existing config in cache
	existingCached, existingEtag := c.existingCacheFileOfName(configFileName)
	if existingEtag != "" {
		// validate etag hasn't changed, so we don't have to request full file
		currentEtag := fetchEtag(url)
		if currentEtag == existingEtag {
			return existingCached, nil
		} else {
			// existing cached is no longer valid
			os.Remove(existingCached)
		}
	}

	newCached, err := c.fetchAndCache(url, configFileName)
	if err != nil {
		return "", nil
	}

	return newCached, nil

}

func (c *cache) existingCacheFileOfName(fileName string) (path string, etag string) {
	cacheFiles, err := os.ReadDir(c.baseDirectory)
	if err != nil {
		fmt.Println("CriticalMoments: Could not read cache directory")
		return "", ""
	}

	for _, file := range cacheFiles {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		fullFilePath := filepath.Join(c.baseDirectory, name)
		if name == fileName+configFileSuffix {
			// no etag, exact match
			return fullFilePath, ""
		}
		// check for etag with NAME--etag--ETAG.config format
		if strings.HasPrefix(name, fileName) && strings.HasSuffix(name, configFileSuffix) && strings.Index(name, etagDelim) == len(fileName) {
			etagEnd := len(name) - len(configFileSuffix)
			etag := name[len(fileName)+len(etagDelim) : etagEnd]
			return fullFilePath, etag
		}

	}
	return "", ""
}

func fetchEtag(url string) string {
	var client = &http.Client{
		Timeout: time.Second * 5,
	}

	response, err := client.Head(url)
	if err != nil {
		return ""
	}

	if response.StatusCode != http.StatusOK {
		return ""
	}

	return cleanEtag(response.Header.Get("ETag"))
}

func (c *cache) fetchAndCache(url string, fileName string) (cachedFile string, err error) {
	var client = &http.Client{
		Timeout: time.Second * 20,
	}

	response, err := client.Get(url)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", errors.New("Failed to fetch config file")
	}

	cacheFileName := fileName + configFileSuffix
	etagRaw := response.Header.Get("ETag")
	etag := cleanEtag(etagRaw)
	if etag != "" {
		cacheFileName = fileName + etagDelim + etag + configFileSuffix
	} else {
		fmt.Println("CriticalMoments: Warning -- your host is not returning an etag header for your config file. This will mean (slightly) more network traffic. We suggest you use a host which supports etag.")
	}

	cacheFileFullPath := filepath.Join(c.baseDirectory, cacheFileName)

	// Read into memory and write to file atomically. We don't want a file
	// with an etag to be written, that isn't complete.
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// Jump though hoops of write then move to make this atomic (via Rename)
	tmpCache := filepath.Join(c.baseDirectory, "tmp")
	err = os.Mkdir(tmpCache, 0744)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	tmpFilePath := filepath.Join(tmpCache, fmt.Sprintf("%v", rand.Int()))
	defer os.Remove(tmpFilePath)
	err = os.WriteFile(tmpFilePath, bodyBytes, 0644)
	if err != nil {
		return "", err
	}
	err = os.Rename(tmpFilePath, cacheFileFullPath)
	if err != nil {
		defer os.Remove(cacheFileFullPath)
		return "", err
	}
	if cacheFileFullPath == "" {
		return "", errors.New("Unknown issue caching config file")
	}

	return cacheFileFullPath, nil
}

func cleanEtag(etag string) string {
	if strings.HasPrefix(etag, "W/") {
		// Weak etag, don't trust this https://www.rfc-editor.org/rfc/rfc7232#section-2.3
		return ""
	}
	if !strings.HasPrefix(etag, "\"") || !strings.HasSuffix(etag, "\"") {
		// Doesn't meet spec, don't trust this https://www.rfc-editor.org/rfc/rfc7232#section-2.3
		return ""
	}
	// Strip quotes
	coreEtag := etag[1 : len(etag)-1]
	fileCleanCoreEtag := filepath.Clean(coreEtag)
	if coreEtag != fileCleanCoreEtag {
		// This etag can't be used in filenames
		return ""
	}
	if splitPath, _ := filepath.Split(coreEtag); splitPath != "" {
		// This etag can't be used in filenames -- has path separator
		return ""
	}
	return coreEtag
}
