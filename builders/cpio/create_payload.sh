#!/bin/sh

mkdir -p flat/base.pkg
( cd root && find . | cpio -o --format odc --owner 0:80 | gzip -c ) > flat/base.pkg/Payload

ls -al flat/base.pkg
