package config

import (
    "path/filepath"
    "io/ioutil"
	"testing"
    "os"
)

func TestMarshalDefaultConfig(t *testing.T) {
    dir, err := ioutil.TempDir("", "history-test")
    if err != nil {
	    t.Error(err)
    }
    defer os.RemoveAll(dir)
    
    file := filepath.Join(dir, "config.toml")

    // Loading the config for the first time
    // creates the default and marshals it
    err = Conf.LoadFile(file)
    if err != nil {
        t.Errorf("Failed to create default config file: %v", err) 
    }

    // Loading the config for the second time
    // reads it from disk and unmarshals it
    err = Conf.LoadFile(file)
    if err != nil {
        t.Errorf("Failed to load default config file: %v", err) 
    }
}

