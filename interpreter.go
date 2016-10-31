package ez_lang

import (
	"errors"
	"fmt"
	"io"
	"lab/ez_lang/tokener"
	"strconv"
)

var errUnexpectedToken = errors.New("Unexpected token")

// EZInterpreter x
type ezInterpreter struct {
	tokener     tokener.Tokener
	globalScope *ezScope
	nowScope    *ezScope
}

// NewEZInterpreter x
func NewEZInterpreter() EZInterpreter {
	tokener, _ := tokener.NewSimpleTokener(_EZ_SPLITER_CHARS, ezTokenizeRules...)
	global := newEZScope(nil)
	i := &ezInterpreter{
		globalScope: global,
		nowScope:    global,
		tokener:     tokener,
	}

	RegisterDefaultFuncs(i)
	return i
}

// RegisterFunc x
func (ezi *ezInterpreter) RegisterFunc(funcName string, ezFunc EZFunc) {
	ezi.globalScope.funcTable[funcName] = ezFunc
}

// Interprete x
func (ezi *ezInterpreter) Interprete(input io.Reader) error {
	if err := ezi.tokener.Tokenize(input); err != nil {
		return err
	}

	err := ezi.interpreteProgram()
	if err != nil {
		i := 0
		errMsg := "Interprete Err at >>>: "
		for ezi.tokener.HasNext() && i < 5 {
			tk, err := ezi.tokener.Next()
			if err != nil {
				break
			}
			errMsg = errMsg + tk.Token() + " "
			i++
		}

		err = errors.New(errMsg + " " + err.Error())
	}
	return err
}

func (ezi *ezInterpreter) interpreteProgram() error {
	for ezi.tokener.HasNext() {
		if err := ezi.interpreteStatement(); err != nil {
			return err
		}
	}
	return nil
}

func (ezi *ezInterpreter) interpreteStatement() error {
	nToken, err := ezi.tokener.Peekth(0)
	if err != nil {
		return err
	}

	if nToken.Type() == _IF_TYPE {
		return ezi.interpreteIf()
	} else if nToken.Type() == _FOR_TYPE {
		return ezi.interpreteFor()
	}

	nnToken, err := ezi.tokener.Peekth(1)
	if err == nil {
		if nToken.Type() == _IDENT_TYPE && nnToken.Type() == _DECLARE_TYPE {
			return ezi.interpreteDeclare()
		} else if nToken.Type() == _IDENT_TYPE && nnToken.Type() == _ASSIGN_TYPE {
			return ezi.interpreteAssign()
		}
	}

	_, err = ezi.interpreteExpression()
	return err
}

func (ezi *ezInterpreter) interpreteIf() error {
	ezi.tokener.Next() // remove "if"

	// check the expression
	result, err := ezi.interpreteExpression()
	if err != nil {
		return err
	}
	if result.isTrue() == false {
		return ezi.skipCodeBlock()
	}

	// the new scope
	nowScope := ezi.nowScope
	ezi.nowScope = newEZScope(nowScope)
	defer func() {
		ezi.nowScope = nowScope
	}()

	leftBrace, err := ezi.tokener.Next()
	if err != nil {
		return err
	}
	if leftBrace.Type() != _LEFT_BRACE_TYPE {
		return errUnexpectedToken
	}

	err = ezi.interpreteProgram()  // interprete recursively
	if err != errUnexpectedToken { // there should be a right brace which cause this err
		return err
	}

	rightBrace, err := ezi.tokener.Next()
	if err != nil {
		return err
	}
	if rightBrace.Type() != _RIGHT_BRACE_TYPE {
		return errUnexpectedToken
	}

	return nil
}

func (ezi *ezInterpreter) interpreteFor() error {
	forBegin := ezi.tokener.Offset()
	for {
		ezi.tokener.Next() // remote "for"

		// check the expression
		result, err := ezi.interpreteExpression()
		if err != nil {
			return err
		}
		if result.isTrue() == false {
			return ezi.skipCodeBlock() // end the loop
		}

		// run the loop's body; create a new scope first
		nowScope := ezi.nowScope
		ezi.nowScope = newEZScope(nowScope)
		defer func() {
			ezi.nowScope = nowScope
		}()

		leftBrace, err := ezi.tokener.Next()
		if err != nil {
			return err
		}
		if leftBrace.Type() != _LEFT_BRACE_TYPE {
			return errUnexpectedToken
		}

		err = ezi.interpreteProgram()  // interprete recursively
		if err != errUnexpectedToken { // there should be a right brace which cause this err
			return err
		}

		rightBrace, err := ezi.tokener.Next()
		if err != nil {
			return err
		}
		if rightBrace.Type() != _RIGHT_BRACE_TYPE {
			return errUnexpectedToken
		}

		err = ezi.tokener.Seek(forBegin, 0) // seek the tokener and run again
		if err != nil {
			return nil
		}
	}
}

