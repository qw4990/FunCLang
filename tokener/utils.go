package tokener

type simpleToken struct {
	token string
	typeX int
}

func (st *simpleToken) Token() string {
	return st.token
}

func (st *simpleToken) Type() int {
	return st.typeX
}
