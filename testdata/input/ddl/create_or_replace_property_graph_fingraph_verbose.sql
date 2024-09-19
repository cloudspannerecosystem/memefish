CREATE OR REPLACE PROPERTY GRAPH FinGraph
  NODE TABLES (
    Account AS Account -- element_alias
      KEY (id) -- element_key in node_element_key in element_keys
      -- label_and_property_list
      LABEL DetailedAccount -- LABEL label_name in element_label
        PROPERTIES (create_time, is_blocked, nick_name AS name) -- derived_property_list
      DEFAULT LABEL -- DEFAULT LABEL in element_label
        NO PROPERTIES -- NO PROPERTIES in element_properties
    ,
    Person
      -- no element_keys
      -- no element_label because of direct element_properties
      PROPERTIES ARE ALL COLUMNS EXCEPT (city) -- properties_are
  )
  EDGE TABLES (
    PersonOwnAccount AS PersonOwnAccount
      KEY (id, account_id)
      SOURCE KEY (id) REFERENCES Person -- source_key without column_name_list
      DESTINATION KEY (account_id) REFERENCES Account -- destination_key without column_name_list
      LABEL Owns
        PROPERTIES ALL COLUMNS,
    AccountTransferAccount
      SOURCE KEY (id) REFERENCES Account (id) -- source_key
      DESTINATION KEY (to_id) REFERENCES Account (id) -- destination_key
      LABEL Transfers -- LABEL label_name in element_label
      -- without element_properties
  )