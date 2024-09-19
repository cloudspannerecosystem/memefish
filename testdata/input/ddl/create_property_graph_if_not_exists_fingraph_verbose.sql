CREATE PROPERTY GRAPH IF NOT EXISTS FinGraph
  NODE TABLES (
    Account AS Account
       DEFAULT LABEL
       PROPERTIES (create_time, is_blocked, nick_name AS name),
    Person
       LABEL Person
       PROPERTIES ALL COLUMNS
  )
  EDGE TABLES (
    PersonOwnAccount
      SOURCE KEY (id) REFERENCES Person (id)
      DESTINATION KEY (account_id) REFERENCES Account (id)
      LABEL Owns PROPERTIES ALL COLUMNS EXCEPT (a),
    AccountTransferAccount
      SOURCE KEY (id) REFERENCES Account (id)
      DESTINATION KEY (to_id) REFERENCES Account (id)
      LABEL Transfers NO PROPERTIES
  )