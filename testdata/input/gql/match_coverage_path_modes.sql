GRAPH FinGraph
-- Path Modes (singular/plural)
MATCH WALK PATH (aw)-[ew]->(bw)
MATCH WALK PATHS (aws)-[ews]->(bws)
MATCH TRAIL PATH (atvar)-[et]->(btvar)
MATCH SIMPLE PATH (asvar)-[es]->(bsvar)
MATCH ACYCLIC PATH (aa)-[ea]->(ba)
RETURN *
