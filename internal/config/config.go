package config

import "os"

const (
	// TODO: this needs to be tied to the engine path directory
	directory = "./.jf/config"
)

func createDirIfNotExist() error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
