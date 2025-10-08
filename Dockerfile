# RabbitMQ instance to support STOMP as well, just for fun
FROM rabbitmq:3.13-management
RUN rabbitmq-plugins enable rabbitmq_stomp
