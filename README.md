# Freethewump.us Web
The web part of the freethewump.us service. There are several environment variables that should be configured:
- `HOST` - The host that uploads will be accepted from. Defaults to `freethewump.us`.
- `REDIS_HOST` - The Redis host. Defaults to `localhost:6379`.
- `REDIS_PASSWORD` - The password for the Redis database. Defaults to none.
- `RETHINK_HOST` - The RethinkDB host. Defaults to `127.0.0.1:28015`.
- `RETHINK_PASSWORD` - The RethinkDB password. Defaults to none.
- `RETHINK_USER` - The RethinkDB user. Defaults to `admin`.
- `AWS_SECRET_ACCESS_KEY` - The AWS secret access key for the default bucket. **This is required even in development.**
- `AWS_ACCESS_KEY_ID` - The AWS access key ID for the default bucket. **This is required even in development.**
- `S3_BUCKET` - The default S3 bucket. **This is required even in development.**
- `S3_ENDPOINT` - The endpoint for the default S3 bucket. **This is required even in development.**

