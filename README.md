# lenv

> **Warning**
`lenv` is in alpha version. Bugs and instability are possible.

`lenv` is a simple tool to manage multiple Java/Python versions on a single machine. It is inspired by [pyenv](https://github.com/pyenv/pyenv) for Python.

```
$ java -version
openjdk version "1.8.0_432-432"
OpenJDK Runtime Environment (build 1.8.0_432-432-b06)
OpenJDK 64-Bit Server VM (build 25.432-b06, mixed mode)

$ lenv java install 11-openjdk
$ lenv java global 11-openjdk
$ java -version

openjdk version "18.0.2" 2022-07-19
OpenJDK Runtime Environment (build 18.0.2+9-61)
OpenJDK 64-Bit Server VM (build 18.0.2+9-61, mixed mode, sharing)
```

## Installation
### Windows
```
iwr -useb https://raw.githubusercontent.com/kiber-io/lenv/main/win_install.ps1 | iex
```
### Linux
```
bash <(curl -s https://raw.githubusercontent.com/kiber-io/lenv/main/linux_install.sh)
```

## Usage
### List all available versions
```
$ lenv list --all
Available Versions:
    23-openjdk
  * 18.0.2-openjdk
    11.0.2-openjdk
    11-openjdk
 -> 8-openjdk
```
`*` - downloaded localy

`->` - currently active version

### Install specific version
```
$ lenv install 11-openjdk
...
Java version 11-openjdk installed
```

### Set specific version as a global
```
$ lenv global 11-openjdk
Java version 11-openjdk set as global
```

## Uninstall
Simply remove the `.lenv` directory from your home directory.