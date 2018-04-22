package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/lab46/monorepo/tools/envoy-confgen/confgen"

	yaml "gopkg.in/yaml.v2"
)

const (
	generateCommand = "generate"
	replaceCommand  = "replace"
)

func main() {
	var (
		confPath      string
		generatedFile string
		toFile        string
		fromFile      string
		command       string
		args          []string
		err           error
	)
	flag.StringVar(&confPath, "gen-file", "", "locate envoy config generator file")
	flag.StringVar(&generatedFile, "conf-name", "envoy_config.json", "name for generated config file")
	flag.StringVar(&toFile, "to", "", "replace this config file")
	flag.StringVar(&fromFile, "from", "", "replace from config file")
	flag.Parse()
	args = flag.Args()

	if len(args) >= 1 {
		command = args[0]
	}

	switch command {
	case generateCommand:
		err = generate(confPath, generatedFile)
	case replaceCommand:
		err = replace(fromFile, toFile)
	default:
		// print information
		log.Println(`
To generate a configuration, use 'generate command'. For example:
	envoy-confgen -gen-file=envoyconf.yaml generate

To replace old confifg with newer config, use replace command. For example:
	envoy-confgen -from=envoy_config.new.json -to=envoy.json replace
			`)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func generate(configPath, fileName string) error {
	conf := confgen.Generator{}
	yamlContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlContent, &conf)
	if err != nil {
		return err
	}

	return confgen.GenerateToFile(conf, fileName)
}

func replace(from, to string) error {
	return nil
}
