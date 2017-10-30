package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"

	"testing"
)

var (
	checkBinPath string
	err          error
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)

	checkBinPath, err = gexec.Build("github.com/vlad-stoian/postfacto-concourse-resource/cmd/check")
	Expect(err).NotTo(HaveOccurred())

	SynchronizedAfterSuite(func() {
	}, func() {
		gexec.CleanupBuildArtifacts()
	})

	RunSpecs(t, "Integration Suite")
}
