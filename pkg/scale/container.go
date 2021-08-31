package scale

import (
	"fmt"

	"github.com/kubernetes-sigs/cri-tools/pkg/framework"
	internalapi "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"

	. "github.com/onsi/ginkgo"
        gmeasure "github.com/onsi/gomega/gmeasure"
)

var _ = framework.KubeDescribe("Container", func() {
	
	f := framework.NewDefaultCRIFramework()

	var rc internalapi.RuntimeService
	var ic internalapi.ImageManagerService

	BeforeEach(func() {
        	rc = f.CRIClient.CRIRuntimeClient
		ic = f.CRIClient.CRIImageClient
	})

	Context("Start containers in batches", func() {
		var batchSizeStart = 10
		var batchSizeIncrement = 5
                var batchSizeMax = 100


		It("Start containers in increments", func() {
			var containerID string
			var podID string
			var podConfig *runtimeapi.PodSandboxConfig

			var containerIDs []string
			var podIDs []string
			var podConfigs []*runtimeapi.PodSandboxConfig

		       	experiment := gmeasure.NewExperiment("Increments")

			for i := batchSizeStart; i <= batchSizeMax; i= i + batchSizeIncrement {
				By(fmt.Sprintf("Creating %v containers.", i ))
				
				//this should go in a helper func
				for j := 0 ; j < i ; j = j + 1 {
					podID, podConfig = framework.CreatePodSandboxForContainer(rc)
					podIDs = append(podIDs, podID)
					podConfigs = append(podConfigs, podConfig)

					containerID = framework.CreateDefaultContainer(rc, ic, podID, podConfig, "batch-")
					containerIDs = append(containerIDs, containerID)
				}
				
				var err error
				stopwatch := experiment.NewStopwatch()
				for j := 0 ; j < i ; j = j + 1 {
					err = rc.StartContainer(containerIDs[j])
					framework.ExpectNoError(err, "Failed starting container number %v in batch", j)
				}
				stopwatch.Record("Starting")
				
				//cleanup - helper func also
				for j := 0 ; j < i ; j = j + 1 {
					_ = rc.StopContainer(containerIDs[j], framework.DefaultStopContainerTimeout)
					_ = rc.RemoveContainer(containerIDs[j])
			                _ = rc.StopPodSandbox(podIDs[j])
					_ = rc.RemovePodSandbox(podIDs[j])
				}
				containerIDs = nil
				podIDs = nil
				podConfigs = nil // nu cred ca e nevoie de pod configs
                               
			}




		})
	})

})
