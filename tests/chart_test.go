package tests

import (
	"path/filepath"
	"testing"

	"github.com/jenkins-x/helm-unit-tester/pkg"
)

func TestChartsWithDifferentValues(t *testing.T) {
	chart := "../jxboot-helmfile-resources"

	pkg.RunTests(t, chart, filepath.Join("test_data"))
}