func (ezi *ezInterpreter) interpreteDeclare() error {
	// ident := expr
	identTk, err := ezi.tokener.Next()
	if err != nil {
		return err
	}

	ezi.tokener.Next() // skip ":="

	result, err := ezi.interpreteExpression()
	if err != nil {
		return err
	}
	ezi.nowScope.varTable[identTk.Token()] = result
	return nil
}

func (ezi *ezInterpreter) interpreteAssign() error {
	// idnet = expr
	identTk, err := ezi.tokener.Next()
	if err != nil {
		return err
	}

	ident := ezi.nowScope.lookupVar(identTk.Token())
	if ident == nil {
		return fmt.Errorf("No such a variable: %v", identTk.Token())
	}

	ezi.tokener.Next() // skip "="

	result, err := ezi.interpreteExpression()
	if err != nil {
		return err
	}
	ident.clone(result)
	return nil
}

func (ezi *ezInterpreter) interpreteExpression() (*ezVar, error) {
	nToken, err := ezi.tokener.Peekth(0)
	if err != nil {
		return nil, err
	}

	if nToken.Type() == _NUMBER_TYPE { // 12.34
		return ezi.interpreteNumberExpr()
	} else if nToken.Type() == _STRING_TYPE { // "hello :)"
		return ezi.interpreteStringExpr()
	}

	if nToken.Type() == _IDENT_TYPE {
		nnToken, err := ezi.tokener.Peekth(1)
		if err == nil && nnToken.Type() == _LEFT_BRACKET_TYPE { // print(123)
			return ezi.interpreteFuncCallExpr()
		}
		return ezi.interpreteIdentExpr() // someVar
	}

	return nil, errUnexpectedToken
}

func (ezi *ezInterpreter) interpreteIdentExpr() (*ezVar, error) {
	identTk, err := ezi.tokener.Next()
	if err != nil {
		return nil, err
	}
	ident := ezi.nowScope.lookupVar(identTk.Token())
	if ident == nil {
		return nil, fmt.Errorf("No such a variable: %v", identTk.Token())
	}
	return ident, nil
}

func (ezi *ezInterpreter) interpreteFuncCallExpr() (*ezVar, error) {
	funcTk, err := ezi.tokener.Next()
	if err != nil {
		return nil, err
	}
	fun := ezi.nowScope.lookupFunc(funcTk.Token())
	if fun == nil {
		return nil, fmt.Errorf("No such a function: %v", funcTk.Token())
	}

	leftBracket, err := ezi.tokener.Next()
	if err != nil {
		return nil, err
	}
	if leftBracket.Type() != _LEFT_BRACKET_TYPE {
		return nil, errUnexpectedToken
	}

	var args []EZVar
	for {
		arg, err := ezi.interpreteExpression()
		if err == nil {
			args = append(args, arg)
		} else if err == errUnexpectedToken { // ',' or ')'
			break
		}
		if err != nil {
			return nil, err
		}

		commaOrBracket, err := ezi.tokener.Next()
		if err != nil {
			return nil, err
		}
		if commaOrBracket.Type() == _COMMA_TYPE {
			// continue
		} else if commaOrBracket.Type() == _RIGHT_BRACKET_TYPE {
			break
		} else {
			return nil, errUnexpectedToken
		}
	}

	result, err := ezi.safelyRun(funcTk.Token(), fun, args...)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return ezNil, nil
	}

	return cloneEZVar(result), nil
}

func (ezi *ezInterpreter) interpreteNumberExpr() (*ezVar, error) {
	numTk, err := ezi.tokener.Next()
	if err != nil {
		return nil, err
	}

	num, err := strconv.ParseFloat(numTk.Token(), 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid number format: %v", numTk.Token())
	}
	return newEZNum(num), nil
}

func (ezi *ezInterpreter) interpreteStringExpr() (*ezVar, error) {
	strTk, err := ezi.tokener.Next()
	if err != nil {
		return nil, err
	}

	str := strTk.Token()
	if len(str) < 2 { // with quotes
		return nil, fmt.Errorf("Internal err: string has no quotes: %v", str)
	}

	return newEZStr(str[1 : len(str)-1]), nil
}

// skipCodeBlock skips code block included by a pair of braces
func (ezi *ezInterpreter) skipCodeBlock() error {
	leftBrace, err := ezi.tokener.Next()
	if err != nil {
		return err
	}
	if leftBrace.Type() != _LEFT_BRACE_TYPE {
		return errUnexpectedToken
	}

	pairCnt := 1
	for {
		token, err := ezi.tokener.Next()
		if err != nil {
			return err
		}

		if token.Type() == _LEFT_BRACE_TYPE {
			pairCnt++
		} else if token.Type() == _RIGHT_BRACE_TYPE {
			pairCnt--
		}

		if pairCnt == 0 {
			break
		}
	}

	return nil
}

func (ezi *ezInterpreter) safelyRun(funName string, fun EZFunc, args ...EZVar) (result EZVar, err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("Panic when run func: %v, err: %v", funName, pErr)
			result = nil
		}
	}()

	result = fun(args...)
	return
}
