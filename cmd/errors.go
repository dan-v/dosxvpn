package cmd

import "errors"

var (
	errorMissingToken  = errors.New("need to have environment variable DIGITALOCEAN_ACCESS_TOKEN set")
	errorMissingRegion = errors.New("need to specify region")
	errorMissingName   = errors.New("need to specify name")
)
