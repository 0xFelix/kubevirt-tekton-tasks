package parse_test

import (
	"fmt"
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/create-datavolume-from-manifest/pkg/utils/parse"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/sharedtest/testobjects/datavolume"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
)

const (
	testStringDVNamespace1   = "dv-namespace-test-1"
	testStringDVNamespace2   = "dv-namespace-test-2"
	testStringWaitForSuccess = "true"
)

var (
	testDVManifest1 = strings.TrimSpace(datavolume.NewBlankDataVolume("testDv1").ToString())
	testDVManifest2 = strings.TrimSpace(datavolume.NewBlankDataVolume("testDv2").WithNamespace(testStringDVNamespace2).ToString())
)

var _ = Describe("CLIOptions", func() {
	Context("invalid cli options", func() {
		DescribeTable("Init return correct assertion errors", func(expectedErrMessage string, options *parse.CLIOptions) {
			err := options.Init()
			Expect(err).Should(HaveOccurred())
			fmt.Println(err.Error())
			Expect(err.Error()).To(ContainSubstring(expectedErrMessage))
		},
			Entry("no dv-manifest", "dv-manifest param has to be specified", &parse.CLIOptions{}),
			Entry("wrong output type", "non-existing is not a valid output type",
				&parse.CLIOptions{
					DataVolumeManifest: testDVManifest1,
					Output:             "non-existing",
				}),
		)
	})
	Context("correct cli options", func() {
		DescribeTable("Init should succeed", func(options *parse.CLIOptions) {
			Expect(options.Init()).To(Succeed())
		},
			Entry("with yaml output", &parse.CLIOptions{
				DataVolumeManifest:  testDVManifest1,
				DataVolumeNamespace: testStringDVNamespace1,
				Output:              "yaml",
			}),
			Entry("with json output", &parse.CLIOptions{
				DataVolumeManifest:  testDVManifest1,
				DataVolumeNamespace: testStringDVNamespace1,
				Output:              "json",
			}),
			Entry("with debug loglevel", &parse.CLIOptions{
				DataVolumeManifest:  testDVManifest1,
				DataVolumeNamespace: testStringDVNamespace1,
				Debug:               true,
			}),
			Entry("with WaitForSuccess", &parse.CLIOptions{
				DataVolumeManifest:  testDVManifest1,
				DataVolumeNamespace: testStringDVNamespace1,
				WaitForSuccess:      testStringWaitForSuccess,
			}),
		)

		It("Init should trim spaces", func() {
			options := &parse.CLIOptions{
				DataVolumeManifest:  " " + testDVManifest1 + " ",
				DataVolumeNamespace: " " + testStringDVNamespace1 + " ",
				WaitForSuccess:      " " + testStringWaitForSuccess + " ",
			}
			Expect(options.Init()).To(Succeed())
			Expect(options.DataVolumeManifest).To(Equal(testDVManifest1), "DataVolumeManifest should equal")
			Expect(options.DataVolumeNamespace).To(Equal(testStringDVNamespace1), "DataVolumeNamespace should equal")
			Expect(options.WaitForSuccess).To(Equal(testStringWaitForSuccess), "WaitForSuccess should equal")
		})

		DescribeTable("CLI options should return correct values", func(fnToCall func() string, result string) {
			Expect(fnToCall()).To(Equal(result), "result should equal")
		},
			Entry("GetDataVolumeManifest should return correct value", (&parse.CLIOptions{DataVolumeManifest: testDVManifest1}).GetDataVolumeManifest, testDVManifest1),
			Entry("GetSourceTemplateNamespace should return correct value", (&parse.CLIOptions{DataVolumeNamespace: testStringDVNamespace1}).GetDataVolumeNamespace, testStringDVNamespace1),
		)

		DescribeTable("GetWaitForSuccess should return correct values", func(fnToCall func() bool, result bool) {
			Expect(fnToCall()).To(Equal(result), "result should equal")
		},
			Entry("should return correct true", (&parse.CLIOptions{WaitForSuccess: "true"}).GetWaitForSuccess, true),
			Entry("should return correct false", (&parse.CLIOptions{WaitForSuccess: "false"}).GetWaitForSuccess, false),
			Entry("should return correct false, when wrong string", (&parse.CLIOptions{WaitForSuccess: "notAValue"}).GetWaitForSuccess, false),
		)

		DescribeTable("CLI options should return correct log level", func(options *parse.CLIOptions, level zapcore.Level) {
			Expect(options.GetDebugLevel()).To(Equal(level), "level should equal")
		},
			Entry("GetDebugLevel should return correct debug level", (&parse.CLIOptions{Debug: true}), zapcore.DebugLevel),
			Entry("GetDebugLevel should return correct info level", (&parse.CLIOptions{Debug: false}), zapcore.InfoLevel),
		)

		It("Init should read the namespace from the manifest", func() {
			options := &parse.CLIOptions{
				DataVolumeManifest: " " + testDVManifest2 + " ",
			}
			Expect(options.Init()).To(Succeed())
			Expect(options.DataVolumeNamespace).To(Equal(testStringDVNamespace2), "DataVolumeNamespace should equal")
		})

		It("Init should try to get the active namespace", func() {
			options := &parse.CLIOptions{
				DataVolumeManifest: " " + testDVManifest1 + " ",
			}

			err := options.Init()
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("can't get active namespace: could not detect active namespace"))
		})
	})
})
