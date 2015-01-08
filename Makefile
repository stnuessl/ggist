#
# Copyright (C) 2014  Steffen Nüssle
# ggist - go gist
#
# This file is part of ggist.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
# 
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.
#

#
# Simple Makefile for ggist to wrap the go toolchain into 
# a more familiar format
#

BIN =	ggist
SRC =	src/ggist.go 			\
	src/gist/gistapi.go 		\
	src/gist/history.go 		\
	src/util/print.go   		\
	src/util/cmdparser.go
	
INSTALL_DIR ?=	/usr/local/bin/

${BIN}: ${SRC}
	GOPATH=`pwd` go build $<

clean:
	rm -rf ${BIN}
install: ${BIN}
	cp ${BIN} ${INSTALL_DIR}

uninstall:
	rm -rf ${INSTALL_DIR}${BIN}

.PHONY: clean uninstall
.SILENT: clean