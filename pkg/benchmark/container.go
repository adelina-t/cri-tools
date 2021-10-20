//go:build container
// +build container

package benchmark

import (
	"github.com/kubernetes-sigs/cri-tools/pkg/framework"
	"github.com/onsi/gomega/gmeasure"
	internalapi "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"

	. "github.com/onsi/ginkgo"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

type ContainerExperimentData struct {

	CreateContainer, StatusContainer, StopContainer, RemoveContainer, StartContainer string
}

var _ = framework.KubeDescribe("Container", func() {

	f := framework.NewDefaultCRIFramework()

	var rc internalapi.RuntimeService
	var ic internalapi.ImageManagerService

	BeforeEach(func() {

		rc = f.CRIClient.CRIRuntimeClient
		ic = f.CRIClient.CRIImageClient
	})

	Context("benchmark about operations on Container", func() {
		It("benchmark about basic operations on Container", func() {
			experiment := gmeasure.NewExperiment("ContainerOps")
			experiment.Sample(func(idx int) {
				var podID string
				var podConfig *runtimeapi.PodSandboxConfig
				var containerID string
				var err error

				podID, podConfig = framework.CreatePodSandboxForContainer(rc)

				By("CreatingContainer")
				stopwatch := experiment.NewStopwatch()
				stopwatch.Reset()
				containerID = framework.CreateDefaultContainer(rc, ic, podID, podConfig, "Benchmark-container-")
				stopwatch.Record("CreateContainer")

				By("StartingContainer")
				stopwatch.Reset()
				err = rc.StartContainer(containerID)
				stopwatch.Record("StartContainer")
				framework.ExpectNoError(err, "failed to start Container: %v", err)

				By("ContainerStatus")
				stopwatch.Reset()
				_, err = rc.ContainerStatus(containerID)
				stopwatch.Record("StatusContainer")
				framework.ExpectNoError(err, "failed to get Container status: %v", err)

				By("ContainerStop")
				stopwatch.Reset()
				err = rc.StopContainer(containerID, framework.DefaultStopContainerTimeout)
				stopwatch.Record("StopContainer")
				framework.ExpectNoError(err, "failed to stop Container: %v", err)

				By("ContainerRemove")
				stopwatch.Reset()
				err = rc.RemoveContainer(containerID)
				stopwatch.Record("RemoveContainer")
				framework.ExpectNoError(err, "failed to remove Container: %v", err)

				By("stop PodSandbox")
				rc.StopPodSandbox(podID)
				By("delete PodSandbox")
				rc.RemovePodSandbox(podID)

			}, gmeasure.SamplingConfig{N: 200, NumParallel: 1})

			data := ContainerExperimentData{
				CreateContainer: fmt.Sprintf("%v", experiment.Get("CreateContainer").Durations),
				StatusContainer: fmt.Sprintf("%v", experiment.Get("StatusContainer").Durations),
				StopContainer: fmt.Sprintf("%v", experiment.Get("StopContainer").Durations),
				RemoveContainer: fmt.Sprintf("%v", experiment.Get("RemoveContainer").Durations),
				StartContainer: fmt.Sprintf("%v", experiment.Get("StartContainer").Durations), 
			}

			file, _ := json.MarshalIndent(data, "", " ")
			_ = ioutil.WriteFile("c:/experiment_container.json", file, 0644)
		})

	})
})
