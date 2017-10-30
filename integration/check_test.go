package integration_test

import (
	"os/exec"

	"bytes"
	"encoding/json"

	"github.com/charlievieth/fs/testdata"
	"github.com/concourse/atc"
	"github.com/concourse/atc/cessna"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Check Resource", func() {

	Context("Executing /opt/resource/check", func() {

		var session *gexec.Session
		var err error

		var stdoutBytes []byte
		var stderrBytes []byte

		var checkRequest cessna.CheckRequest

		BeforeEach(func() {
			checkRequest = cessna.CheckRequest{
				Source: atc.Source{
					"password": os.Getenv("TEST_RETRO_PASSWORD"),
					"id":       os.Getenv("TEST_RETRO_ID"),
				},
				Version: atc.Version{},
			}
		})

		JustBeforeEach(func() {
			cmd := exec.Command(checkBinPath)

			jsonRequest, jsonErr := json.Marshal(checkRequest)
			Expect(jsonErr).ToNot(HaveOccurred())

			cmd.Stdin = bytes.NewBuffer(jsonRequest)

			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())

			Eventually(session, "10s").Should(gexec.Exit(0))

			stdoutBytes = session.Out.Contents()
			stderrBytes = session.Err.Contents()
		})

		It("unmarshals correctly", func() {
			var checkResponse cessna.CheckResponse

			err := json.Unmarshal(stdoutBytes, &checkResponse)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
