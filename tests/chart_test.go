package tests

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jenkins-x/helm-unit-tester/pkg"
	"github.com/jenkins-x/jx/v2/pkg/apis/jenkins.io/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestChartsWithDifferentValues(t *testing.T) {
	chart := "../jxboot-helmfile-resources"

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
