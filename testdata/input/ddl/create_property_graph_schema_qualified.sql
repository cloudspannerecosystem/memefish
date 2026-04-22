CREATE PROPERTY GRAPH MyGraph
  NODE TABLES (
    myschema.NodeTable KEY (id),
    SimpleNode KEY (id)
  )
  EDGE TABLES (
    myschema.EdgeTable
      KEY (src_id, dst_id)
      SOURCE KEY (src_id) REFERENCES myschema.NodeTable (id)
      DESTINATION KEY (dst_id) REFERENCES SimpleNode (id)
  )
