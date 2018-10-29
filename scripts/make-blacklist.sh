#!/bin/bash

function list_1() {
    curl 'https://tspprs.com/dl/{ads,malware,phishing,ransomware,scam,tracking}' | grep -v '\-\-\_curl\_\-\-'
}

function list_2() {
    curl 'http://pgl.yoyo.org/adservers/serverlist.php?mimetype=plaintext' | grep -v '#' | sed 's/.* //'
}

function lower() {
    tr '[:upper:]' '[:lower:]'
}

function trim() {
    sed 's/ *//' | sed 's/ *$//'
}

function filter() {
    sed '/[^a-zA-Z0-9\._\-]/d' | sed '/^$/d' | sed '/^-.*$/d' | sed '/^.*-$/d'
}

export LC_ALL='C'
cat <(list_1) <(list_2) \
    | lower \
    | trim \
    | filter \
    | sort \
    | uniq \
    > ../data/blacklist.txt
