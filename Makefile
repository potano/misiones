# Copyright © 2022 Michael Thompson
# SPDX-License-Identifier: GPL-2.0-or-later

PREFIX ?= /usr
BINDIR ?= $(PREFIX)/bin
MAN1DIR ?= $(PREFIX)/share/man/man1
INFODIR ?= $(PREFIX)/share/info

executables = misiones
man_files = misiones.1
info_files = misiones.info

bin: $(addprefix bin/, $(executables))

bin/misiones: cmd/misiones.go great/* parser/* sexp/* vectordata/*
	go build -o ./bin/misiones ./cmd/...

man: $(addprefix doc/, $(man_files))

%.1: %_manpage.adoc
	a2x -f manpage $<

doc/misiones.info: doc/misiones.adoc
	asciidoc -b docbook -d book -a data-uri -o doc/misiones.xml doc/misiones.adoc
	docbook2x-texi doc/misiones.xml --encoding=UTF-8 --to-stdout >doc/misiones.texi
	makeinfo --no-split -o doc/misiones.info doc/misiones.texi

info: doc/misiones.info

install-bin: $(addprefix bin/, $(executables)) | $(DESTDIR)$(BINDIR)
	install -m 755 $(addprefix bin/, $(executables)) $(DESTDIR)$(BINDIR)

install-man: $(addprefix doc/, $(man_files)) | $(DESTDIR)$(MAN1DIR)
	install -m 644 $(addprefix doc/, $(man_files)) $(DESTDIR)$(MAN1DIR)

install-info: $(addprefix doc/, $(info_files)) | $(DESTDIR)$(INFODIR)
	install -m 644 $(addprefix doc/, $(info_files)) $(DESTDIR)$(INFODIR)

$(DESTDIR)$(BINDIR) $(DESTDIR)$(MAN1DIR) $(DESTDIR)$(INFODIR):
	install -d -m 755 $@

all: bin man info

install: install-bin install-man install-info

dist: bin man info
	go run makedist.go

clean:
	rm -rf dist bin/*
	rm doc/\*.{html,info,man,texi,xml}

test:
	go test ./...


.PHONY: bin man info all install-bin install-man install-info install dist clean test

