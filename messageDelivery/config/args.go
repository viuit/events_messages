package config

import (
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	ConfigFile string `short:"c" long:"configFile" description:"The path to the config file"`
}

// New returns a new Config struct
func NewArgs() *Options {
	//Load args
	// var opts Options
	var opts Options
	p := flags.NewParser(&opts, flags.Default&^flags.HelpFlag)
	_, err := p.Parse()
	if err != nil {
		log.Fatalf("fail to parse args: %v", err)
	}
	// log.Printf("Args: %v\n", opts)

	return &opts
}
