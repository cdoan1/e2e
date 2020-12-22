package main_test

import (
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Given a hub cluster web console", func() {

	var page *agouti.Page
	var console, login string

	BeforeEach(func() {
		var err error

		console = "https://multicloud-console.apps." + baseDomain + "/multicloud/"
		login = "https://oauth-openshift.apps." + baseDomain + "/login"

		defaultOptions := []string{
			"ignore-certificate-errors",
			"disable-gpu",
			"no-sandbox",
			"incognito",
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

	It("should allow the user to login to web console", func() {

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
			Expect(page).To(HaveURL(getConsoleURL(console, "/welcome")))
			page.Refresh()
			Eventually(page.FindByClass("welcome")).Should(BeFound())
			// wait(5)
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

		It("should allow the user to navigate and view the Overview page (mvp)", func() {
			By("navigating to /multicloud/overview", func() {
				Expect(page.Navigate(getConsoleURL(console, "/overview"))).To(Succeed())
				Expect(page).To(HaveURL(getConsoleURL(console, "/overview")))
				Eventually(page.FindByClass("overview-header-title")).Should(BeFound())
				wait(5)
			})
		})

		It("should allow the user to navigate and view the Toplogy page (mvp)", func() {
			By("navigating to /multicloud/topology", func() {
				Expect(page.Navigate(getConsoleURL(console, "/topology"))).To(Succeed())
				Expect(page).To(HaveURL(getConsoleURL(console, "/topology/")))
				Expect(page.FindByClass("bx--detail-page-header-title")).Should(HaveText("Topology"))
				wait(5)
			})
		})

		It("should allow the user to navigate and view the Clusters page (mvp)", func() {
			By("navigating to /multicloud/clusters", func() {
				Expect(page.Navigate(getConsoleURL(console, "/clusters"))).To(Succeed())
				Expect(page).To(HaveURL(getConsoleURL(console, "/clusters")))
				Expect(page.FindByClass("bx--detail-page-header-title")).Should(HaveText("Clusters"))
				Expect(page.FindByClass("create-import-cluster-dropdown")).Should(BeFound())
				wait(5)
			})
		})

		It("should allow the user to navigate and view the Application page (mvp)", func() {
			By("navigating to /multicloud/applications", func() {
				Expect(page.Navigate(getConsoleURL(console, "/applications"))).To(Succeed())
				Expect(page).To(HaveURL(getConsoleURL(console, "/applications/")))
				Expect(page.FindByClass("bx--detail-page-header-title")).Should(HaveText("Applications"))
				Expect(page.FindByButton("Create application")).Should(BeFound())
				Expect(page.Find("table")).To(BeVisible())
				// wait(5)
			})
		})

		It("should allow the user to navigate and view the Policy page (mvp)", func() {
			By("navigating to /multicloud/policies/all", func() {
				Expect(page.Navigate(getConsoleURL(console, "/policies/all"))).To(Succeed())
				Expect(page).To(HaveURL(getConsoleURL(console, "/policies/all")))
				Expect(page.FindByClass("bx--detail-page-header-title")).Should(HaveText("Governance and risk"))
				Expect(page.FindByButton("Create policy")).Should(BeFound())
				Eventually(page.FindByClass("grc-view-by-policies-table")).Should(BeFound())
				Eventually(page.FindByButton("Policies")).Should(BeFound())
				// wait(5)
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
	time.Sleep(time.Duration(seconds*1000) * time.Millisecond)
}
