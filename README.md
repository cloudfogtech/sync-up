## Environment Variables

| Name                                   | Example                                  | Required | Default   | Description |
|----------------------------------------|------------------------------------------|----------|-----------|-------------|
| SYNCUP_PORT                            | 13003                                    | Optional | 13003     |             |
| SYNCUP_DB                              | sqlite                                   | Optional | sqlite    |             |
| SYNCUP_DB_DSN                          | sqlite.db                                | Optional | sqlite.db |             |
| SYNCUP_USERNAME                        | admin                                    | Optional |           |             |
| SYNCUP_PASSWORD                        | admin                                    | Optional |           |             |
| SYNCUP_SECRET_KEY                      | b8555b47784f612c36c7ca5561c9a3cfd806800a | Optional |           |             |
| SYNCUP_{service_id}_TYPE               | redis-rdb-c                              | Required |           |             |
| SYNCUP_{service_id}_CONTAINER          | redis                                    | Required |           |             |
| SYNCUP_{service_id}_LOCAL_DIR_PATH     | /temp                                    | Required |           |             |
| SYNCUP_{service_id}_CRON               | 0 0 1 * * *                              | Optional |           |             |
| SYNCUP_{service_id}_PASSWORD           | abc1234                                  | Required |           |             |
| SYNCUP_{service_id}_RDB_FILE_PATH      | /data/dump.rdb                           | Optional |           |             |
| SYNCUP_{service_id}_RC                 | rclone                                   | Required |           |             |
| SYNCUP_{service_id}_RC_DIR_PATH        | /backup                                  | Required |           |             |
| SYNCUP_{service_id}_RC_BANDWIDTH_LIMIT | 1M                                       | Optional |           |             |
| SYNCUP_{service_id}_RC_REMOTE_NAME     | webdav                                   | Required |           |             |
| SYNCUP_{service_id}_RC_REMOTE_PATH     | /backup/redis                            | Required |           |             |