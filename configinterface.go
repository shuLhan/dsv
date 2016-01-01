package dsv

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

/*
ConfigInterface for reader and writer for initializing the config from JSON.
*/
type ConfigInterface interface {
	GetConfigPath() string
	SetConfigPath(dir string)
}

/*
ConfigOpen configuration file and initialize the attributes.
*/
func ConfigOpen(rw ConfigInterface, fcfg string) error {
	cfg, e := ioutil.ReadFile(fcfg)

	if nil != e {
		return e
	}

	// Get directory where the config reside.
	rw.SetConfigPath(path.Dir(fcfg))

	return ConfigParse(rw, cfg)
}

/*
ConfigParse from JSON string.
*/
func ConfigParse(rw ConfigInterface, cfg []byte) error {
	return json.Unmarshal([]byte(cfg), rw)
}

/*
ConfigCheckPath if no path in file, return the config path plus file name,
otherwise leave it unchanged.
*/
func ConfigCheckPath(comin ConfigInterface, file string) string {
	dir := path.Dir(file)

	if dir == "." {
		cfgPath := comin.GetConfigPath()
		if cfgPath != "" && cfgPath != "." {
			return cfgPath + "/" + file
		}
	}

	// nothing happen.
	return file
}
