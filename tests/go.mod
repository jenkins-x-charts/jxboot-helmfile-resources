module github.com/jenkins-x/jxboot-helmfile-resources

go 1.12

require (
	github.com/gophercloud/gophercloud v0.0.0-20190126172459-c818fa66e4c8 // indirect
	github.com/jenkins-x/helm-unit-tester v0.0.4
	github.com/jenkins-x/jx v0.0.0-20200129202546-993ff917ca15
	github.com/stretchr/testify v1.4.0
	google.golang.org/appengine v1.5.0 // indirect
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/yaml v1.1.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190528110122-9ad12a4af326
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190221084156-01f179d85dbc
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190528110200-4f3abb12cae2
)
