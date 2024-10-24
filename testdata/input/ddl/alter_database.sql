ALTER DATABASE dbname SET OPTIONS (
    optimizer_version=2,
    optimizer_statistics_package='auto_20191128_14_47_22UTC',
    version_retention_period='7d',
    enable_key_visualizer=true,
    default_leader='europe-west1'
  )