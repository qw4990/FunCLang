package func_lang

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/qw4990/func_lang/tokener"
)

var errUnexpectedToken = errors.New("Unexpected token")

type funcaller struct {
	tokener     tokener.Tokener
	globalScope *fcScope
	nowScope    *fcScope

	result   Var
	returned bool
}

// NewFunCaller ...
func NewFunCaller() FunCaller {
	tokener, _ := tokener.NewSimpleTokener(_SPLITER_CHARS, tokenizeRules...)
	global := newFCScope(nil)
	i := &funcaller{
		globalScope: global,
		nowScope:    global,
		tokener:     tokener,
	}

	RegisterDefaultFuncs(i)
	return i
}

// RegisterFunc x
func (fci *funcaller) RegisterFunc(funcName string, f Func) {
	fci.globalScope.funcTable[funcName] = f
}

// Call x
func (fci *funcaller) Call(funcBody io.Reader) (Var, error) {
	if err := fci.tokener.Tokenize(funcBody); err != nil {
		return nil, err
	}

	err := fci.interpreteProgram()
	if err != nil {
		i := 0
		errMsg := "Interprete Err at >>>: "
		for fci.tokener.HasNext() && i < 5 {
			tk, err := fci.tokener.Next()
			if err != nil {
				break
			}
			errMsg = errMsg + tk.Token() + " "
			i++
		}

		err = errors.New(errMsg + " " + err.Error())
		return nil, err
	}
	return fci.result, nil
}

func (fci *funcaller) interpreteProgram() error {
	for fci.tokener.HasNext() {
		if err := fci.interpreteStatement(); err != nil {
			return err
		}
		if fci.returned == true {
			return nil
		}
	}
	return nil
}

func (fci *funcaller) interpreteStatement() error {
	nToken, err := fci.tokener.Peekth(0)
	if err != nil {
		return err
	}

	if nToken.Type() == _RETURN_TYPE {
		return fci.interpreteReturn()
	} else if nToken.Type() == _IF_TYPE {
		return fci.interpreteIf()
	} else if nToken.Type() == _FOR_TYPE {
		return fci.interpreteFor()
	}

	nnToken, err := fci.tokener.Peekth(1)
	if err == nil {
		if nToken.Type() == _IDENT_TYPE && nnToken.Type() == _DECLARE_TYPE {
			return fci.interpreteDeclare()
		} else if nToken.Type() == _IDENT_TYPE && nnToken.Type() == _ASSIGN_TYPE {
			return fci.interpreteAssign()
		}
	}

	_, err = fci.interpreteExpression()
	return err
}

func (fci *funcaller) interpreteReturn() error {
	fci.tokener.Next() // remove "return"
	result, err := fci.interpreteExpression()
	if err != nil {
		return err
	}

	fci.returned = true
	fci.result = result
	return nil
}

func (fci *funcaller) interpreteIf() error {
	fci.tokener.Next() // remove "if"

	// check the expression
	result, err := fci.interpreteExpression()
	if err != nil {
		return err
	}
	if result.IsTrue() == false {
		return fci.skipCodeBlock()
	}

	// the new scope
	nowScope := fci.nowScope
	fci.nowScope = newFCScope(nowScope)
	defer func() {
		fci.nowScope = nowScope
	}()

	leftBrace, err := fci.tokener.Next()
	if err != nil {
		return err
	}
	if leftBrace.Type() != _LEFT_BRACE_TYPE {
		return errUnexpectedToken
	}

	err = fci.interpreteProgram()  // interprete recursively
	if err != errUnexpectedToken { // there should be a right brace which cause this err
		return err
	}
	if fci.returned == true {
		return nil
	}

	rightBrace, err := fci.tokener.Next()
	if err != nil {
		return err
	}
	if rightBrace.Type() != _RIGHT_BRACE_TYPE {
		return errUnexpectedToken
	}

	return nil
}

func (fci *funcaller) interpreteFor() error {
	forBegin := fci.tokener.Offset()
	for {
		fci.tokener.Next() // remote "for"

		// check the expression
		result, err := fci.interpreteExpression()
		if err != nil {
			return err
		}
		if result.IsTrue() == false {
			return fci.skipCodeBlock() // end the loop
		}

		// run the loop's body; create a new scope first
		nowScope := fci.nowScope
		fci.nowScope = newFCScope(nowScope)
		defer func() {
			fci.nowScope = nowScope
		}()

		leftBrace, err := fci.tokener.Next()
		if err != nil {
			return err
		}
		if leftBrace.Type() != _LEFT_BRACE_TYPE {
			return errUnexpectedToken
		}

		err = fci.interpreteProgram()  // interprete recursively
		if err != errUnexpectedToken { // there should be a right brace which cause this err
			return err
		}
		if fci.returned == true {
			return nil
		}

		rightBrace, err := fci.tokener.Next()
		if err != nil {
			return err
		}
		if rightBrace.Type() != _RIGHT_BRACE_TYPE {
			return errUnexpectedToken
		}

		err = fci.tokener.Seek(forBegin, 0) // seek the tokener and run again
		if err != nil {
			return nil
		}
	}
}

