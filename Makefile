all: package

package:
	gcloud container builds submit ./build-dir/ \
		--config ./cloudbuild.yml \
		--substitutions=_PACKAGE_NAME=$(PACKAGE_NAME),_PACKAGE_IDENTIFIER=$(PACKAGE_IDENTIFIER),_PACKAGE_VERSION=$(PACKAGE_VERSION)

builders: xar bomutils cpio build-info

build-info:
	gcloud container builds submit . --config ./builders/build-info/cloudbuild.yml

xar:
	gcloud container builds submit ./builders/xar/ --config ./builders/xar/cloudbuild.yml

bomutils:
	gcloud container builds submit ./builders/bomutils/ --config ./builders/bomutils/cloudbuild.yml

cpio:
	gcloud container builds submit ./builders/cpio/ --config ./builders/cpio/cloudbuild.yml


.PHONY: xar bomutils cpio build-info package
