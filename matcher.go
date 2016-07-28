package mongofil

type Matcher interface {
	Match(doc []byte) bool
}
