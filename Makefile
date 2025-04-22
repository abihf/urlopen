DESTDIR ?= /
PREFIX ?= /usr
BINDIR ?= $(DESTDIR)$(PREFIX)/bin
SHAREDDIR ?= $(DESTDIR)$(PREFIX)/share

build: dist/urlopen

dist/urlopen: main.go go.mod go.sum
	mkdir -p dist
	go build -o dist/urlopen main.go

install: dist/urlopen
	install -Dm755 dist/urlopen $(BINDIR)/urlopen
	install -Dm644 urlopen.desktop $(SHAREDDIR)/applications/urlopen.desktop