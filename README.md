# jxboot-helmfile-resources

Helm chart for resources used by [Jenkins X v3](https://jenkins-x.io/v3/)

See chart [README](./jxboot-resources/README.md) for install and config options.

## Building

To build this chart run:

``` 
make
```

to run the unit tests:

```
make test

```      

### Renerating test data

The following command will regenerate the test files in `tests/test_data/*/expected/**/*.yaml` if the underlying charts change:

```
make test-regen

```      
