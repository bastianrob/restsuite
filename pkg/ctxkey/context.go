package ctxkey

//Context scope collection
const (
	OrganizationName = ContextKey("organization_name")
)

//ContextKey alias
type ContextKey string

func (c ContextKey) String() string {
	return "restify_context_key_" + string(c)
}
