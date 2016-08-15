/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/util/wait"
	"k8s.io/kubernetes/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	FederatedIngressTimeout = 60 * time.Second
)

var _ = framework.KubeDescribe("Federation ingresses [Feature:Federation]", func() {
	defer GinkgoRecover()

	var (
		clusters                               map[string]*cluster // All clusters, keyed by cluster name
		primaryClusterName, federationName, ns string
		gceController                          *GCEIngressController
		jig                                    *testJig
		conformanceTests                       []conformanceTests
	)

	f := framework.NewDefaultFederatedFramework("federation-ingress")

	// e2e cases for federation ingress controller
	var _ = Describe("Federated Ingresses", func() {
		// register clusters in federation apiserver
		BeforeEach(func() {
			framework.SkipUnlessFederated(f.Client)
			if federationName = os.Getenv("FEDERATION_NAME"); federationName == "" {
				federationName = DefaultFederationName
			}
			clusters = map[string]*cluster{}
			primaryClusterName = registerClusters(clusters, UserAgentName, federationName, f)

			// TODO: require ingress to be supported by extensions client
			jig = newTestJig(f.FederationClientset_1_4.ExtensionsClient.RESTClient)
			ns = f.Namespace.Name
		})

		// create backend pod, service and ingress,
		// this requires federation replicaset controller, service controller and ingress controller working properly
		conformanceTests = createComformanceTests(jig, ns)

		AfterEach(func() {
			unregisterClusters(clusters, f)
		})

		Describe("Ingress creation", func() {

			BeforeEach(func() {
				framework.SkipUnlessFederated(f.Client)
				framework.SkipUnlessProviderIs("gce", "gke")
				By("Initializing gce controller")
				//TODO: we do not need federation ingress controller, just check underlying k8s ingress objects are created
				gceController = &GCEIngressController{ns: ns, Project: framework.TestContext.CloudConfig.ProjectID, c: jig.client}
				gceController.init()
			})

			AfterEach(func() {
				framework.SkipUnlessFederated(f.Client)
				if CurrentGinkgoTestDescription().Failed {
					describeIng(ns)
				}
				if jig.ing == nil {
					By("No ingress created, no cleanup necessary")
					return
				}
			})

			It("should succeed", func() {
				framework.SkipUnlessFederated(f.Client)
				for _, t := range conformanceTests {
					By(t.entryLog)
					t.execute()
					By(t.exitLog)
					// confirm corresponding objects in underlying clusters are created, and spec are equal
					waitForIngressShardsOrFail(f.Namespace.Name, jig.ing, clusters)
					// confirm the traffic are working properly, which means the underlying services, backend pods, ingress have been handled properly
					jig.waitForIngress()

				}
			})

			It("ingress should be rebuilt when deletion in underlying cluster", func() {
				for _, c := range clusters {
					err := c.ExtensionsClient.RESTClient.Delete().Namespace(ns).Resource(ingressResourceName).Name(jig.ing.Name).Do().Error()
					if err != nil {
						framework.Failf("Unable to delete underlying ingress: %v", err)
					}
					break
				}
				// the delete ingress should be recovered by federation ingress controller, so wait and
				// confirm corresponding objects in underlying clusters are created, and spec are equal
				waitForIngressShardsOrFail(f.Namespace.Name, jig.ing, clusters)
				// confirm the traffic are working properly, which means the underlying services, backend pods, ingress have been handled properly
				jig.waitForIngress()
			})

			It("shoud create ingress with given static-ip ", func() {
				ip := gceController.staticIP(ns)
				By(fmt.Sprintf("allocated static ip %v: %v through the GCE cloud provider", ns, ip))

				jig.createIngress(filepath.Join(ingressManifestPath, "static-ip"), ns, map[string]string{
					"kubernetes.io/ingress.global-static-ip-name": ns,
					"kubernetes.io/ingress.allow-http":            "false",
				})

				// confirm corresponding objects in underlying clusters are created, and spec are equal
				waitForIngressShardsOrFail(f.Namespace.Name, jig.ing, clusters)

				By("waiting for Ingress to come up with ip: " + ip)
				httpClient := buildInsecureClient(reqTimeout)
				ExpectNoError(jig.pollURL(fmt.Sprintf("https://%v/", ip), "", lbPollTimeout, httpClient, false))

				By("should reject HTTP traffic")
				ExpectNoError(jig.pollURL(fmt.Sprintf("http://%v/", ip), "", lbPollTimeout, httpClient, true))
			})

			It("federation ingress deletion should delete all underlying ingresses", func() {
				err := jig.client.Delete().Namespace(ns).Resource(ingressResourceName).Name(jig.ing.Name).Do().Error()
				if err != nil {
					framework.Failf("Unable to delete federation ingress: %v", err)
				}
				// all underlying ingress should be deleted by federation ingress controller
				waitForIngressShardsGoneOrFail(f.Namespace.Name, jig.ing, clusters)
			})
		})
	})
})

