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
	"reflect"
	"strconv"
	"time"

	"k8s.io/kubernetes/federation/client/clientset_generated/federation_release_1_4"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	"k8s.io/kubernetes/pkg/client/clientset_generated/release_1_3"
	"k8s.io/kubernetes/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/util/wait"
	"k8s.io/kubernetes/test/e2e/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	FederatedIngressTimeout        = 60 * time.Second
	FederatedIngressName           = "federated-ingress"
	FederatedIngressServiceName    = "federated-ingress-service"
	FederatedIngressServicePodName = "federated-ingress-service-test-pod"
)

var _ = framework.KubeDescribe("Federation ingresses [Feature:Federation]", func() {
	defer GinkgoRecover()

	var (
		clusters                               map[string]*cluster // All clusters, keyed by cluster name
		primaryClusterName, federationName, ns string
	)

	f := framework.NewDefaultFederatedFramework("federated-ingress")

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
			ns = f.Namespace.Name
		})

		AfterEach(func() {
			unregisterClusters(clusters, f)
		})

		Describe("Ingress objects", func() {
			AfterEach(func() {
				framework.SkipUnlessFederated(f.Client)
				// Delete registered ingresses.
				ingressList, err := f.FederationClientset_1_4.Extensions().Ingresses(f.Namespace.Name).List(api.ListOptions{})
				Expect(err).NotTo(HaveOccurred())
				for _, ingress := range ingressList.Items {
					err := f.FederationClientset_1_4.Extensions().Ingresses(f.Namespace.Name).Delete(ingress.Name, &api.DeleteOptions{})
					Expect(err).NotTo(HaveOccurred())
					for _, c := range clusters {
						err := c.Extensions().Ingresses(f.Namespace.Name).Delete(ingress.Name, &api.DeleteOptions{})
						Expect(err).NotTo(HaveOccurred())
					}
				}
			})

			It("should be created and deleted successfully in federation", func() {
				framework.SkipUnlessFederated(f.Client)
				ingress := createIngressOrFail(f.FederationClientset_1_4, f.Namespace.Name)
				By(fmt.Sprintf("Creation of ingress %q in namespace %q succeeded.  Deleting ingress.", ingress.Name, f.Namespace.Name))
				// Cleanup
				err := f.FederationClientset_1_4.Extensions().Ingresses(f.Namespace.Name).Delete(ingress.Name, &api.DeleteOptions{})
				framework.ExpectNoError(err, "Error deleting ingress %q in namespace %q", ingress.Name, ingress.Namespace)
				By(fmt.Sprintf("Deletion of ingress %q in namespace %q succeeded.", ingress.Name, f.Namespace.Name))
			})

			It("should create and update matching ingresses in underlying clusters", func() {
				framework.SkipUnlessFederated(f.Client)
				ingress := createIngressOrFail(f.FederationClientset_1_4, f.Namespace.Name)
				defer func() { // Cleanup
					By(fmt.Sprintf("Deleting ingress %q in namespace %q", ingress.Name, f.Namespace.Name))
					err := f.FederationClientset_1_4.Ingresses(f.Namespace.Name).Delete(ingress.Name, &api.DeleteOptions{})
					framework.ExpectNoError(err, "Error deleting ingress %q in namespace %q", ingress.Name, f.Namespace.Name)
				}()
				// wait for ingress shards being created
				waitForIngressShardsOrFail(f.Namespace.Name, ingress, clusters)
				ingress = updateIngressOrFail(f.FederationClientset_1_4, f.Namespace.Name)
				waitForIngressShardsUpdatedOrFail(f.Namespace.Name, ingress, clusters)
			})
		})

		var _ = Describe("DNS", func() {

			var (
				service *v1.Service
				ingress *v1beta1.Ingress
			)

			BeforeEach(func() {
				framework.SkipUnlessFederated(f.Client)
				// create backend pod
				createBackendPodsOrFail(clusters, f.Namespace.Name, FederatedIngressServicePodName)
				// create backend service
				service = createServiceOrFail(f.FederationClientset_1_4, f.Namespace.Name, FederatedIngressServiceName)
				// create ingress object
				ingress = createIngressOrFail(f.FederationClientset_1_4, f.Namespace.Name)
				// wait for services objects sync
				waitForServiceShardsOrFail(f.Namespace.Name, service, clusters)
				// wait for ingress objects sync
				waitForIngressShardsOrFail(f.Namespace.Name, ingress, clusters)
			})

			AfterEach(func() {
				framework.SkipUnlessFederated(f.Client)
				deleteBackendPodsOrFail(clusters, f.Namespace.Name)
				if service != nil {
					deleteServiceOrFail(f.FederationClientset_1_4, f.Namespace.Name, service.Name)
					service = nil
				} else {
					By("No service to delete. Service is nil")
				}
				if ingress != nil {
					deleteIngressOrFail(f.FederationClientset_1_4, f.Namespace.Name, ingress.Name)
					ingress = nil
				} else {
					By("No ingress to delete. Ingress is nil")
				}
			})

			It("should be able to discover a federated ingress service", func() {
				framework.SkipUnlessFederated(f.Client)
				// we are about the ingress name
				svcDNSNames := []string{
					FederatedIngressServiceName,
					fmt.Sprintf("%s.%s", FederatedIngressServiceName, f.Namespace.Name),
					fmt.Sprintf("%s.%s.svc.cluster.local.", FederatedIngressServiceName, f.Namespace.Name),
					fmt.Sprintf("%s.%s.%s", FederatedIngressServiceName, f.Namespace.Name, federationName),
					fmt.Sprintf("%s.%s.%s.svc.cluster.local.", FederatedIngressServiceName, f.Namespace.Name, federationName),
				}
				for i, DNSName := range svcDNSNames {
					discoverService(f, DNSName, true, "federated-ingress-e2e-discovery-pod-"+strconv.Itoa(i))
				}
			})

			Context("non-local federated ingress service", func() {
				BeforeEach(func() {
					framework.SkipUnlessFederated(f.Client)

					// Delete all the backend pods from the shard which is local to the discovery pod.
					deleteOneBackendPodOrFail(clusters[primaryClusterName])

				})

				It("should be able to discover a non-local federated ingress", func() {
					framework.SkipUnlessFederated(f.Client)

					svcDNSNames := []string{
						fmt.Sprintf("%s.%s.%s", FederatedIngressServiceName, f.Namespace.Name, federationName),
						fmt.Sprintf("%s.%s.%s.svc.cluster.local.", FederatedIngressServiceName, f.Namespace.Name, federationName),
					}
					for i, name := range svcDNSNames {
						discoverService(f, name, true, "federated-ingress-e2e-discovery-pod-"+strconv.Itoa(i))
					}
				})

				// TODO(mml): This currently takes 9 minutes.  Consider reducing the
				// TTL and/or running the pods in parallel.
				Context("[Slow] missing local ingress", func() {
					It("should never find DNS entries for a missing local ingress", func() {
						framework.SkipUnlessFederated(f.Client)

						localSvcDNSNames := []string{
							FederatedIngressServiceName,
							fmt.Sprintf("%s.%s", FederatedIngressServiceName, f.Namespace.Name),
							fmt.Sprintf("%s.%s.svc.cluster.local.", FederatedIngressServiceName, f.Namespace.Name),
						}
						for i, name := range localSvcDNSNames {
							discoverService(f, name, false, "federated-ingress-e2e-discovery-pod-"+strconv.Itoa(i))
						}
					})
				})
			})
		})
	})
})

