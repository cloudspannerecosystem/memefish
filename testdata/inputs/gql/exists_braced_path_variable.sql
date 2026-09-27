GRAPH FinGraph
RETURN
  EXISTS { p = (a)-[e]->(b) } AS local_graph,
  EXISTS { GRAPH FinGraph p = (a)-[e]->(b) } AS explicit_graph
