package tests

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"io/ioutil"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v2"
)

const (
	resourcesSeparator = "---"

	// DefaultDirWritePermissions for directories
	DefaultDirWritePermissions = 0760

	// DefaultFileWritePermissions default permissions when creating a file
	DefaultFileWritePermissions = 0644
)

func TestDefaultCharts(t *testing.T) {
	chart := "../jxboot-helmfile-resources"

	outDir, err := ioutil.TempDir("", "helm-test-")
	require.NoError(t, err, "could not create temp dir")
	t.Logf("writing generated helm templates to %s", outDir)

	testDir := "test_data"
	files, err := ioutil.ReadDir(testDir)
	require.NoError(t, err, "could not read dir %s", testDir)
	for _, f := range files {
		if f.IsDir() {
			name := f.Name()
			valuesDir := filepath.Join(testDir, name, "values")
			expectedDir := filepath.Join(testDir, name, "expected")

			testOutDir := filepath.Join(outDir, name)
			resultsDir, _, err := assertHelmTemplate(t, chart, testOutDir, valuesDir)
			require.NoError(t, err, "failed to generate helm templates")
			require.NotEmpty(t, resultsDir, "no resultsDir returned")

			assertYamlExpected(t, expectedDir, resultsDir, name)
		}
	}
}

// assertYamlExpected asserts that the expectedDir of generated YAML is contained in the resultsDir
func assertYamlExpected(t *testing.T, expectedDir string, resultsDir string, testName string) {
	err := filepath.Walk(expectedDir,
	    func(path string, info os.FileInfo, err error) error {
	    if err != nil {
	        return err
	    }
	    if !info.IsDir() && filepath.Ext(path) == ".yaml" {
	    	t.Logf("testing expected file %s", path)

	    	relPath, err := filepath.Rel(expectedDir, path)
	    	if err != nil {
	    	  return errors.Wrapf(err, "failed to get base path for %s", path)
	    	}

	    	actualFile := filepath.Join(resultsDir, relPath)
			assertFilesSame(t, path, actualFile, testName)
		}
	    return nil
	})
	require.NoError(t, err, "failed to verify expected files")
}

func assertFilesSame(t *testing.T, expectedFile, actualFile, testName string) {
	if assert.FileExists(t, expectedFile, testName) && assert.FileExists(t, actualFile, testName) {
		expectedData,err := ioutil.ReadFile(expectedFile)
		require.NoError(t, err, "failed to load file %s", expectedFile)

		actualData,err := ioutil.ReadFile(actualFile)
		require.NoError(t, err, "failed to load file %s", actualFile)


		diff := cmp.Diff(string(actualData), string(expectedData))
		if diff != "" {
			t.Logf("generated: %s does not match expected: %s", actualFile, expectedFile)
			t.Logf("%s\n", diff)
			assert.Fail(t, "file %s is not the same as file %s", actualFile, expectedFile)
		}
	}
}

func assertHelmTemplate(t *testing.T, chart string, outDir, valuesDir string) (string, []string, error) {
	fileNames := []string{}

	args := []string{"template", chart, "--output-dir", outDir}

	files, err := ioutil.ReadDir(valuesDir)
	require.NoError(t, err, "could not read dir %s", valuesDir)
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, ".yaml") {
			args = append(args, "--values", filepath.Join(valuesDir, name))
		}
	}

	commandLine := "helm " + strings.Join(args, " ")
	t.Logf("invoking: %s\n", commandLine)

	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	assert.NoError(t, err, "failed to run: %s", commandLine)

	// TODO now lets assert that the files are valid
	templatesDir, err := filepath.Glob(filepath.Join(outDir, "*", "templates"))
	require.NoError(t, err, "could not find templates dir in dir %s", outDir)
	require.NotEmpty(t, templatesDir, "no */templates dir in dir %s", outDir)
	templateDir := templatesDir[0]

	resultDir := filepath.Join(outDir, "results")


	files, err = ioutil.ReadDir(templateDir)
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, ".yaml") {
			path := filepath.Join(templateDir, name)

			names, err := splitObjectsInFiles(t, path, resultDir)
			require.NoError(t, err, "failed to split the yaml file into resulting files for %s", path)
			fileNames = append(fileNames, names...)
		}
	}
	return resultDir, fileNames, err
}

