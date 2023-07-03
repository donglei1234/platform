package mage

import (
	"gitlab.spacerouter.net/ca/dev/sr/ops/mage/pkg/mage"
)

// SrcPackage is the package to run mage commands against.
const SrcPackage = "github.com/donglei1234/platform/services/chat"

// srcPath is the file path of SrcPackage.
var srcPath = mage.PathOfGoPackage(SrcPackage)

func Test() error {
	return mage.Test(SrcPackage)
}

func TestCover() error {
	return mage.TestCover(SrcPackage)
}

func Build() error {
	return mage.Build(SrcPackage)
}

func Install() error {
	return mage.Install(srcPath)
}

func DockerBuild() error {
	return mage.DockerBuild(srcPath)
}

func DockerPush() error {
	return mage.DockerPush(srcPath)
}

func Package() error {
	return mage.Package(srcPath)
}

func Deploy() error {
	return mage.Deploy(srcPath)
}

func Delete() error {
	return mage.Delete(srcPath)
}
