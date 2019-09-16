#!/bin/sh
executable="$1"
string="$2"
# Quoting and robustness can be improved, but it works for the purpose I used it for
for segment in $(readelf --segments $executable | grep "\." | grep -v Requesting | grep -v Segment | cut -b6- | xargs echo | tr ' ' '\n' | sort | uniq); do echo $segment; readelf -p $segment $1 | grep $string; done
