# dpbuild

dpbuild is a CLI tool for building Debian packages using pbuilder as the backend.

## Features

- Build Debian packages in a clean chroot environment
- Support multiple Ubuntu/Debian distributions and architectures
- Manage build environments through YAML configuration
- Create source packages with orig.tar.xz and dsc files
- Build packages with consistent directory structure output

## Prerequisites

- pbuilder
- dpkg-dev
- `sudo` privileges for running pbuilder commands

### Debian

`ubuntu-keyring` is required for building Ubuntu packages. Note that it is only available in Debian GNU/Linux sid.

```bash
# apt install pbuilder dpkg-dev ubuntu-keyring debian-archive-keyring sudo
```

### Ubuntu

`debian-keyring` is required for building Debian packages.

```bash
# apt install pbuilder dpkg-dev debian-keyring debian-archive-keyring sudo
```

## Installation

### Using go install

```bash
$ go install github.com/zinrai/dpbuild/cmd/dpbuild@latest
```

### Build from source

```bash
$ go build -o dpbuild cmd/dpbuild/main.go
```

## Configuration

Create `~/.dpbuild/config.yaml`. Below is an example configuration for building packages for Ubuntu (jammy, noble) and Debian (bookworm):

```yaml
environments:
  - distribution: jammy
    architecture: amd64
    mirror: http://jp.archive.ubuntu.com/ubuntu
    components: [main, universe]
  - distribution: noble
    architecture: amd64
    mirror: http://jp.archive.ubuntu.com/ubuntu
    components: [main, universe]
  - distribution: bookworm
    architecture: amd64
    mirror: http://httpredir.debian.org/debian
    components: [main]
```

## Usage

### Initialize build environments

Create pbuilder base tarballs for all environments in config.yaml:

```bash
$ dpbuild init
```

### Update build environments

Update pbuilder base tarballs for all environments in config.yaml:

```bash
$ dpbuild update
```

### Create source package

Run in a directory containing the `debian` directory:

```bash
$ dpbuild source
```

This command:
- Creates `*.orig.tar.xz` in the parent directory
- Generates a `*.dsc` file

### Build package

Build a Debian package using pbuilder:

```bash
$ dpbuild pkg --dist bookworm --arch amd64 --dsc ../package_1.0-1.dsc
```

The built packages will be stored in `packages/${distribution}-${architecture}/`.

## License

This project is licensed under the MIT License - see the [LICENSE](https://opensource.org/license/mit) for details.
