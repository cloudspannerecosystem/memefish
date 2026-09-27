-- https://cloud.google.com/spanner/docs/full-text-search/search-indexes#search-index-schema-definitions
CREATE TABLE Albums (
                        AlbumId STRING(MAX) NOT NULL,
                        SingerId INT64 NOT NULL,
                        ReleaseTimestamp INT64 NOT NULL,
                        AlbumTitle STRING(MAX),
                        Rating FLOAT64,
                        AlbumTitle_Tokens TOKENLIST AS (TOKENIZE_FULLTEXT(AlbumTitle)) HIDDEN,
                        Rating_Tokens TOKENLIST AS (TOKENIZE_NUMBER(Rating)) HIDDEN
) PRIMARY KEY(AlbumId)