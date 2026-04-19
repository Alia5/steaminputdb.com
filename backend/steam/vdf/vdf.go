package vdf

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type SyntaxError struct {
	msg    string
	Offset int
}

func (e *SyntaxError) Error() string { return e.msg }

var (
	reInt   = regexp.MustCompile(`^-?\d+$`)
	reFloat = regexp.MustCompile(`^-?\d+\.\d+$`)
	reBool  = regexp.MustCompile(`(?i)^(true|false)$`)
)

func autoType(s string) any {
	if reInt.MatchString(s) {
		if n, err := strconv.ParseInt(s, 10, 64); err == nil {
			return n
		}
	}
	if reFloat.MatchString(s) {
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
	}
	if reBool.MatchString(s) {
		return strings.EqualFold(s, "true")
	}
	return s
}

func Parse(data string) (map[string]any, error) {
	p := &parser{input: data}
	return p.parseMap()
}

type parser struct {
	input string
	pos   int
}

func (p *parser) skipWhitespaceAndComments() {
	for p.pos < len(p.input) {
		b := p.input[p.pos]
		if b == ' ' || b == '\t' || b == '\r' || b == '\n' {
			p.pos++
			continue
		}
		if b == '/' && p.pos+1 < len(p.input) && p.input[p.pos+1] == '/' {
			for p.pos < len(p.input) && p.input[p.pos] != '\n' {
				p.pos++
			}
			continue
		}
		break
	}
}

func (p *parser) readToken() (string, error) {
	p.skipWhitespaceAndComments()
	if p.pos >= len(p.input) {
		return "", &SyntaxError{"unexpected end of input", p.pos}
	}
	if p.input[p.pos] == '"' {
		return p.readQuotedString()
	}
	return p.readUnquotedString()
}

func (p *parser) readQuotedString() (string, error) {
	p.pos++
	var sb strings.Builder
	for p.pos < len(p.input) {
		b := p.input[p.pos]
		if b == '\\' && p.pos+1 < len(p.input) {
			p.pos++
			switch p.input[p.pos] {
			case '"':
				sb.WriteByte('"')
			case '\\':
				sb.WriteByte('\\')
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			default:
				sb.WriteByte('\\')
				sb.WriteByte(p.input[p.pos])
			}
			p.pos++
			continue
		}
		if b == '"' {
			p.pos++
			return sb.String(), nil
		}
		sb.WriteByte(b)
		p.pos++
	}
	return "", &SyntaxError{"unterminated string", p.pos}
}

func (p *parser) readUnquotedString() (string, error) {
	start := p.pos
	for p.pos < len(p.input) {
		b := p.input[p.pos]
		if b == ' ' || b == '\t' || b == '\r' || b == '\n' || b == '"' || b == '{' || b == '}' {
			break
		}
		p.pos++
	}
	if p.pos == start {
		return "", &SyntaxError{fmt.Sprintf("expected token, got %q", p.input[p.pos]), p.pos}
	}
	return p.input[start:p.pos], nil
}

func (p *parser) parseMap() (map[string]any, error) {
	result := make(map[string]any)
	for {
		p.skipWhitespaceAndComments()
		if p.pos >= len(p.input) {
			return result, nil
		}
		if p.input[p.pos] == '}' {
			return result, nil
		}
		key, err := p.readToken()
		if err != nil {
			return nil, err
		}
		p.skipWhitespaceAndComments()
		if p.pos >= len(p.input) {
			return nil, &SyntaxError{fmt.Sprintf("unexpected end after key %q", key), p.pos}
		}
		if p.input[p.pos] == '{' {
			p.pos++
			children, err := p.parseMap()
			if err != nil {
				return nil, err
			}
			p.skipWhitespaceAndComments()
			if p.pos >= len(p.input) || p.input[p.pos] != '}' {
				return nil, &SyntaxError{fmt.Sprintf("expected '}' for key %q", key), p.pos}
			}
			p.pos++
			existing, exists := result[key]
			if !exists {
				result[key] = children
			} else {
				result[key] = arrayify(existing, children)
			}
		} else {
			val, err := p.readToken()
			if err != nil {
				return nil, err
			}
			typed := autoType(val)
			existing, exists := result[key]
			if !exists {
				result[key] = typed
			} else {
				result[key] = arrayify(existing, typed)
			}
		}
	}
}

func arrayify(existing, incoming any) any {
	if arr, ok := existing.([]any); ok {
		return append(arr, incoming)
	}
	return []any{existing, incoming}
}
