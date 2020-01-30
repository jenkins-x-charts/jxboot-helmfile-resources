# jxboot-resources

Helm chart for resources used by [jx boot](https://jenkins-x.io/getting-started/boot/)

See chart [readme](./jxboot-resources/README.md) for install and config options.


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