func splitObjectsInFiles(t *testing.T, inputFile string, resultDir string) ([]string, error) {
	fileNames := make([]string, 0)
	f, err := os.Open(inputFile)
	if err != nil {
		return fileNames, errors.Wrapf(err, "opening inputFile %q", inputFile)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var buf bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()
		if line == resourcesSeparator {
			// ensure that we actually have YAML in the buffer
			data := buf.Bytes()
			if isWhitespaceOrComments(data) {
				buf.Reset()
				continue
			}

			fileName, err := writeBufferToFile(t, data, resultDir)
			require.NoError(t, err, "failed to write buffer to file")
			if fileName != "" {
				fileNames = append(fileNames, fileName)
			}
			buf.Reset()
		} else {
			_, err := buf.WriteString(line)
			if err != nil {
				return fileNames, errors.Wrapf(err, "writing line from inputFile %q into a buffer", inputFile)
			}
			_, err = buf.WriteString("\n")
			if err != nil {
				return fileNames, errors.Wrapf(err, "writing a new line in the buffer")
			}
		}
	}
	if buf.Len() > 0 && !isWhitespaceOrComments(buf.Bytes()) {
		data := buf.Bytes()
		fileName, err := writeBufferToFile(t, data, resultDir)
		require.NoError(t, err, "failed to write buffer to file")
		if fileName != "" {
			fileNames = append(fileNames, fileName)
		}
	}
	return fileNames, nil
}

func writeBufferToFile(t *testing.T, data []byte, resultDir string) (string, error) {
	m := yaml.MapSlice{}
	err := yaml.Unmarshal(data, &m)
	require.NoError(t, err, "failed to unmarshal YAML: %s", string(data))

		if len(m) == 0 {
			return "", nil
		}

	name := getYamlValueString(&m, "metadata", "name")
	kind := getYamlValueString(&m, "kind")
	require.NotEmpty(t, name, "resource with missing name: %s", string(data))
	require.NotEmpty(t, kind, "resource with missing kind: %s", string(data))

	outDir := filepath.Join(resultDir, kind)
	err = os.MkdirAll(outDir, DefaultDirWritePermissions)
	require.NoError(t, err, "failed to make output dir %s", outDir)

	fileName := filepath.Join(outDir, name+".yaml")

	err = ioutil.WriteFile(fileName, data, DefaultFileWritePermissions)
	require.NoError(t, err, "creating file %q", fileName)
	return fileName, err
}

func getYamlValueString(mapSlice *yaml.MapSlice, keys ...string) string {
	value := getYamlValue(mapSlice, keys...)
	answer, ok := value.(string)
	if ok {

		return answer
	}
	return ""
}

func getYamlValue(mapSlice *yaml.MapSlice, keys ...string) interface{} {
	if mapSlice == nil {
		return nil
	}
	if mapSlice == nil {
		return fmt.Errorf("No map input!")
	}
	m := mapSlice
	lastIdx := len(keys) - 1
	for idx, k := range keys {
		last := idx >= lastIdx
		found := false
		for _, mi := range *m {
			if mi.Key == k {
				found = true
				if last {
					return mi.Value
				} else {
					value := mi.Value
					if value == nil {
						return nil
					} else {
						v, ok := value.(yaml.MapSlice)
						if ok {
							m = &v
						} else {
							v2, ok := value.(*yaml.MapSlice)
							if ok {
								m = v2
							} else {
								return nil
							}
						}
					}
				}
			}
		}
		if !found {
			return nil
		}
	}
	return nil
}

// isWhitespaceOrComments returns true if the data is empty, whitespace or comments only
func isWhitespaceOrComments(data []byte) bool {
	if len(data) == 0 {
		return true
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		t := strings.TrimSpace(line)
		if t != "" && !strings.HasPrefix(t, "#") {
			return false
		}
	}
	return true
}
