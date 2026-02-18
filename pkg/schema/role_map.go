package schema

// RoleMap maps OMC agent names to canonical roles.
// See references/event-schema.md section 11.
var RoleMap = map[string]Role{
	"planner":                RolePlanner,
	"executor":               RoleExecutor,
	"deep-executor":          RoleExecutor,
	"explore":                RoleExplorer,
	"architect":              RoleArchitect,
	"debugger":               RoleDebugger,
	"verifier":               RoleVerifier,
	"designer":               RoleDesigner,
	"code-reviewer":          RoleReviewer,
	"style-reviewer":         RoleReviewer,
	"quality-reviewer":       RoleReviewer,
	"api-reviewer":           RoleReviewer,
	"performance-reviewer":   RoleReviewer,
	"security-reviewer":      RoleGuard,
	"test-engineer":          RoleTester,
	"writer":                 RoleWriter,
	"analyst":                RolePlanner,
	"product-manager":        RolePlanner,
	"product-analyst":        RolePlanner,
	"ux-researcher":          RolePlanner,
	"information-architect":  RolePlanner,
	"build-fixer":            RoleExecutor,
	"scientist":              RoleExplorer,
	"dependency-expert":      RoleExplorer,
	"git-master":             RoleExecutor,
	"qa-tester":              RoleTester,
	"critic":                 RoleReviewer,
}

// LookupRole returns the canonical role for an OMC agent name.
// Unknown agents map to RoleCustom with a warning flag.
func LookupRole(agentName string) (Role, bool) {
	if r, ok := RoleMap[agentName]; ok {
		return r, true
	}
	return RoleCustom, false
}
