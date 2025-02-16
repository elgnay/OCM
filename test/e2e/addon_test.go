package e2e

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/util/rand"

	operatorapiv1 "open-cluster-management.io/api/operator/v1"
)

var _ = Describe("Manage the managed cluster addons", func() {
	var klusterletName, clusterName, agentNamespace, addOnName string

	BeforeEach(func() {
		klusterletName = fmt.Sprintf("e2e-klusterlet-%s", rand.String(6))
		clusterName = fmt.Sprintf("e2e-managedcluster-%s", rand.String(6))
		agentNamespace = fmt.Sprintf("open-cluster-management-agent-%s", rand.String(6))
		addOnName = fmt.Sprintf("e2e-addon-%s", rand.String(6))

		_, err := t.CreateApprovedKlusterlet(klusterletName, clusterName, agentNamespace, operatorapiv1.InstallModeDefault)
		Expect(err).ToNot(HaveOccurred())
	})
	AfterEach(func() {
		By(fmt.Sprintf("clean klusterlet %v resources after the test case", klusterletName))
		Expect(t.cleanKlusterletResources(klusterletName, clusterName)).To(BeNil())
	})

	It("Create one managed cluster addon and make sure it is available", func() {
		By(fmt.Sprintf("create the addon %v on the managed cluster namespace %v", addOnName, clusterName))
		err := t.CreateManagedClusterAddOn(clusterName, addOnName, addOnName)
		Expect(err).ToNot(HaveOccurred())

		By(fmt.Sprintf("create the addon lease %v on addon install namespace %v", addOnName, addOnName))
		err = t.CreateManagedClusterAddOnLease(addOnName, addOnName)
		Expect(err).ToNot(HaveOccurred())

		By(fmt.Sprintf("wait the addon %v available condition to be true", addOnName))
		Eventually(func() error {
			return t.CheckManagedClusterAddOnStatus(clusterName, addOnName)
		}, t.EventuallyTimeout*5, t.EventuallyInterval*5).Should(Succeed())
	})

	It("Create one managed cluster addon and make sure it is available in Hosted mode", func() {
		By(fmt.Sprintf("create the addon %v on the managed cluster namespace %v", addOnName, clusterName))
		err := t.CreateManagedClusterAddOn(clusterName, addOnName, addOnName)
		Expect(err).ToNot(HaveOccurred())

		By(fmt.Sprintf("create the addon lease %v on addon install namespace %v", addOnName, addOnName))
		err = t.CreateManagedClusterAddOnLease(addOnName, addOnName)
		Expect(err).ToNot(HaveOccurred())

		By(fmt.Sprintf("wait the addon %v available condition to be true", addOnName))
		Eventually(func() error {
			return t.CheckManagedClusterAddOnStatus(clusterName, addOnName)
		}, t.EventuallyTimeout*5, t.EventuallyInterval*5).Should(Succeed())
	})
})
