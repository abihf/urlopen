package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	yaml "gopkg.in/yaml.v2"
)

type configItem struct {
	Prefix  []string `yaml:"prefix"`
	Command string   `yaml:"command"`
}

type config map[string]configItem

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <url>\n", os.Args[0])
		os.Exit(1)
	}
	panic(open(os.Args[1], os.Args[2:]))
}

func open(url string, args []string) error {
	c, err := openConfig()
	if err != nil {
		return err
	}

	cmdStr, err := findCommand(url, c)
	if err != nil {
		return err
	}

	fullCmd, err := exec.LookPath(cmdStr)
	if err != nil {
		return err
	}

	fullArgs := append([]string{fullCmd, url}, args...)
	env := os.Environ()

	fmt.Printf("args %v", fullArgs)

	return syscall.Exec(fullCmd, fullArgs, env)
}

func findCommand(url string, c config) (string, error) {
	for _, item := range c {
		for _, p := range item.Prefix {
			fmt.Printf("Prefix: %s [%s]\n", url, p)
			if strings.HasPrefix(url, p) {
				return item.Command, nil
			}
		}
	}
	return "", fmt.Errorf("Can not found command for %s", url)
}

func openConfig() (config, error) {
	confFile := os.Getenv("HOME") + "/.config/urlopen.yml"
	file, err := os.Open(confFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var c config
	err = yaml.NewDecoder(file).Decode(&c)
	return c, err
}
