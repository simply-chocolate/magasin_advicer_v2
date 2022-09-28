package utils

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"
)

type AdviceCache struct {
	LastAdviceDocNum string `json:"lastAdviceDocNum"`
}

func path() string {
	basePath, _ := filepath.Abs("cache")

	return filepath.Join(basePath, filepath.Base("advice.json"))
}

// Read the advice number currently stored in the cache
func ReadAdviceCache() (AdviceCache, error) {
	file, err := ioutil.ReadFile(path())
	if err != nil {
		return AdviceCache{}, err
	}

	data := AdviceCache{}
	if err = json.Unmarshal([]byte(file), &data); err != nil {
		return AdviceCache{}, err
	}

	return data, nil
}

// Write advice number to the cache
func WriteAdviceCache(advice AdviceCache) error {
	adviceCache, err := ReadAdviceCache()
	if err != nil {
		return err
	}

	if err = mergo.Merge(&adviceCache, advice, mergo.WithOverride); err != nil {
		return err
	}

	data, err := json.Marshal(adviceCache)
	if err != nil {
		return err
	}

	if _, err = os.Create(path()); err != nil {
		return err
	}

	if err = ioutil.WriteFile(path(), data, fs.ModeAppend); err != nil {
		return err
	}

	return nil
}
