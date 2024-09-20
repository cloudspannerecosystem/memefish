package memefish

import "github.com/cloudspannerecosystem/memefish/token"

func SeparateRawStatements(filepath, s string) ([]string, error) {
	lex := &Lexer{
		File: &token.File{
			FilePath: filepath,
			Buffer:   s,
		},
	}

	var result []string
	var firstPos token.Pos
	for {
		if lex.Token.Kind == ";" {
			result = append(result, s[firstPos:lex.Token.Pos])
			lex.nextToken()
			firstPos = lex.Token.Pos
			continue
		}

		err := lex.NextToken()
		if err != nil {
			return nil, err
		}

		if lex.Token.Kind == token.TokenEOF {
			if lex.Token.Pos == firstPos {
				break
			}
			result = append(result, s[firstPos:lex.Token.Pos])
			firstPos = lex.Token.Pos
			break
		}
	}
	if len(result) == 0 {
		return []string{""}, nil
	}
	return result, nil
}
