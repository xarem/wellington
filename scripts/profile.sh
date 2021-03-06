#!/bin/bash

export rmnpath=$RMN_BASE_PATH
export guipath=$rmnpath/www/gui

FILES=$(find $rmnpath/www/gui/sass -name "[^_]*\.scss")

time wt compile --cpuprofile=wt.prof --multi --gen $guipath/build/im --font $guipath/font-face -b $guipath/build/css/ -p $guipath/sass --images-dir $guipath/im/sass $FILES

go tool pprof --png $GOPATH/bin/wt wt.prof > profile.png
open profile.png