/*
   equivalent returns true if the two ingresss are equivalent.  Fields which are expected to differ between
   federated ingresss and the underlying cluster ingresss (e.g. ClusterIP, LoadBalancerIP etc) are ignored.
*/
func equivalentIngress(federationIngress, clusterIngress extensions.Ingress) bool {
	return reflect.DeepEqual(clusterIngress.Spec, federationIngress.Spec)
}

/*
   waitForIngressOrFail waits until a ingress is either present or absent in the cluster specified by clientset.
   If the condition is not met within timout, it fails the calling test.
*/
func waitForIngressOrFail(client *restclient.RESTClient, namespace string, ingress *extensions.Ingress, present bool, timeout time.Duration) {
	By(fmt.Sprintf("Fetching a federated ingress shard of ingress %q in namespace %q from cluster", ingress.Name, namespace))
	var clusterIngress *extensions.Ingress
	var ok bool
	err := wait.PollImmediate(framework.Poll, timeout, func() (bool, error) {
		o, err := client.Get().Namespace(namespace).Resource(ingressResourceName).Name(ingress.Name).Do().Get()
		if clusterIngress, ok = o.(*extensions.Ingress); !ok {
			return false, nil
		}
		if (!present) && errors.IsNotFound(err) { // We want it gone, and it's gone.
			By(fmt.Sprintf("Success: shard of federated ingress %q in namespace %q in cluster is absent", ingress.Name, namespace))
			return true, nil // Success
		}
		if present && err == nil { // We want it present, and the Get succeeded, so we're all good.
			By(fmt.Sprintf("Success: shard of federated ingress %q in namespace %q in cluster is present", ingress.Name, namespace))
			return true, nil // Success
		}
		By(fmt.Sprintf("Ingress %q in namespace %q in cluster.  Found: %v, waiting for Found: %v, trying again in %s (err=%v)", ingress.Name, namespace, clusterIngress != nil && err == nil, present, framework.Poll, err))
		return false, nil
	})
	framework.ExpectNoError(err, "Failed to verify ingress %q in namespace %q in cluster: Present=%v", ingress.Name, namespace, present)

	if present && clusterIngress != nil {
		Expect(equivalentIngress(*clusterIngress, *ingress))
	}
}

/*
   waitForIngressShardsOrFail waits for the ingress to appear in all clusters
*/
func waitForIngressShardsOrFail(namespace string, ingress *extensions.Ingress, clusters map[string]*cluster) {
	framework.Logf("Waiting for ingress %q in %d clusters", ingress.Name, len(clusters))
	for _, c := range clusters {
		waitForIngressOrFail(c.Clientset.ExtensionsClient.RESTClient, namespace, ingress, true, FederatedIngressTimeout)
	}
}

/*
   waitForIngressShardsGoneOrFail waits for the ingress to disappear in all clusters
*/
func waitForIngressShardsGoneOrFail(namespace string, ingress *extensions.Ingress, clusters map[string]*cluster) {
	framework.Logf("Waiting for ingress %q in %d clusters", ingress.Name, len(clusters))
	for _, c := range clusters {
		waitForIngressOrFail(c.Clientset.ExtensionsClient.RESTClient, namespace, ingress, false, FederatedIngressTimeout)
	}
}

func deleteIngressOrFail(client *restclient.RESTClient, namespace string, ingressName string) {
	var err error
	if client == nil || len(namespace) == 0 || len(ingressName) == 0 {
		Fail(fmt.Sprintf("Internal error: invalid parameters passed to deleteIngressOrFail: clientset: %v, namespace: %v, ingress: %v", client, namespace, ingressName))
	}
	err = client.Delete().Namespace(namespace).Resource(ingressResourceName).Do().Error()
	framework.ExpectNoError(err, "Error deleting ingress %q from namespace %q", ingressName, namespace)
}
