package operator

import (
	"flag"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	. "github.com/redhat-nfvpe/test-network-function/test-network-function/operator-test/config"
	ginkgoreporters "kubevirt.io/qe-tools/pkg/ginkgo-reporters"
	"path"
	"testing"
)

const (
	defaultCliArgValue            = ""
	CnfCertificationTestSuiteName = "CNF Certification Operator Test Suite"
	junitFlagKey                  = "junit"
	JunitXMLFileName              = "cnf-operator-certification-tests_junit.xml"
	reportFlagKey                 = "report"
)

var junitPath *string
var reportPath *string
var TConfig *TestConfig

func init() {
	junitPath = flag.String(junitFlagKey, defaultCliArgValue,
		"the path for the junit format report")
	reportPath = flag.String(reportFlagKey, defaultCliArgValue,
		"the path of the report file containing details for failed tests")
	TConfig, _ = GetConfig()
}
//TestOperatorTest ...Main test operator function
func TestOperatorTest(t *testing.T) {
	RegisterFailHandler(Fail)
	var ginkgoReporters []Reporter
	if ginkgoreporters.Polarion.Run {
		ginkgoReporters = append(ginkgoReporters, &ginkgoreporters.Polarion)
	}
	if *junitPath != "" {
		junitFile := path.Join(*junitPath, JunitXMLFileName)
		ginkgoReporters = append(ginkgoReporters, reporters.NewJUnitReporter(junitFile))
	}
	RunSpecsWithDefaultAndCustomReporters(t, CnfCertificationTestSuiteName, ginkgoReporters)
}
