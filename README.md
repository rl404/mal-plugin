# mal-plugin

[![Go Reference](https://pkg.go.dev/badge/github.com/rl404/mal-plugin.svg)](https://pkg.go.dev/github.com/rl404/mal-plugin)

Additional libraries used by:
- [go-malscraper](https://github.com/rl404/go-malscraper)
- [mal-api](https://github.com/rl404/mal-api)

Contain interfaces for:
- Cache
  - [Redis](https://redis.io/)
  - [In-memory](https://github.com/allegro/bigcache)
  - [Memcached](https://memcached.org/)
  - No cache mock
- PubSub
  - [Redis](https://redis.io/)
  - [RabbitMQ](https://www.rabbitmq.com/)
  - [NSQ](https://nsq.io/)

Though, you can use them for other projects.

**See usage examples in `example` folder.*