CREATE PROPERTY GRAPH g NODE TABLES (
  T DEFAULT LABEL OPTIONS (description = 'A node', synonyms = ['vertex', 'entity'])
    PROPERTIES (id OPTIONS (description = 'Identifier'), name AS title OPTIONS (synonyms = ['caption', 'heading'], description = 'Title')),
  U DEFAULT LABEL OPTIONS (synonyms = ['other']) NO PROPERTIES
)
