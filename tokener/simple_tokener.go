package tokener

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
)

type simpleTokener struct {
	spliterSet map[byte]bool
	tokens     []Token
	offset     int
	rules      []Rule
	reges      []*regexp.Regexp
}

// NewSimpleTokener x
func NewSimpleTokener(spliterChars string, rules ...Rule) (Tokener, error) {
	rs := make([]Rule, 0, len(rules))
	rs = append(rs, rules...)
	reges := make([]*regexp.Regexp, 0, len(rules))
	for _, r := range rs {
		e, err := regexp.Compile(r.RegExp())
		if err != nil {
			return nil, fmt.Errorf("Compile Regexp: %v", err)
		}
		e.Longest()
		reges = append(reges, e)
	}

	spSet := make(map[byte]bool, len(spliterChars))
	for _, r := range spliterChars {
		spSet[byte(r)] = true
	}

	return &simpleTokener{
		spliterSet: spSet,
		tokens:     make([]Token, 0, 1024),
		rules:      rs,
		reges:      reges,
	}, nil
}

func (st *simpleTokener) Seek(offset, whence int) error {
	old := st.offset
	if whence == 0 {
		st.offset = offset
	} else if whence == 1 {
		st.offset = offset + st.offset
	} else if whence == 2 {
		st.offset = len(st.tokens) - offset
	} else {
		return fmt.Errorf("Invalid whence argument: %v", whence)
	}

	if st.offset < 0 || st.offset > len(st.tokens) {
		st.offset = old
		return fmt.Errorf("Bad offset, the offset should between [0, %v] and your offset is: %v", len(st.tokens), st.offset)
	}
	return nil
}

func (st *simpleTokener) Offset() int {
	return st.offset
}

func (st *simpleTokener) HasNext() bool {
	return st.offset < len(st.tokens)
}

func (st *simpleTokener) Next() (Token, error) {
	if len(st.tokens) == st.offset {
		return nil, io.EOF
	}

	st.offset++
	return st.tokens[st.offset-1], nil
}

func (st *simpleTokener) Peekth(kth int) (Token, error) {
	if len(st.tokens) < st.offset+kth+1 {
		return nil, io.EOF
	}
	return st.tokens[st.offset+kth], nil
}

func (st *simpleTokener) Tokenize(input io.Reader) error {
	data, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	return st.tokenize(data)
}

func (st *simpleTokener) tokenize(data []byte) error {
	for len(data) > 0 {
		// remove leading spliter chars
		for len(data) > 0 && st.spliterSet[data[0]] {
			data = data[1:]
		}

		if len(data) == 0 { // tokenization success and finish
			return nil
		}

		// find the longest matched regexp
		longest := 0
		typex := 0
		for i, r := range st.reges {
			pos := r.FindIndex(data)
			if pos != nil && pos[0] == 0 && // matched from the begining
				pos[1] > longest { // and is longger
				longest = pos[1]
				typex = st.rules[i].Type()
			}
		}

		if longest == 0 { // no matched regexp, return an error
			return fmt.Errorf("Tokenize error begining at >>>: %v", data)
		}

		st.tokens = append(st.tokens, &simpleToken{
			token: string(data[:longest]),
			typeX: typex,
		})

		data = data[longest:]
	}
	return nil
}
