package tests

import (
	"github.com/jenkins-x/jx-api/v4/pkg/util"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jenkins-x/helm-unit-tester/pkg"
	"github.com/jenkins-x/jx-api/v4/pkg/apis/jenkins.io/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestChartsWithDifferentValues(t *testing.T) {
	chart := filepath.Join("..", "charts", "jxboot-helmfile-resources")

	_, testcases := pkg.RunTests(t, chart, filepath.Join("test_data"))

	envs := []string{"dev", "production", "staging"}

	for _, tc := range testcases {
		remoteCluster := false
		expectedEnvironmentScheduler := "in-repo"
		expectedDefaultScheduler := "jx-meta-pipeline"

		if tc.Name == "lighthouse-jx" {
			expectedEnvironmentScheduler = "environment"
			expectedDefaultScheduler = "default"
		}

		switch tc.Name {
		case "remote-env":
			remoteCluster = true

		case "custom-env", "no-envs":
			continue

		case "gke-domain":
			assertChartRepoIngress(t, tc, true, false)

		case "gke-domain-bucketrepo":
			assertChartRepoIngress(t, tc, false, true)

		case "gke-domain-no-repo":
			assertChartRepoIngress(t, tc, true, false)

		case "gke-domain-none-repo":
			assertChartRepoIngress(t, tc, true, false)
		}

		dir := filepath.Join(tc.OutDir, "results", "jenkins.io", "v1")
		for _, e := range envs {
			file := filepath.Join(dir, "Environment", e+".yaml")
			assert.FileExists(t, file)
			data, err := ioutil.ReadFile(file)
			require.NoError(t, err, "failed to load file %s", file)
			env := &v1.Environment{}
			err = yaml.Unmarshal(data, env)
			require.NoError(t, err, "failed to parse file %s", file)

			validationErrors, err := util.ValidateYaml(env, data)
			require.NoError(t, err, "failed to validate %s for test %s", file, tc.Name)

			for _, ve := range validationErrors {
				t.Logf("test %s file %s validation error: %s\n", tc.Name, file, ve)
			}
			assert.Emptyf(t, validationErrors, "validation errors for file %s for test %s", file, tc.Name)

			if env.Name == "dev" {
				assert.Equal(t, expectedDefaultScheduler, env.Spec.TeamSettings.DefaultScheduler.Name, "env.Spec.TeamSettings.DefaultScheduler.Name: %s", env.Name)
			}

			assert.Equal(t, remoteCluster, env.Spec.RemoteCluster, "env.Spec.RemoteCluster for environment %s", env)
		}

		for _, e := range envs {
			file := filepath.Join(dir, "SourceRepository", e+".yaml")
			assert.FileExists(t, file)
			data, err := ioutil.ReadFile(file)
			require.NoError(t, err, "failed to load file %s", file)
			sr := &v1.SourceRepository{}
			err = yaml.Unmarshal(data, sr)
			require.NoError(t, err, "failed to parse file %s", file)

			assert.Equal(t, expectedEnvironmentScheduler, sr.Spec.Scheduler.Name, "sr.Spec.Scheduler.Name for environment: %s", sr)
		}
	}
}

func assertChartRepoIngress(t *testing.T, tc *pkg.TestCase, expectChartMuseum bool, expectBucketRepo bool) {
	dir := filepath.Join(tc.OutDir, "results", "networking.k8s.io", "v1beta1", "Ingress")

	assertFileExists(t, expectChartMuseum, filepath.Join(dir, "chartmuseum.yaml"), tc.Name)
	assertFileExists(t, expectBucketRepo, filepath.Join(dir, "bucketrepo.yaml"), tc.Name)
}

func assertFileExists(t *testing.T, exists bool, path, name string) {
	if exists {
		assert.FileExists(t, path)
		t.Logf("expected file is created %s for test %s\n", path, name)
	} else {
		assert.NoFileExists(t, path)
		t.Logf("no file is exist %s for test %s\n", path, name)
	}
}
