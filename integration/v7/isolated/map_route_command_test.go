package isolated

import (
	. "code.cloudfoundry.org/cli/cf/util/testhelpers/matchers"
	"code.cloudfoundry.org/cli/integration/helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("map-route command", func() {
	Context("help", func() {
		It("appears in cf help -a", func() {
			session := helpers.CF("help", "-a")
			Eventually(session).Should(Exit(0))
			Expect(session).To(HaveCommandInCategoryWithDescription("map-route", "ROUTES", "Map a route to an app"))
		})
		It("Displays command usage to output", func() {
			session := helpers.CF("map-route", "--help")
			Eventually(session).Should(Say(`NAME:`))
			Eventually(session).Should(Say(`map-route - Map a route to an app\n`))
			Eventually(session).Should(Say(`\n`))

			Eventually(session).Should(Say(`USAGE:`))
			Eventually(session).Should(Say(`cf map-route APP_NAME DOMAIN \[--hostname HOSTNAME\] \[--path PATH\]\n`))
			Eventually(session).Should(Say(`\n`))

			Eventually(session).Should(Say(`EXAMPLES:`))
			Eventually(session).Should(Say(`cf map-route my-app example.com                              # example.com`))
			Eventually(session).Should(Say(`cf map-route my-app example.com --hostname myhost            # myhost.example.com`))
			Eventually(session).Should(Say(`cf map-route my-app example.com --hostname myhost --path foo # myhost.example.com/foo`))
			Eventually(session).Should(Say(`\n`))

			Eventually(session).Should(Say(`OPTIONS:`))
			Eventually(session).Should(Say(`--hostname, -n\s+Hostname for the HTTP route \(required for shared domains\)`))
			Eventually(session).Should(Say(`--path\s+Path for the HTTP route`))
			Eventually(session).Should(Say(`\n`))

			Eventually(session).Should(Say(`SEE ALSO:`))
			Eventually(session).Should(Say(`create-route, routes, unmap-route`))

			Eventually(session).Should(Exit(0))
		})
	})

	When("the environment is not setup correctly", func() {
		It("fails with the appropriate errors", func() {
			helpers.CheckEnvironmentTargetedCorrectly(true, true, ReadOnlyOrg, "map-route", "app-name", "domain-name")
		})
	})

	When("The environment is set up correctly", func() {
		var (
			orgName    string
			spaceName  string
			domainName string
			hostName   = "jarjar"
			path       = "/binks"
			userName   string
			appName    = helpers.NewAppName()
		)

		BeforeEach(func() {
			orgName = helpers.NewOrgName()
			spaceName = helpers.NewSpaceName()
			helpers.SetupCF(orgName, spaceName)
			userName, _ = helpers.GetCredentials()
			domainName = helpers.DefaultSharedDomain()

			helpers.WithHelloWorldApp(func(dir string) {
				session := helpers.CustomCF(helpers.CFEnv{WorkingDirectory: dir}, "push", appName)
				Eventually(session).Should(Exit(0))
			})
		})

		When("the route already exists", func() {
			var (
				route helpers.Route
			)
			BeforeEach(func() {
				route = helpers.NewRoute(spaceName, domainName, hostName, path)
				route.V7Create()
			})

			AfterEach(func() {
				route.Delete()
			})

			It("maps the route to an app", func() {
				session := helpers.CF("map-route", appName, domainName, "--hostname", route.Host, "--path", route.Path)

				Eventually(session).Should(Say(`Mapping route %s.%s%s to app %s in org %s / space %s as %s\.\.\.`, hostName, domainName, path, appName, orgName, spaceName, userName))
				Eventually(session).Should(Say(`OK`))
				Eventually(session).Should(Exit(0))
			})
		})

		When("the route doesnt exist", func() {
			It("creates the route and maps it to an app", func() {
				session := helpers.CF("map-route", appName, domainName, "--hostname", "test", "--path", "route")
				Eventually(session).Should(Say(`Creating route %s.%s/%s for org %s / space %s as %s\.\.\.`, "test", domainName, "route", orgName, spaceName, userName))
				Eventually(session).Should(Say(`OK`))
				Eventually(session).Should(Say(`Mapping route %s.%s/%s to app %s in org %s / space %s as %s\.\.\.`, "test", domainName, "route", appName, orgName, spaceName, userName))
				Eventually(session).Should(Say(`OK`))
				Eventually(session).Should(Exit(0))
			})
		})

	})
})