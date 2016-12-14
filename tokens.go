package func_lang

import "github.com/qw4990/func_lang/tokener"

const (
	_ASSIGN_TYPE   = iota
	_ASSIGN_REGEX  = "="
	_DECLARE_TYPE  = iota
	_DECLARE_REGEX = ":="

	_IF_TYPE   = iota
	_IF_REGEX  = "if"
	_FOR_TYPE  = iota
	_FOR_REGEX = "for"

	_LEFT_BRACKET_TYPE   = iota
	_LEFT_BRACKET_REGEX  = "[(]"
	_RIGHT_BRACKET_TYPE  = iota
	_RIGHT_BRACKET_REGEX = "[)]"
	_COMMA_TYPE          = iota
	_COMMA_REGEX         = ","
	_LEFT_BRACE_TYPE     = iota
	_LEFT_BRACE_REGEX    = "{"
	_RIGHT_BRACE_TYPE    = iota
	_RIGHT_BRACE_REGEX   = "}"

	_RETURN_TYPE  = iota
	_RETURN_REGEX = "return"
	_NUMBER_TYPE  = iota
	_NUMBER_REGEX = "[0-9]+([.][0-9]+)?"
	_STRING_TYPE  = iota
	_STRING_REGEX = "[\"'][^\"]*[\"']"
	_IDENT_TYPE   = iota
	_IDENT_REGEX  = "[A-Za-z_][A-Za-z_0-9]*"

	_SPLITER_CHARS = "\r\n\t; "
)

type tokenizeRule struct {
	typex  int
	regexp string
}

func (tr *tokenizeRule) Type() int {
	return tr.typex
}

func (tr *tokenizeRule) RegExp() string {
	return tr.regexp
}

// the order of these rules is important !!!
var tokenizeRules = []tokener.Rule{
	&tokenizeRule{_ASSIGN_TYPE, _ASSIGN_REGEX},
	&tokenizeRule{_DECLARE_TYPE, _DECLARE_REGEX},

	&tokenizeRule{_IF_TYPE, _IF_REGEX},
	&tokenizeRule{_FOR_TYPE, _FOR_REGEX},

	&tokenizeRule{_LEFT_BRACKET_TYPE, _LEFT_BRACKET_REGEX},
	&tokenizeRule{_RIGHT_BRACKET_TYPE, _RIGHT_BRACKET_REGEX},
	&tokenizeRule{_COMMA_TYPE, _COMMA_REGEX},
	&tokenizeRule{_LEFT_BRACE_TYPE, _LEFT_BRACE_REGEX},
	&tokenizeRule{_RIGHT_BRACE_TYPE, _RIGHT_BRACE_REGEX},

	&tokenizeRule{_RETURN_TYPE, _RETURN_REGEX},
	&tokenizeRule{_NUMBER_TYPE, _NUMBER_REGEX},
	&tokenizeRule{_STRING_TYPE, _STRING_REGEX},
	&tokenizeRule{_IDENT_TYPE, _IDENT_REGEX},
}
