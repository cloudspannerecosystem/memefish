package memefish

import "github.com/cloudspannerecosystem/memefish/token"

// SeparateRawStatements separates s to statements without parsing.
// Statements are terminated by `;`, `<eof>` or `;<eof>` and the minimum output will be []string{""}.
// See [terminating semicolons].
// This function won't panic but return error if lexer become error state.
// filepath can be used in error message.
//
// [terminating semicolons]: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#terminating_semicolons
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
			err := lex.NextToken()
			if err != nil {
				return nil, err
			}
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
