package parser

import (
	"fmt"
	"github.com/xianlianghe0123/jsonpath/internal/ast"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

const (
	dot                = '.'
	star               = '*'
	dollar             = '$'
	comma              = ','
	colon              = ':'
	add                = '+'
	sub                = '-'
	leftSquareBracket  = '['
	rightSquareBracket = ']'
	leftParenthesis    = '('
	rightParenthesis   = ')'
	singleQuotes       = '\''
	doubleQuotes       = '"'
)

var stateMachine = map[tokenType][]tokenType{
	TokenStart:     {TokenRoot, TokenAll, TokenField},
	TokenRoot:      {TokenDot, TokenRecursion, TokenSquare},
	TokenAll:       {TokenDot, TokenRecursion, TokenSquare},
	TokenDot:       {TokenAll, TokenField, TokenSquare},
	TokenRecursion: {TokenAll, TokenField, TokenSquare},
	TokenField:     {TokenDot, TokenRecursion, TokenSquare},
	TokenSquare:    {TokenDot, TokenRecursion},
}

func transfer(from, to tokenType) bool {
	for _, t := range stateMachine[from] {
		if t == to {
			return true
		}
	}
	return false
}

type Parser struct {
	input  []rune
	offset int
	status tokenType
	ast    *ast.AST
	err    error
	once   sync.Once
}

func NewParser(jsonPath string) *Parser {
	return &Parser{
		input:  []rune(jsonPath),
		offset: 0,
		status: TokenStart,
		ast:    nil,
		err:    nil,
		once:   sync.Once{},
	}
}

func (p *Parser) Parse() (*ast.AST, error) {
	p.once.Do(func() {
		n, err := p.parse()
		if err != nil {
			p.err = err
			return
		}
		p.ast = ast.NewAST(n)
	})
	return p.ast, p.err
}

func (p *Parser) parse() (ast.Node, error) {
	if p.offset == len(p.input) {
		return ast.NewEnd(), nil
	}
	switch p.input[p.offset] {
	case dot:
		t, err := p.scanDots()
		if err != nil {
			return nil, err
		}
		if !transfer(p.status, t.TokenType) {
			return nil, fmt.Errorf(`syntax error: unexpected token %s`, t.Value)
		}
		p.status = t.TokenType
		n, err := p.parse()
		if err != nil {
			return nil, err
		}
		if t.TokenType == TokenRecursion {
			return ast.NewRecursion(n), nil
		}
		return n, nil
	case leftSquareBracket:
		if !transfer(p.status, TokenSquare) {
			return nil, fmt.Errorf(`syntax error: unexpected token %s`, string(leftSquareBracket))
		}
		p.status = TokenSquare
		node, err := p.parseSquare()
		if err != nil {
			return nil, err
		}
		return node, err
	default:
		t, err := p.scanField()
		if err != nil {
			return nil, err
		}
		if !transfer(p.status, t.TokenType) {
			return nil, fmt.Errorf(`syntax error: unexpected token %s`, t.Value)
		}
		p.status = t.TokenType
		n, err := p.parse()
		if err != nil {
			return nil, err
		}
		switch t.TokenType {
		case TokenAll:
			return ast.NewAll(n), nil
		case TokenField:
			return ast.NewSingleField(t.Value, n), nil
		case TokenRoot:
			return ast.NewRoot(n), nil
		}
		return nil, fmt.Errorf("parser error")
	}
}

func (p *Parser) pop(offset int) string {
	t := p.offset
	p.offset = offset
	return string(p.input[t:offset])
}

func (p *Parser) scanDots() (*Token, error) {
	off := p.offset + 1
	for ; off < len(p.input) && p.input[off] == dot; off++ {
	}
	switch off - p.offset {
	case 1:
		return &Token{
			TokenType: TokenDot,
			Value:     p.pop(off),
		}, nil
	case 2:
		return &Token{
			TokenType: TokenRecursion,
			Value:     p.pop(off),
		}, nil
	default:
		return nil, fmt.Errorf(`syntax error near %q: unexpected token %s`, string(p.input[p.offset:]), string(p.input[p.offset:off]))
	}
}

func (p *Parser) scanField() (*Token, error) {
	i := p.offset + 1
	for ; i < len(p.input); i++ {
		if p.input[i] == dot ||
			p.input[i] == leftSquareBracket {
			break
		}
	}
	tokenType := TokenField
	value := p.pop(i)
	switch value {
	case "$":
		tokenType = TokenRoot
	case "*":
		tokenType = TokenAll
	}
	return &Token{
		TokenType: tokenType,
		Value:     value,
	}, nil
}

func (p *Parser) parseSquare() (ast.Node, error) {
	p.offset++
	if p.offset >= len(p.input) {
		return nil, fmt.Errorf("syntax err near %s", string(p.input[p.offset:]))
	}
	switch p.input[p.offset] {
	case singleQuotes, doubleQuotes:
		node, err := p.parseFields()
		if err != nil {
			return nil, err
		}
		return node, nil
	case colon, sub, '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		node, err := p.parseIndexes()
		if err != nil {
			return nil, err
		}
		return node, nil
	case star:
		p.offset++
		p.skipSpace()
		if p.offset == len(p.input) || p.input[p.offset] != rightSquareBracket {
			return nil, fmt.Errorf("syntax err near %s", string(p.input[p.offset:]))
		}
		p.offset++
		n, err := p.parse()
		if err != nil {
			return nil, err
		}
		return ast.NewAll(n), nil
	default:
		return nil, fmt.Errorf("syntax err near %s: expected string or integer", string(p.input[p.offset:]))
	}
}

func (p *Parser) skipSpace() {
	for ; p.offset < len(p.input) && unicode.IsSpace(p.input[p.offset]); p.offset++ {
	}
}

func (p *Parser) parseFields() (ast.Node, error) {
	fields := make([]string, 0, 1)
	for {
		str, err := p.scanString()
		if err != nil {
			return nil, err
		}
		fields = append(fields, str)
		if p.offset == len(p.input) {
			return nil, fmt.Errorf("syntax err near %s: could not found ]", string(p.input[p.offset:]))
		}
		p.skipSpace()
		if p.input[p.offset] == rightSquareBracket {
			p.offset++
			break
		}
		if p.input[p.offset] != comma {
			return nil, fmt.Errorf("syntax err near %s", string(p.input[p.offset:]))
		}
		p.offset++
		p.skipSpace()
	}
	n, err := p.parse()
	if err != nil {
		return nil, err
	}
	if len(fields) > 1 {
		return ast.NewMultiFields(fields, n), nil
	}
	return ast.NewSingleField(fields[0], n), nil
}

func (p *Parser) scanString() (string, error) {
	if p.input[p.offset] != singleQuotes && p.input[p.offset] != doubleQuotes {
		return "", fmt.Errorf(`syntax error near %q: could not find quotes`, string(p.input[p.offset:]))
	}
	result := make([]rune, 0)
	i := p.offset + 1
	for ; i < len(p.input) && p.input[i] != p.input[p.offset]; i++ {
		if p.input[i] == '\\' {
			i++
		}
		result = append(result, p.input[i])
	}
	if i == len(p.input) {
		return "", fmt.Errorf(`syntax error near %q: unmatched quotes`, string(p.input[p.offset:]))
	}
	p.offset = i + 1
	return string(result), nil
}

type indexesData struct {
	isSlice bool
	index   int
	start   int
	end     int
	step    int
}

func (p *Parser) parseIndexes() (ast.Node, error) {
	data := make([]*indexesData, 0, 1)
	for slice := make([]int, 0, 3); ; {
		p.skipSpace()
		integer := 0
		if !strings.ContainsRune(":,]", p.input[p.offset]) {
			var err error
			integer, err = p.scanInteger()
			if err != nil {
				return nil, err
			}
		}
		if len(slice) == 3 {
			return nil, fmt.Errorf(`syntax error near %q`, string(p.input[p.offset:]))
		}
		slice = append(slice, integer)
		if p.offset == len(p.input) {
			return nil, fmt.Errorf("syntax err near %q : could not found ]", string(p.input[p.offset:]))
		}
		p.skipSpace()
		p.offset++
		switch p.input[p.offset-1] {
		case colon:
		case comma, rightSquareBracket:
			if len(slice) == 1 {
				data = append(data, &indexesData{isSlice: false, index: slice[0]})
			} else {
				end, step := 0, 0
				if len(slice) > 1 {
					end = slice[1]
					if len(slice) > 2 {
						step = slice[2]
					}
				}
				data = append(data, &indexesData{isSlice: true, start: slice[0], end: end, step: step})
			}
			slice = slice[:0]
			if p.input[p.offset-1] == rightSquareBracket {
				goto exit
			}
		default:
			return nil, fmt.Errorf("syntax err near %s", string(p.input[p.offset:]))
		}
	}
exit:
	n, err := p.parse()
	if err != nil {
		return nil, err
	}
	nodes := make([]ast.Node, 0, len(data))
	for _, d := range data {
		if !d.isSlice {
			nodes = append(nodes, ast.NewIndexField(d.index, n))
		} else {
			nodes = append(nodes, ast.NewSlice(d.start, d.end, d.step, n))
		}
	}
	if len(nodes) == 1 {
		return nodes[0], nil
	}
	return ast.NewIndexes(nodes), nil
}

func (p *Parser) scanInteger() (int, error) {
	off := p.offset + 1
	for ; off < len(p.input) && unicode.IsNumber(p.input[off]); off++ {
	}
	integer, err := strconv.Atoi(string(p.input[p.offset:off]))
	if err != nil {
		return 0, fmt.Errorf(`syntax error near %q: could not parse integer`, string(p.input[p.offset:]))
	}
	p.offset = off
	return integer, nil
}
