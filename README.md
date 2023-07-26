# govanity

A simple server for serving Go vanity URLs.

## Motivation

There are quite a lot of similar projects already. However most of them had one or another drawback. The most common issue was that each and every subpackage had to be mapped manually.

govanity is a simple implementation that matches pagages against known **prefixes** allowing arbitrary subpackages to work instantly.

## Usage

You can run `govanity` like this:

```shell
govanity server [mappings...]
```

Each of the `mappings` is a string of the format `prefix=[vcs:]repo-root` where `prefix` is the import name of a package, `repo-root` is the location of the repository containing the code and `vcs` is the version control system. If `vcs` is not specified it defaults to `git`.

An import prefix is matched against its host and path, so for a mapping like `codello.dev/govanity=https://github.com/codello/govanity` `govanity` would respond at `GET go.codello.dev/govanity`. Alternatively you can omit the host of the import prefix. Then `govanity` will use the host of the request in its place.