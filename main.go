package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"syscall"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <url>\n", os.Args[0])
		os.Exit(1)
	}
	panic(open(os.Args[1:]))
}

func open(urls []string) error {
	c, err := loadConfig()
	if err != nil {
		return err
	}

	for _, rawURL := range urls {
		err = c.Open(rawURL)
		if err != nil {
			return err
		}
	}

	return nil
}

type Config struct {
	Default  string              `yaml:"default"`
	Browsers map[string]*Browser `yaml:"browsers"`

	Routes []*Route `yaml:"routes"`
}

func (c *Config) Open(rawURL string) error {
	parsedUrl, err := url.Parse(rawURL)
	if err != nil {
		fmt.Printf("can not parse url %s: %v\n", rawURL, err)
	} else {
		for _, route := range c.Routes {
			if route.Match(parsedUrl) {
				return c.Browsers[route.BrowserName].Open(rawURL)
			}
		}
	}

	return c.Browsers[c.Default].Open(rawURL)
}


type DirFunc func() string

var configDirs = []DirFunc{
	func() string {
		return os.Getenv("HOME") + "/.config"
	},
}

func loadConfig() (*Config, error) {
	for _, fc := range configDirs {
		conf, err := readConfig(fc() + "/urlopen.yml")
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			fmt.Printf("can not open config %v\n", err)
		}
		return conf, nil
	}
	return nil, fmt.Errorf("can not find config file")
}

func readConfig(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var c Config
	err = yaml.NewDecoder(file).Decode(&c)
	return &c, err
}



type Browser struct {
	Command string   `yaml:"command"`
	Args    []string `yaml:"args"`
}

func (b *Browser) Open(u string) error {
	bin, err := exec.LookPath(b.Command)
	if err != nil {
		return err
	}

	fullArgs := append([]string{bin}, b.Args...)
	fullArgs = append(fullArgs, u)

	return syscall.Exec(bin, fullArgs, os.Environ())
}

type Route struct {
	BrowserName string `yaml:"browser"`

	Domain       *string `yaml:"domain"`
	DomainSuffix *string `yaml:"domainSuffix"`
	PathPrefix   *string `yaml:"pathPrefix"`
	Scheme       *string `yaml:"scheme"`
}


func (r *Route) Match(u *url.URL) bool {
	return (r.Domain == nil || u.Host == *r.Domain) &&
		(r.DomainSuffix == nil || strings.HasSuffix(u.Host, *r.DomainSuffix)) &&
		(r.PathPrefix == nil || strings.HasPrefix(u.Path, *r.PathPrefix)) &&
		(r.Scheme == nil || u.Scheme == *r.Scheme)
}
