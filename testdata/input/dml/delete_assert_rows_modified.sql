DELETE FROM UserTasks WHERE UserId = @userId AND MessageId = @messageId ASSERT_ROWS_MODIFIED 1
