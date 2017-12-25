#!/bin/sh

( cd flat/base.pkg && xar --compression none -cf "../../${1}" * )
