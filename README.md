Build macOS packages using [Google Container Builder](https://cloud.google.com/container-builder/).

# Example

```
# create package root
mkdir -p build-dir/root

# put stuff in package root
mkdir -p build-dir/root/Users/Shared/foo
touch build-dir/root/Users/Shared/foo/bar

# submit cloudbuild job to create package
# the package will be uploaded to the 'gs://mac-packages' bucket. Make sure to replace with a bucket name you have write permissions to.
make package PACKAGE_NAME=foo-1.2.3.pkg PACKAGE_VERSION=1.2.3 PACKAGE_IDENTIFIER=co.acme.foo
```

# How it works

A macOS flat package is a xar archive with a specific structure(see References). This repo takes advantage of several linux utilities to build a GCP Container Builder pipeline which creates a new macOS package.
All the build steps assume that the package root is located inside `build-dir/root/` and if there are package scripts, they're in `build-dir/scripts/`.

The provided Makefile `package` target abstracts the `gcloud container builds submit` step, which can be seen in full below:

```
gcloud container builds submit ./build-dir/ \
	--config ./cloudbuild.yml \
	--substitutions=_PACKAGE_NAME=foo-1.2.3.pkg,_PACKAGE_IDENTIFIER=co.acme.foo,_PACKAGE_VERSION=1.2.3
```

The full pipeline can be seen in the `cloudbuild.yml` file at the root of the repo. The steps are ordered using the `id` and `waitFor` directives in each build step.

## Builders

The package pipeline is composed of several build steps, each of which is made up of a containerized linux utility. Every builder has an assocated `Dockerfile`, `cloudbuild.yml` and `make $builder` target.

```
builders
├── bomutils
│   ├── Dockerfile
│   ├── cloudbuild.yml
│   └── create_bom.sh
├── build-info
│   ├── Dockerfile
│   └── cloudbuild.yml
├── cpio
│   ├── Dockerfile
│   ├── cloudbuild.yml
│   ├── create_payload.sh
│   └── create_scripts.sh
└── xar
    ├── Dockerfile
    ├── cloudbuild.yml
    └── create_xar.sh
```

### build-info

The `build-info` utility is a Go script which takes a few CLI arguments and traverses the `root` and `scripts` folders to build a `PackageInfo` file required by the package archive.
The build-info utility is the least complete of all the steps, but could be updated to fit more complex requirements. 
Create an [issue](https://github.com/kolide/gcp-pkgbuild/issues/new) or [pull request](https://github.com/kolide/gcp-pkgbuild/issues/new).  

# TODO

- [ ] Support all the PackageInfo directives.
- [ ] Build distribution style packages. Right now only simple flat packages are supported.
- [ ] Sign packages using [Google KMS](https://cloud.google.com/kms/) secrets.
- [ ] Add utility for bumping the package version number. 
- [ ] Automat Munki import with build pipeline.

# References

* http://s.sudre.free.fr/Stuff/Ivanhoe/FLAT.html
* http://bomutils.dyndns.org/tutorial.html
* https://gist.github.com/bruienne/ec5205408b9e52bd5cfc
* https://groob.io/posts/osx_packages/
