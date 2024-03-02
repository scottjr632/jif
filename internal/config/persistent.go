package config

import (
	"encoding/gob"
	"os"
)

const (
	version          = "0.0.1"
	configBinaryName = "workspace"
)

var fullPath = directory + "/" + configBinaryName + "_" + version

type Persistent struct {
	BranchPrefix string
}

func (p *Persistent) SetBranchPrefix(prefix string) error {
	p.BranchPrefix = prefix
	return p.Save()
}

func (p *Persistent) GetBranchPrefix() string {
	return p.BranchPrefix
}

func (p *Persistent) Save() error {
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	defer file.Close()
	enc := gob.NewEncoder(file)
	return enc.Encode(*p)
}

func LoadPersistenConfig() (*Persistent, error) {
	if err := createDirIfNotExist(); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fileInfo.Size() == 0 {
		return &Persistent{}, nil
	}
	dec := gob.NewDecoder(file)
	var config Persistent
	if err := dec.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
