# javaenv
`javaenv` is a simple tool to manage multiple Java versions on a single machine. It is inspired by [pyenv](https://github.com/pyenv/pyenv) for Python.

```
$ java -version
openjdk version "1.8.0_432-432"
OpenJDK Runtime Environment (build 1.8.0_432-432-b06)
OpenJDK 64-Bit Server VM (build 25.432-b06, mixed mode)

$ javaenv install 11-openjdk
$ javaenv global 11-openjdk

openjdk version "18.0.2" 2022-07-19
OpenJDK Runtime Environment (build 18.0.2+9-61)
OpenJDK 64-Bit Server VM (build 18.0.2+9-61, mixed mode, sharing)
```

## Installation
### Windows
```
iwr -useb https://raw.githubusercontent.com/kiber-io/javaenv/main/win_install.ps1 | iex
```

## Usage
### List all available Java versions
```
$ javaenv list --all
Available Versions:
    23-openjdk
  * 18.0.2-openjdk
    11.0.2-openjdk
    11-openjdk
 -> 8-openjdk
```
`*` - downloaded localy

`->` - currently active version

### Install a Java version
```
$ javaenv install 11-openjdk
...
Java version 11-openjdk installed
```

### Set a Java version as the global version
```
$ javaenv global 11-openjdk
Java version 11-openjdk set as global
```