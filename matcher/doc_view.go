package matcher

type DocView interface {
	Get(path string) (any, bool, error)
}
