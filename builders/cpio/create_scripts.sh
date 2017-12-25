#!/bin/sh

if [ ! -d 'scripts' ]; then
    exit 0
fi

mkdir -p flat/base.pkg
( cd scripts && find . | cpio -o --format odc --owner 0:80 | gzip -c ) > flat/base.pkg/Scripts

ls -al flat/base.pkg/
