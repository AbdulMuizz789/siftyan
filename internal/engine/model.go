package engine

// Model defines how conflict rules are applied based on distribution
type Model string

const (
	ModelSaaS     Model = "saas"
	ModelBinary   Model = "binary"
	ModelInternal Model = "internal"
)

// Rule defines a license compatibility rule
type Rule struct {
	TargetType LicenseType
	Impact     string
	Description string
	Suggestions []string
}

// GetRulesForModel returns the conflict rules based on the distribution model
func GetRulesForModel(m string) []Rule {
	model := Model(m)
	switch model {
	case ModelBinary:
		return []Rule{
			{StrongCopyleftLT, "HIGH", "Strong Copyleft (GPL) in a binary distribution triggers the 'viral' clause.", []string{"Switch project license to GPL", "Replace dependency"}},
			{NetworkCopyleftLT, "HIGH", "Network Copyleft (AGPL) is a high risk in distributed software.", []string{"Replace with a permissive alternative", "Consult legal"}},
			{WeakCopyleftLT, "MEDIUM", "LGPL detected. Ensure you are using dynamic linking to comply with the linking exception.", []string{"Verify dynamic linking", "Replace if static linking is required"}},
		}
	case ModelSaaS:
		return []Rule{
			{NetworkCopyleftLT, "HIGH", "AGPL-3.0 affects SaaS. You may be required to disclose source code even if not distributing binaries.", []string{"Replace dependency", "Isolate as a microservice"}},
			{StrongCopyleftLT, "LOW", "GPL is generally safe for internal SaaS use, but verify no code is shipped to clients (e.g., JS).", []string{"Verify no client-side exposure"}},
		}
	default: // internal
		return []Rule{
			{NetworkCopyleftLT, "MEDIUM", "Internal use of AGPL is usually safe, but check for potential future exposure.", []string{"Monitor usage"}},
			{StrongCopyleftLT, "LOW", "GPL is safe for internal use.", []string{"No action required"}},
		}
	}
}
