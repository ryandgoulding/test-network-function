package operator

import (
	expect "github.com/google/goexpect"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/redhat-nfvpe/test-network-function/internal/reel"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf"
	"github.com/redhat-nfvpe/test-network-function/pkg/tnf/interactive"
	"time"
)

const (
	defaultTimeoutSeconds = 10
)

// The default test timeout.
var defaultTimeout = time.Duration(defaultTimeoutSeconds) * time.Second

var context *interactive.Context
var err error

var _ = ginkgo.Describe("operator_test", func() {
	ginkgo.When("A local shell is spawned", func() {
		goExpectSpawner := interactive.NewGoExpectSpawner()
		var spawner interactive.Spawner = goExpectSpawner
		context, err = interactive.SpawnShell(&spawner, defaultTimeout, expect.Verbose(true))
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(context).ToNot(gomega.BeNil())
		gomega.Expect(context.GetExpecter()).ToNot(gomega.BeNil())
	})

	ginkgo.When("Operator is already installed", func() {
		ginkgo.It("Checks if the CSV is installed successfully", func() {
			csv := NewCsv(TConfig.Csv.Name, TConfig.Csv.Namespace, defaultTimeout)
			csv.ExpectStatus = TConfig.Csv.Status
			gomega.Expect(csv).ToNot(gomega.BeNil())
			test, err := tnf.NewTest(context.GetExpecter(), csv, []reel.Handler{csv}, context.GetErrorChannel())
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(test).ToNot(gomega.BeNil())
			testResult, err := test.Run()
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(testResult).To(gomega.Equal(tnf.SUCCESS))
		})
	})

})
