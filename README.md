# Argo CD Checker

A tiny tool that verifies that the Application and ApplicationSet resources have a valid destination path, and that `kustomize build` runs successfully in the application and component folders.

For example:

```
$ check-argocd --base-dir=$(pwd) --apps=apps-of-apps,apps --components=components
INFO 🏁 Checking Argo CD configuration base-dir=/path/to/repository
INFO 👀 checking Applications and ApplicationSets path=apps-of-apps
INFO 👀 checking Applications and ApplicationSets path=apps
INFO 👀 checking Components path=components
```

## Building

Requires Go version 1.20.x (1.20.11 or higher) - download for your development environment [here](https://golang.org/dl).

To install, execute:

```
$ make install
```

This builds the `check-argocd` binary and copies it into `$GOPATH/bin`.


## License

The code is available under the terms of the [Apache License 2.0](LICENCE).