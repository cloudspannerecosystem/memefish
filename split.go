package memefish

import "github.com/cloudspannerecosystem/memefish/token"

type RawStatement struct {
	Pos, End  token.Pos
	Statement string
}

// SplitRawStatements splits an input string to statement strings at terminating semicolons without parsing.
// It preserves all comments.
// Statements are terminated by `;`, `<eof>` or `;<eof>` and the minimum output will be []string{""}.
// See [terminating semicolons].
// This function won't panic but return error if lexer become error state.
// filepath can be empty, it is only used in error message.
//
// [terminating semicolons]: https://cloud.google.com/spanner/docs/reference/standard-sql/lexical#terminating_semicolons
func SplitRawStatements(filepath, s string) ([]*RawStatement, error) {
	lex := &Lexer{
		File: &token.File{
			FilePath: filepath,
			Buffer:   s,
		},
	}

	var result []*RawStatement
	var firstPos token.Pos
	for {
		if lex.Token.Kind == ";" {
			result = append(result, &RawStatement{Pos: firstPos,
				End:       lex.Token.Pos,
				Statement: s[firstPos:lex.Token.Pos],
			})
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
				result = append(result, &RawStatement{Pos: firstPos,
					End:       lex.Token.Pos,
					Statement: s[firstPos:lex.Token.Pos],
				})
			}
			break
		}
	}
	if len(result) == 0 {
		return []*RawStatement{
			{Pos: token.Pos(0), End: token.Pos(0), Statement: ""},
		}, nil
	}
	return result, nil
}