/*
   equivalent returns true if the two ingresss are equivalent.  Fields which are expected to differ between
   federated ingresss and the underlying cluster ingresss (e.g. ClusterIP, LoadBalancerIP etc) are ignored.
*/
func equivalentIngress(federatedIngress, clusterIngress v1beta1.Ingress) bool {
	return reflect.DeepEqual(clusterIngress.Spec, federatedIngress.Spec)
}

/*
   waitForIngressOrFail waits until a ingress is either present or absent in the cluster specified by clientset.
   If the condition is not met within timout, it fails the calling test.
*/
func waitForIngressOrFail(clientset *release_1_3.Clientset, namespace string, ingress *v1beta1.Ingress, present bool, timeout time.Duration) {
	By(fmt.Sprintf("Fetching a federated ingress shard of ingress %q in namespace %q from cluster", ingress.Name, namespace))
	var clusterIngress *v1beta1.Ingress
	err := wait.PollImmediate(framework.Poll, timeout, func() (bool, error) {
		clusterIngress, err := clientset.Ingresses(namespace).Get(ingress.Name)
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
func waitForIngressShardsOrFail(namespace string, ingress *v1beta1.Ingress, clusters map[string]*cluster) {
	framework.Logf("Waiting for ingress %q in %d clusters", ingress.Name, len(clusters))
	for _, c := range clusters {
		waitForIngressOrFail(c.Clientset, namespace, ingress, true, FederatedIngressTimeout)
	}
}

/*
   waitForIngressShardsOrFail waits for the ingress to appear in all clusters
*/
func waitForIngressShardsUpdatedOrFail(namespace string, ingress *v1beta1.Ingress, clusters map[string]*cluster) {
	framework.Logf("Waiting for ingress %q in %d clusters", ingress.Name, len(clusters))
	for _, c := range clusters {
		waitForIngressUpdateOrFail(c.Clientset, namespace, ingress, FederatedIngressTimeout)
	}
}

/*
   waitForIngressOrFail waits until a ingress is either present or absent in the cluster specified by clientset.
   If the condition is not met within timout, it fails the calling test.
*/
func waitForIngressUpdateOrFail(clientset *release_1_3.Clientset, namespace string, ingress *v1beta1.Ingress, timeout time.Duration) {
	By(fmt.Sprintf("Fetching a federated ingress shard of ingress %q in namespace %q from cluster", ingress.Name, namespace))
	err := wait.PollImmediate(framework.Poll, timeout, func() (bool, error) {
		clusterIngress, err := clientset.Ingresses(namespace).Get(ingress.Name)
		if err == nil { // We want it present, and the Get succeeded, so we're all good.
			if equivalentIngress(*clusterIngress, *ingress) {
				By(fmt.Sprintf("Success: shard of federated ingress %q in namespace %q in cluster is updated", ingress.Name, namespace))
				return true, nil
			}
			By(fmt.Sprintf("Ingress %q in namespace %q in cluster, waiting for service being updated, trying again in %s (err=%v)", ingress.Name, namespace, framework.Poll, err))
			return false, nil
		}
		By(fmt.Sprintf("Ingress %q in namespace %q in cluster, waiting for service being updated, trying again in %s (err=%v)", ingress.Name, namespace, framework.Poll, err))
		return false, nil
	})
	framework.ExpectNoError(err, "Failed to verify ingress %q in namespace %q in cluster", ingress.Name, namespace)
}

/*
   waitForIngressShardsGoneOrFail waits for the ingress to disappear in all clusters
*/
func waitForIngressShardsGoneOrFail(namespace string, ingress *v1beta1.Ingress, clusters map[string]*cluster) {
	framework.Logf("Waiting for ingress %q in %d clusters", ingress.Name, len(clusters))
	for _, c := range clusters {
		waitForIngressOrFail(c.Clientset, namespace, ingress, false, FederatedIngressTimeout)
	}
}

func deleteIngressOrFail(clientset *federation_release_1_4.Clientset, namespace string, ingressName string) {
	if clientset == nil || len(namespace) == 0 || len(ingressName) == 0 {
		Fail(fmt.Sprintf("Internal error: invalid parameters passed to deleteIngressOrFail: clientset: %v, namespace: %v, ingress: %v", clientset, namespace, ingressName))
	}
	err := clientset.Ingresses(namespace).Delete(ingressName, api.NewDeleteOptions(0))
	framework.ExpectNoError(err, "Error deleting ingress %q from namespace %q", ingressName, namespace)
}

func createIngressOrFail(clientset *federation_release_1_4.Clientset, namespace string) *v1beta1.Ingress {
	if clientset == nil || len(namespace) == 0 {
		Fail(fmt.Sprintf("Internal error: invalid parameters passed to createIngressOrFail: clientset: %v, namespace: %v", clientset, namespace))
	}
	By(fmt.Sprintf("Creating federated ingress %q in namespace %q", FederatedIngressName, namespace))

	ingress := &v1beta1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name: FederatedIngressName,
		},
		Spec: v1beta1.IngressSpec{
			Backend: &v1beta1.IngressBackend{
				ServiceName: "testingress-service",
				ServicePort: intstr.FromInt(80),
			},
		},
	}

	_, err := clientset.Extensions().Ingresses(namespace).Create(ingress)
	framework.ExpectNoError(err, "Creating ingress %q in namespace %q", ingress.Name, namespace)
	By(fmt.Sprintf("Successfully created federated ingress %q in namespace %q", FederatedIngressName, namespace))
	return ingress
}

func updateIngressOrFail(clientset *federation_release_1_4.Clientset, namespace string) (newIng *v1beta1.Ingress) {
	var err error
	if clientset == nil || len(namespace) == 0 {
		Fail(fmt.Sprintf("Internal error: invalid parameters passed to createIngressOrFail: clientset: %v, namespace: %v", clientset, namespace))
	}
	ingress := &v1beta1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name: FederatedIngressName,
		},
		Spec: v1beta1.IngressSpec{
			Backend: &v1beta1.IngressBackend{
				ServiceName: "updated-testingress-service",
				ServicePort: intstr.FromInt(80),
			},
		},
	}

	for i := 0; i < 3; i++ {
		_, err = clientset.Extensions().Ingresses(namespace).Get(FederatedIngressName)
		if err != nil {
			framework.Failf("failed to get ingress %q: %v", FederatedIngressName, err)
		}
		newIng, err = clientset.Extensions().Ingresses(namespace).Update(ingress)
		if err == nil {
			describeIng(namespace)
			return
		}
		if !errors.IsConflict(err) && !errors.IsServerTimeout(err) {
			framework.Failf("failed to update ingress %q: %v", FederatedIngressName, err)
		}
	}
	framework.Failf("too many retries updating ingress %q", FederatedIngressName)
	return newIng
}
