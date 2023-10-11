package utils

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/imdario/mergo"
)

type AdviceCache struct {
	LastAdviceDocNum string `json:"lastAdviceDocNum"`
}

func path(fileName string) string {
	basePath, _ := filepath.Abs("cache")

	return filepath.Join(basePath, filepath.Base(fmt.Sprintf("%s.json", fileName)))
}

// Read the advice number currently stored in the cache
func ReadAdviceCache(fileName string) (AdviceCache, error) {
	file, err := os.ReadFile(path(fileName))
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
func WriteAdviceCache(advice AdviceCache, filePath string) error {
	adviceCache, err := ReadAdviceCache(filePath)
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

	if _, err = os.Create(path(filePath)); err != nil {
		return err
	}

	if err = os.WriteFile(path(filePath), data, fs.ModeAppend); err != nil {
		return err
	}

	return nil
}
