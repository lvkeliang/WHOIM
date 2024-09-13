# WHOIM

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "type": "text",
  "status": "sent",
  "sender_id": "user_123",
  "receiver_id": "user_456",
  "device_id": "device_abc",
  "content": "Hello, how are you?",
  "target_type": "single",
  "timestamp": "2024-01-01T12:00:00Z",
  "extra_data": {}
}
```

docker run -d --name rmqnamesrv --ip 172.19.0.3 --network WHOIM -p 9876:9876 -v ~/rocketmq/data/namesrv/logs:/home/rocketmq/logs -v ~/rocketmq/data/namesrv/store:/home/rocketmq/store rocketmqinc/rocketmq sh mqnamesrv


docker run -d --name rmqbroker \
--network WHOIM \
-p 10911:10911 \
-v E:/rocketmq/data/broker/logs:/home/rocketmq/logs \
-v E:/rocketmq/data/broker/store:/home/rocketmq/store \
-v E:/rocketmq/conf/broker.conf:/home/rocketmq/rocketmq-4.4.0/conf/broker.conf \
apache/rocketmq:latest sh mqbroker -n 172.19.0.3:9876 -c /home/rocketmq/rocketmq-5.3.0/conf/broker.conf

docker run -d --name rmqbroker --network WHOIM -p 10911:10911 -p 10909:10909 -v E:/rocketmq/data/broker/logs:/home/rocketmq/logs -v E:/rocketmq/data/broker/store:/home/rocketmq/store -v E:/rocketmq/conf/broker.conf:/home/rocketmq/rocketmq-4.4.0/conf/broker.conf apache/rocketmq:latest sh mqbroker -n 172.19.0.3:9876 -c /home/rocketmq/rocketmq-5.3.0/conf/broker.conf

brokerIP1=127.0.0.1