package featureflags

import (
	"io/ioutil"

	"github.com/lab46/monorepo/gopkg/log"
	yaml "gopkg.in/yaml.v2"
)

var fflags map[string]bool

// GetAllFeatureFlags return available feature flags
func GetAllFeatureFlags() map[string]bool {
	return fflags
}

// ReadFromYAMLFile load flags from yaml file
func ReadFromYAMLFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(content, fflags)
	if err != nil {
		return err
	}
	log.Debugf("[featureflags] loaded from %s", filename)
	return nil
}

// IsActive check if a feature is on or not
func IsActive(featurename string) bool {
	return fflags[featurename]
}
