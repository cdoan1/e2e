package main_test

import (
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"

	"k8s.io/klog"
)

var _ = Describe("Given a hub cluster web console", func() {

	var page *agouti.Page
	var console, console2, login string
	var version string

	BeforeEach(func() {
		var err error

		console = "https://multicloud-console.apps." + baseDomain + "/multicloud/"
		console2 = "https://multicloud-console.apps." + baseDomain + "/"
		login = "https://oauth-openshift.apps." + baseDomain + "/login"

		defaultOptions := []string{
			"ignore-certificate-errors",
			"disable-gpu",
			"no-sandbox",
			"incognito",
			"window-size=1280,1024",
		}

		if testHeadless {
			defaultOptions = append(defaultOptions, "headless")
		}

		page, err = agoutiDriver.NewPage(agouti.Desired(agouti.Capabilities{
			"chromeOptions": map[string][]string{
				"args": defaultOptions,
			},
		}))

		Expect(err).NotTo(HaveOccurred())
		// SetDefaultEventuallyTimeout(testUITimeout)
	})

	AfterEach(func() {
		Expect(page.Destroy()).To(Succeed())
	})

	It("should allow the user to login to web console (2.1.3, 2.3.0) ", func() {

		By("redirecting the user to the OpenShift login form", func() {
			Expect(page.Navigate(console)).To(Succeed())
			if Expect(page.URL()).To(ContainSubstring("/oauth")) {
				page.AllByClass("idp").At(testIdentityProvider).Click()
			}
			Expect(page.URL()).To(HavePrefix(login))
		})

		By("allowing the user to fill out the login form and submit it", func() {
			Eventually(page.FindByID("inputUsername")).Should(BeFound())
			Eventually(page.FindByID("inputPassword")).Should(BeFound())
			Expect(page.FindByID("inputUsername").Fill(kubeadminUser)).To(Succeed())
			Expect(page.FindByID("inputPassword").Fill(kubeadminCredential)).To(Succeed())
			_, err := page.FindByClass("form-horizontal").Active()
			if err != nil {
				Expect(page.FindByClass("pf-c-form").Submit()).To(Succeed())
				// Expect(page.FindByClass("pf-c-button").Submit()).To(Succeed())
			} else {
				Expect(page.FindByClass("form-horizontal").Submit()).To(Succeed())
			}
		})

		By("viewing the Getting Started page after a successful login", func() {
			Expect(page).To(HaveURL(getConsoleURL(console, "/")))

			
			page.Refresh()

			// wait(120)

			// Eventually(page.FindByClass("clusters")).Should(BeFound())
			// Expect(page.FindByID("acm-info-dropdown").Click())
			// Expect(page.FindByID("acm-about").Click())
			// Eventually(page.FindByClass("version-details")).Should(BeFound())
			// Eventually(page.FindByClass("version-details__no")).Should(BeFound())
			
			// wait for the version to populate
			wait(2)

			// version, _ = page.FindByClass("version-details__no").Text()
			// if version == "2.1.1" {
			// 	klog.V(1).Infof("A version: %s ...", version)
			// } else {
			// 	klog.V(1).Infof("B version: %s ...", version)
			// }

			Expect(page.Screenshot("./results/.test.login.screenshot.png")).To(Succeed())
		})
	})

	When("the user is already authenticated", func() {

		BeforeEach(func() {
			By("redirecting the user to the OpenShift login form", func() {
				Expect(page.Navigate(console)).To(Succeed())
				if Expect(page.URL()).To(ContainSubstring("/oauth")) {
					page.AllByClass("idp").At(testIdentityProvider).Click()
				}
				Expect(page.URL()).To(HavePrefix(login))
			})

			By("allowing the user to fill out the login form and submit it", func() {
				Eventually(page.FindByID("inputUsername")).Should(BeFound())
				Eventually(page.FindByID("inputPassword")).Should(BeFound())
				Expect(page.FindByID("inputUsername").Fill(kubeadminUser)).To(Succeed())
				Expect(page.FindByID("inputPassword").Fill(kubeadminCredential)).To(Succeed())
				_, err := page.FindByClass("form-horizontal").Active()
				if err != nil {
					Expect(page.FindByClass("pf-c-form").Submit()).To(Succeed())
				} else {
					Expect(page.FindByClass("form-horizontal").Submit()).To(Succeed())
				}
			})
		})

		// NOTE: we're looking for elements by class when we really need to be looking by ID
		//       ID are more static!

		It("should allow the user to navigate and view the Overview page (mvp, 2.3.0)", func() {
			By("navigating to /multicloud/overview", func() {
				Expect(page.Navigate(getConsoleURL(console, "/overview"))).To(Succeed())
				Expect(page).To(HaveURL(getConsoleURL(console2, "/overview")))
				// Eventually(page.FindByClass("overview-header-title")).Should(BeFound())
				// wait(5)
				Expect(page.Screenshot("./results/.test.overview.screenshot.png")).To(Succeed())
			})
		})

		It("should allow the user to navigate and view the Toplogy page (mvp)", func() {
			By("navigating to /multicloud/topology", func() {
				Expect(page.Navigate(getConsoleURL(console, "/topology"))).To(Succeed())
				Expect(page).To(HaveURL(getConsoleURL(console, "/topology/")))
				Expect(page.FindByClass("bx--detail-page-header-title")).Should(HaveText("Topology"))
				wait(5)
				Expect(page.Screenshot("./results/.test.topology.screenshot.png")).To(Succeed())
			})
		})

		// 2.2 OK, 2.1 OK
		It("should allow the user to navigate and view the Clusters page (2.1,2.1.1,2.1.3)", func() {
			By("navigating to /multicloud/clusters", func() {
				Expect(page.Navigate(getConsoleURL(console, "/clusters"))).To(Succeed())
				Expect(page).To(HaveURL(getConsoleURL(console, "/clusters")))
				Expect(page.FindByClass("bx--detail-page-header-title")).Should(HaveText("Clusters"))
				Eventually(page.FindByClass("create-import-cluster-dropdown"), 5*time.Second).Should(BeFound())
				Expect(page.Screenshot("./results/.test.cluster.screenshot.png")).To(Succeed())
			})
		})

		It("should allow the user to navigate and view the Application page (2.1.3)", func() {
			By("navigating to /multicloud/applications", func() {
				Expect(page.Navigate(getConsoleURL(console, "/applications"))).To(Succeed())
				Expect(page).To(HaveURL(getConsoleURL(console, "/applications/")))
			})
			By("should have header label Application", func() {
				Expect(page.FindByClass("bx--detail-page-header-title")).Should(HaveText("Applications"))
			})
			By("should have a button name Create application", func() {
				Eventually(page.FindByButton("Create application")).Should(BeFound())
			})
			By("should have a table with applications", func() {

				if version == "2.1.1" {
					Eventually(page.FirstByClass("bx--data-table-v2")).Should(BeFound())
					Eventually(page.FirstByClass("bx--pagination__left").FindByClass("bx--select--inline").FindByClass("bx--select-input"), 60*time.Second).Should(BeFound())
				} else {
					// 2.2
					Eventually(page.Find("table")).Should(BeVisible())
					Expect(page.Screenshot("./results/.test.application.screenshot.png")).To(Succeed())
				}
			})
			By("screen capture", func() {
				Expect(page.Screenshot("./results/.test.application.screenshot.png")).To(Succeed())
			})
		})

		It("should allow the user to navigate and view the Policy page (2.2)", func() {
			By("navigating to /multicloud/policies/all", func() {
				Expect(page.Navigate(getConsoleURL(console, "/policies/all"))).To(Succeed())
				Expect(page).To(HaveURL(getConsoleURL(console, "/policies/all")))
				Expect(page.FindByClass("bx--detail-page-header-title")).Should(HaveText("Governance and risk"))
				Expect(page.FindByButton("Create policy")).Should(BeFound())

				// 2.2 specific begin
				Eventually(page.FindByClass("grc-view-by-policies-table")).Should(BeFound())
				Eventually(page.FindByButton("Policies")).Should(BeFound())
				// 2.2 specific end

				Expect(page.Screenshot("./results/.test.policies.screenshot.png")).To(Succeed())
			})
		})

		It("should allow the user to navigate and view the Policy page (2.1,2.1.2)", func() {
			By("navigating to /multicloud/policies/all", func() {
				Expect(page.Navigate(getConsoleURL(console, "/policies/all"))).To(Succeed())
				Expect(page).To(HaveURL(getConsoleURL(console, "/policies/all")))
				Expect(page.FindByClass("bx--detail-page-header-title")).Should(HaveText("Governance and risk"))
				Expect(page.FindByButton("Create policy")).Should(BeFound())

				// 2.1 specific begin
				Eventually(page.FirstByClass("bx--data-table-v2")).Should(BeFound())
				Eventually(page.FirstByClass("bx--pagination__left").FindByClass("bx--select--inline").FindByClass("bx--select-input"), 60*time.Second).Should(BeFound())
				// 2.1 specific end

				Expect(page.Screenshot("./results/.test.policies.screenshot.png")).To(Succeed())
			})
		})

	})

})

func getConsoleURL(console, path string) string {
	if strings.HasSuffix(console, "/") {
		if strings.HasPrefix(path, "/") {
			return console + path[1:len(path)]
		}
		return console + path[1:len(path)]
	} else if strings.HasPrefix(path, "/") {
		return console + path
	}
	return console + "/" + path
}

func wait(seconds int) {
	klog.V(1).Infof("waiting %d seconds ...", seconds)
	time.Sleep(time.Duration(seconds*1000) * time.Millisecond)
}
