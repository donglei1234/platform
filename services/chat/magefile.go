//go:build mage
// +build mage

// This is the CI/CD script for the "chat" package. Targets are configured
// through environment variables.
// `+"\n"+`
// - Run the "verify" target to ensure you have all command line tools installed
// on your PATH necessary to run all targets.
// `+"\n"+`
// - Run the "usage" target to see a list of configurable settings.
// `+"\n"+`
// - To see more info about a target, including relevant environment variables
// and their evaluation order, run:
// mage -h <target>
// `+"\n"+`
// - To run a target in verbose mode use:
// mage -v <target>
// `+"\n"+`
// - To run a target in debug mode use:
// mage -debug <target>
// `+"\n"+`
// - Multiple targets may be run sequentially, exiting on failure:
// mage <target1> <target2> <target3>
package main

import (
	chat "github.com/donglei1234/platform/services/chat/build/mage"
	"github.com/magefile/mage/mg"
	m "gitlab.spacerouter.net/ca/dev/sr/ops/mage/pkg/mage"
)

// Show settings that can be configured via environment variables.
func Usage() error {
	return m.Usage()
}

// Verify all necessary command line tools are installed.
// If all tools are available on the command line there will be no output.
// Otherwise, any missing tools will be displayed.
func Verify() error {
	return m.Verify()
}

// Run all tests.
// Accepts all Go tool environment variables.
func Test() error {
	return chat.Test()
}

// Run all tests and generate a coverage profile.
// Accepts all Go tool environment variables.
// Generates raw <TEST_COVER_PROFILE_FILE>, aggregate <TEST_COVER_AGGREGATE_FILE>,
// and HTML <TEST_COVER_HTML_FILE> coverage profiles within <TEST_ARTIFACT_DIR>.
func TestCover() error {
	return chat.TestCover()
}

// Compile all packages.
// Accepts all Go tool environment variables.
func Build() error {
	return chat.Build()
}

// Compile and install all packages.
// Accepts all Go tool environment variables.
// Collects all "install" entries from all <SPACEROUTER_METADATA_FILE> files
// and installs one binary per entry to <INSTALL_DIR | GOBIN | GOPATH/bin>.
func Install() error {
	return chat.Install()
}

type Docker mg.Namespace

// Log in to the Docker registry.
// Logs into the <DOCKER_REGISTRY | CI_REGISTRY> with username <DOCKER_USERNAME | CI_REGISTRY_USER>
// and password <DOCKER_PASSWORD | GITLAB_API_TOKEN | CI_BUILD_TOKEN>.
func (Docker) Login() error {
	return m.DockerLogin()
}

// Build all Docker images.
// Collects all "images" entries from all <SPACEROUTER_METADATA_FILE> files
// and builds one Docker image per entry. Images will be tagged with the current
// Git tag or <DOCKER_TAG> if set. If an entry's "mainPackage" binary does not
// exist it will be installed to <INSTALL_DIR | GOBIN | GOPATH/bin>.
func (Docker) Build() error {
	return chat.DockerBuild()
}

// Push all Docker images.
// Collects all "images" entries from all <SPACEROUTER_METADATA_FILE> files
// and pushes one Docker image per entry. Images will be tagged with the current
// Git tag or <DOCKER_TAG> if set.
func (Docker) Push() error {
	return chat.DockerPush()
}

// Package all Helm charts.
// Collects all "deployments" entries from all <SPACEROUTER_METADATA_FILE> files
// and recursively updates the Helm dependencies of each entry.
func Package() error {
	return chat.Package()
}

// Deploy all Helm charts.
// Collects all "deployments" entries from all <SPACEROUTER_METADATA_FILE> files
// and deploy each.
// The Docker image tag used will be the current Git tag or <DOCKER_TAG> if set.
// Helm will use the <DEPLOY_KUBECONFIG | KUBECONFIG> file for Kubernetes authentication.
// The deployment will target <DEPLOY_NAMESPACE | KUBE_NAMESPACE> and be configured
// with an environment name <DEPLOY_ENVIRONMENT_NAME | CI_ENVIRONMENT_NAME | DEPLOY_NAMESPACE | KUBE_NAMESPACE>,
// an environment subdomain <DEPLOY_ENVIRONMENT_SUBDOMAIN | DEPLOY_ENVIRONMENT_NAME | CI_ENVIRONMENT_NAME | DEPLOY_NAMESPACE | KUBE_NAMESPACE>,
// and a deployment name of <DEPLOY_DEPLOYMENT | DEPLOY_ENVIRONMENT_NAME | CI_ENVIRONMENT_NAME | DEPLOY_NAMESPACE | KUBE_NAMESPACE>.
func Deploy() error {
	return chat.Deploy()
}

// Delete all Helm releases.
// Collects all "deployments" entries from all <SPACEROUTER_METADATA_FILE> files
// and deletes any associated deployed releases.
// Helm will use the <DEPLOY_KUBECONFIG | KUBECONFIG> file for Kubernetes authentication.
// The operation will target <DEPLOY_NAMESPACE | KUBE_NAMESPACE>.
func Delete() error {
	return chat.Delete()
}