func (fci *funcaller) interpreteDeclare() error {
	// ident := expr
	identTk, err := fci.tokener.Next()
	if err != nil {
		return err
	}

	fci.tokener.Next() // skip ":="

	result, err := fci.interpreteExpression()
	if err != nil {
		return err
	}
	fci.nowScope.varTable[identTk.Token()] = result
	return nil
}

func (fci *funcaller) interpreteAssign() error {
	// idnet = expr
	identTk, err := fci.tokener.Next()
	if err != nil {
		return err
	}

	ident := fci.nowScope.lookupVar(identTk.Token())
	if ident == nil {
		return fmt.Errorf("No such a variable: %v", identTk.Token())
	}

	fci.tokener.Next() // skip "="

	result, err := fci.interpreteExpression()
	if err != nil {
		return err
	}
	ident.clone(result)
	return nil
}

func (fci *funcaller) interpreteExpression() (*fcVar, error) {
	nToken, err := fci.tokener.Peekth(0)
	if err != nil {
		return nil, err
	}

	if nToken.Type() == _NUMBER_TYPE { // 12.34
		return fci.interpreteNumberExpr()
	} else if nToken.Type() == _STRING_TYPE { // "hello :)"
		return fci.interpreteStringExpr()
	}

	if nToken.Type() == _IDENT_TYPE {
		nnToken, err := fci.tokener.Peekth(1)
		if err == nil && nnToken.Type() == _LEFT_BRACKET_TYPE { // Println(123)
			return fci.interpreteFuncCallExpr()
		}
		return fci.interpreteIdentExpr() // someVar
	}

	return nil, errUnexpectedToken
}

func (fci *funcaller) interpreteIdentExpr() (*fcVar, error) {
	identTk, err := fci.tokener.Next()
	if err != nil {
		return nil, err
	}
	ident := fci.nowScope.lookupVar(identTk.Token())
	if ident == nil {
		return nil, fmt.Errorf("No such a variable: %v", identTk.Token())
	}
	return ident, nil
}

func (fci *funcaller) interpreteFuncCallExpr() (*fcVar, error) {
	funcTk, err := fci.tokener.Next()
	if err != nil {
		return nil, err
	}
	fun := fci.nowScope.lookupFunc(funcTk.Token())
	if fun == nil {
		return nil, fmt.Errorf("No such a function: %v", funcTk.Token())
	}

	leftBracket, err := fci.tokener.Next()
	if err != nil {
		return nil, err
	}
	if leftBracket.Type() != _LEFT_BRACKET_TYPE {
		return nil, errUnexpectedToken
	}

	var args []Var
	nToken, err := fci.tokener.Peekth(0)
	if err != nil {
		return nil, err
	}
	if nToken.Type() == _RIGHT_BRACKET_TYPE {
		fci.tokener.Next() // empty args
	} else {
		for {
			arg, err := fci.interpreteExpression()
			if err == nil {
				args = append(args, arg)
			} else {
				return nil, err
			}

			commaOrBracket, err := fci.tokener.Next()
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
	}

	result, err := fci.safelyRun(funcTk.Token(), fun, args...)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return fcNil, nil
	}

	return cloneVar(result), nil
}

func (fci *funcaller) interpreteNumberExpr() (*fcVar, error) {
	numTk, err := fci.tokener.Next()
	if err != nil {
		return nil, err
	}

	num, err := strconv.ParseFloat(numTk.Token(), 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid number format: %v", numTk.Token())
	}
	return newFCNum(num), nil
}

func (fci *funcaller) interpreteStringExpr() (*fcVar, error) {
	strTk, err := fci.tokener.Next()
	if err != nil {
		return nil, err
	}

	str := strTk.Token()
	if len(str) < 2 { // with quotes
		return nil, fmt.Errorf("Internal err: string has no quotes: %v", str)
	}

	return newFCStr(str[1 : len(str)-1]), nil
}

// skipCodeBlock skips code block included by a pair of braces
func (fci *funcaller) skipCodeBlock() error {
	leftBrace, err := fci.tokener.Next()
	if err != nil {
		return err
	}
	if leftBrace.Type() != _LEFT_BRACE_TYPE {
		return errUnexpectedToken
	}

	pairCnt := 1
	for {
		token, err := fci.tokener.Next()
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

func (fci *funcaller) safelyRun(funName string, fun Func, args ...Var) (result Var, err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			err = fmt.Errorf("Panic when run func: %v, err: %v", funName, pErr)
			result = nil
		}
	}()

	result = fun(args...)
	return
}
