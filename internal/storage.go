package internal

import (
	"os"
	"log"
	"encoding/json"
)

type ShortcutStore struct {
	Path string
}

type ShortcutData struct {
	Shortcuts Shortcuts `json:"shortcuts"`
}

func (ss *ShortcutStore) Init() error {
	defaultSS := NewDefaultShortcuts()
	if _, err := os.Stat(ss.Path); os.IsNotExist(err) {
		log.Printf("file doesn't exist %s: %v", ss.Path, err)
		if err := ss.SaveShortcuts(defaultSS); err != nil {
			return err
		}
		return nil
	}


	if _, err := ss.LoadShortcuts(); err != nil {
		return err
	}

	return nil
}

func (ss *ShortcutStore) SaveShortcuts(shorts Shortcuts) error {
	// remove help shortcut

	file, err := os.Create(ss.Path)
	if err != nil {
		return err
	}

	defer file.Close()

	temp := Shortcuts{}

	for k, v := range shorts {
		if k != "help" {
			temp[k] = v
		}
	}

	sc := ShortcutData{ Shortcuts: temp }

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(&sc); err != nil {
		return  err
	}

	return nil

}

func (ss *ShortcutStore) LoadShortcuts() (Shortcuts, error) {
	var sc ShortcutData
	file, err := os.Open(ss.Path)
	if err != nil {
		if !os.IsExist(err) {
			return NewDefaultShortcuts(), nil
		}
		return sc.Shortcuts, err
	}

	defer file.Close()
	finfo, err := file.Stat()
	if err != nil {
		return sc.Shortcuts, err
	}

	if finfo.Size() == 0 {
		return NewDefaultShortcuts(), nil
	}


	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&sc); err != nil {
		return sc.Shortcuts, err
	}

	if sc.Shortcuts == nil {
		return NewDefaultShortcuts(), nil
	}

	sc.Shortcuts["help"] = "/"

	if v, ok := sc.Shortcuts["*"]; !ok || v == "" {
		sc.Shortcuts["*"] = DefaultSearchProvider
	}

	return sc.Shortcuts, nil
}
