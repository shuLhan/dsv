package dsv

import (
	"encoding/json"
	"io/ioutil"
	"path"
)

/*
CommonInterface for reader and writer for initializing the config from JSON.
*/
type CommonInterface interface {
	GetPath () string
	SetPath (dir string)
	Init () error
}

/*
CheckPath if is name only without path, prefix it with path of configuration.
*/
func CheckPath (comin CommonInterface, file string) (string) {
	dir := path.Dir (file)

	if dir == "." {
		if comin.GetPath () != "" && comin.GetPath () != "." {
			return comin.GetPath () +"/"+ file
		}
	}

	// nothing happen.
	return file
}

/*
Open configuration file.
*/
func Open (comin CommonInterface, fcfg string) error {
	cfg, e := ioutil.ReadFile (fcfg)

	if nil != e {
		return e
	}

	// Get directory where the config reside.
	comin.SetPath (path.Dir (fcfg))

	e = ParseConfig (comin, cfg)

	if nil != e {
		return e
	}

	e = comin.Init ()

	return e
}

/*
ParseConfig from JSON string.
*/
func ParseConfig (comin CommonInterface, cfg []byte)  (e error) {
	return json.Unmarshal ([]byte (cfg), comin)
}
