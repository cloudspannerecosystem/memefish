SELECT ChangeRecord FROM READ_SingersNameStream (
  start_timestamp => "2022-05-01T09:00:00Z",
  end_timestamp => NULL,
  partition_token => NULL,
  heartbeat_milliseconds => 10000
)