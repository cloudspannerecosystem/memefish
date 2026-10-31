CREATE PROPERTY GRAPH g NODE TABLES (
  T PROPERTIES (id OPTIONS (synonyms = ['key']), id + 1 AS next_id OPTIONS (description = 'Next identifier')),
  U LABEL Named PROPERTIES (name OPTIONS (description = 'Name'))
)
