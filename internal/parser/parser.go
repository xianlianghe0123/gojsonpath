package parser

import (
	"fmt"
	"github.com/xianlianghe0123/jsonpath/internal/ast"
	"strconv"
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

type Parser struct {
	input  []rune
	offset int
	ast    ast.AST
	err    error
	once   sync.Once
}

func NewParser(jsonPath string) *Parser {
	return &Parser{
		input: []rune(jsonPath),
		once:  sync.Once{},
	}
}

func (s *Parser) Parse() (ast.AST, error) {
	if s.err != nil {
		return nil, s.err
	}
	if len(s.ast) == 0 {
		s.err = s.parse()
		if s.err != nil {
			return nil, s.err
		}
	}
	return s.ast, nil
}

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

func (s *Parser) parse() error {
	current := TokenStart
	for s.offset < len(s.input) {
		switch s.input[s.offset] {
		case dot:
			t, err := s.scanDots()
			if err != nil {
				return err
			}
			if !transfer(current, t.TokenType) {
				return fmt.Errorf(`syntax error: unexpected token %s`, t.Value)
			}
			if t.TokenType == TokenRecursion {
				s.ast = append(s.ast, ast.NewRecursion())
			}
			current = t.TokenType
		case leftSquareBracket:
			if !transfer(current, TokenSquare) {
				return fmt.Errorf(`syntax error: unexpected token %s`, string(leftSquareBracket))
			}
			err := s.parseSquare()
			if err != nil {
				return err
			}
			current = TokenSquare
		default:
			t, err := s.scanField()
			if err != nil {
				return err
			}
			if !transfer(current, t.TokenType) {
				return fmt.Errorf(`syntax error: unexpected token %s`, t.Value)
			}
			switch t.TokenType {
			case TokenAll:
				s.ast = append(s.ast, ast.NewAll())
			case TokenField:
				s.ast = append(s.ast, ast.NewSingleField(t.Value))
			case TokenRoot:
				s.ast = append(s.ast, ast.NewRoot())
			}
			current = t.TokenType
		}
	}
	return nil
}

func (s *Parser) pop(offset int) string {
	t := s.offset
	s.offset = offset
	return string(s.input[t:offset])
}

func (s *Parser) scanDots() (*Token, error) {
	off := s.offset + 1
	for ; off < len(s.input) && s.input[off] == dot; off++ {
	}
	switch off - s.offset {
	case 1:
		return &Token{
			TokenType: TokenDot,
			Value:     s.pop(off),
		}, nil
	case 2:
		return &Token{
			TokenType: TokenRecursion,
			Value:     s.pop(off),
		}, nil
	default:
		return nil, fmt.Errorf(`syntax error near %q: unexpected token %s`, string(s.input[s.offset:]), string(s.input[s.offset:off]))
	}
}

func (s *Parser) scanField() (*Token, error) {
	i := s.offset + 1
	for ; i < len(s.input); i++ {
		if s.input[i] == dot ||
			s.input[i] == leftSquareBracket {
			break
		}
	}
	tokenType := TokenField
	value := s.pop(i)
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

func (s *Parser) parseSquare() error {
	s.offset++
	if s.offset >= len(s.input) {
		return fmt.Errorf("syntax err near %s", string(s.input[s.offset:]))
	}
	switch s.input[s.offset] {
	case singleQuotes, doubleQuotes:
		node, err := s.parseFields()
		if err != nil {
			return err
		}
		s.ast = append(s.ast, node)
	case colon, sub, '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		node, err := s.parseIndexes()
		if err != nil {
			return err
		}
		s.ast = append(s.ast, node)
	case star:
		s.offset++
		s.skipSpace()
		if s.offset == len(s.input) || s.input[s.offset] != rightSquareBracket {
			return fmt.Errorf("syntax err near %s", string(s.input[s.offset:]))
		}
		s.ast = append(s.ast, ast.NewAll())
		s.offset++
	default:
		return fmt.Errorf("syntax err near %s: expected string or integer", string(s.input[s.offset:]))
	}
	return nil
}

func (s *Parser) skipSpace() {
	for ; s.offset < len(s.input) && unicode.IsSpace(s.input[s.offset]); s.offset++ {
	}
}

func (s *Parser) parseFields() (ast.Node, error) {
	fields := make([]string, 0, 1)
	for {
		str, err := s.scanString()
		if err != nil {
			return nil, err
		}
		fields = append(fields, str)
		if s.offset == len(s.input) {
			return nil, fmt.Errorf("syntax err near %s: could not found ]", string(s.input[s.offset:]))
		}
		s.skipSpace()
		if s.input[s.offset] == rightSquareBracket {
			s.offset++
			break
		}
		if s.input[s.offset] != comma {
			return nil, fmt.Errorf("syntax err near %s", string(s.input[s.offset:]))
		}
		s.offset++
		s.skipSpace()
	}
	if len(fields) > 1 {
		return ast.NewMultiFields(fields...), nil
	}
	return ast.NewSingleField(fields[0]), nil
}

func (s *Parser) scanString() (string, error) {
	if s.input[s.offset] != singleQuotes && s.input[s.offset] != doubleQuotes {
		return "", fmt.Errorf(`syntax error near %q: could not find quotes`, string(s.input[s.offset:]))
	}
	result := make([]rune, 0)
	i := s.offset + 1
	for ; i < len(s.input) && s.input[i] != s.input[s.offset]; i++ {
		if s.input[i] == '\\' {
			i++
		}
		result = append(result, s.input[i])
	}
	if i == len(s.input) {
		return "", fmt.Errorf(`syntax error near %q: unmatched quotes`, string(s.input[s.offset:]))
	}
	s.offset = i + 1
	return string(result), nil
}

func (s *Parser) parseIndexes() (ast.Node, error) {
	nodes := make([]ast.Node, 0, 1)
	slice := make([]int, 0, 3)
	isSlice := false
	for {
		s.skipSpace()
		integer := 0
		if s.input[s.offset] != colon && s.input[s.offset] != rightSquareBracket {
			var err error
			integer, err = s.scanInteger()
			if err != nil {
				return nil, err
			}
		}
		if len(slice) == 3 {
			return nil, fmt.Errorf(`syntax error near %q`, string(s.input[s.offset:]))
		}
		slice = append(slice, integer)
		if s.offset == len(s.input) {
			return nil, fmt.Errorf("syntax err near %q : could not found ]", string(s.input[s.offset:]))
		}
		s.skipSpace()
		s.offset++
		switch s.input[s.offset-1] {
		case colon:
			isSlice = true
		case comma, rightSquareBracket:
			if !isSlice {
				nodes = append(nodes, ast.NewIndexField(slice[0]))
			} else {
				end, step := 0, 0
				if len(slice) > 1 {
					end = slice[1]
					if len(slice) > 2 {
						step = slice[2]
					}
				}
				nodes = append(nodes, ast.NewSliceField(slice[0], end, step))
			}
			slice = slice[:0]
			isSlice = false
			if s.input[s.offset-1] == rightSquareBracket {
				goto exit
			}
		default:
			return nil, fmt.Errorf("syntax err near %s", string(s.input[s.offset:]))
		}
	}
exit:
	if _, ok := nodes[0].(*ast.IndexField); ok && len(nodes) == 1 {
		return nodes[0], nil
	}
	return ast.NewIndexesField(nodes...), nil
}

func (s *Parser) scanInteger() (int, error) {
	off := s.offset + 1
	for ; off < len(s.input) && unicode.IsNumber(s.input[off]); off++ {
	}
	integer, err := strconv.Atoi(string(s.input[s.offset:off]))
	if err != nil {
		return 0, fmt.Errorf(`syntax error near %q: could not parse integer`, string(s.input[s.offset:]))
	}
	s.offset = off
	return integer, nil
}
