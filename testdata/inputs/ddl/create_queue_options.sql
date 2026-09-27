CREATE QUEUE IF NOT EXISTS UserTasks (
  UserId INT64 NOT NULL,
  MessageId STRING(36) NOT NULL,
  Payload JSON NOT NULL
) PRIMARY KEY (UserId ASC, MessageId DESC),
INTERLEAVE IN PARENT Users ON DELETE CASCADE,
ROW DELETION POLICY (OLDER_THAN(DeliverTime, INTERVAL 7 DAY)),
OPTIONS (receive_mode = 'PULL', disable_send = false, disable_delivery = false, locality_group = 'ssd')
