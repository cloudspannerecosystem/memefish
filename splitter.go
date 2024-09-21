package memefish

import "github.com/cloudspannerecosystem/memefish/token"

// SplitRawStatements splits input to statement strings at terminating semicolons without parsing.
// Statements are terminated by `;`, `<eof>` or `;<eof>` and the minimum output will be []string{""}.
// See [terminating semicolons].
// This function won't panic but return error if lexer become error state.
// filepath can be used in error message.
//
// [terminating semicolons]: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#terminating_semicolons
func SplitRawStatements(filepath, input string) ([]string, error) {
	lex := &Lexer{
		File: &token.File{
			FilePath: filepath,
			Buffer:   input,
		},
	}

	var result []string
	var firstPos token.Pos
	for {
		if lex.Token.Kind == ";" {
			result = append(result, input[firstPos:lex.Token.Pos])
			if err := lex.NextToken(); err != nil {
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
			if lex.Token.Pos != firstPos {
				result = append(result, input[firstPos:lex.Token.Pos])
			}
			break
		}
	}
	if len(result) == 0 {
		return []string{""}, nil
	}
	return result, nil
}
