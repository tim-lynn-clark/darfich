package darfich

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/tim-lynn-clark/darfich/ability"
	"github.com/tim-lynn-clark/darfich/utils"
)

var (
	version string = "v0.1.2"       //https://go.dev/doc/modules/version-numbers
	build   string = "202307292113" // YYYYMMDDHHMM
)

// Config represents configuration values for the DarfIch package.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool
	// Filter defines a function to skip middleware.
	// Optional. Default: nil
	Filter     func(*fiber.Ctx) bool
	ContextKey string
	RuleSet    *ability.Set
}

// New Create a new middleware handler
func New(config Config) func(*fiber.Ctx) error {
	fmt.Println("version=", version)
	fmt.Println("build=", build)

	// Return new Fiber handler
	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if config.Next != nil && config.Next(c) {
			return c.Next()
		}

		// Filter request to skip middleware
		if config.Filter != nil && config.Filter(c) {
			return c.Next()
		}

		// Pull current user out of the context
		currentUser := c.Locals(config.ContextKey).(utils.DtoCurrentUser)

		method := c.Method()
		path := c.Path()

		// Search through rules using key for matching rule
		var allow bool
		for _, rule := range config.RuleSet.Rules {
			if fiber.RoutePatternMatch(path, string(rule.Route)) &&
				rule.Method == utils.HttpMethod(method) &&
				rule.Role == utils.Role(currentUser.RoleName) {

				if rule.Action == utils.ActionAllow {
					allow = true
				}
			}
		}

		// If rule is found, and action is allow, continue
		if allow {
			return c.Next()
		}
		// If no rule is found, return 403 Forbidden
		return c.SendStatus(fiber.StatusForbidden)
	}
}
