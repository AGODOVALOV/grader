FROM rabbitmq:4.2.4-management-alpine

RUN set -eux; \
    apk add --no-cache curl; \
    curl -L -o /opt/rabbitmq/plugins/rabbitmq_delayed_message_exchange-4.2.0.ez \
      https://github.com/rabbitmq/rabbitmq-delayed-message-exchange/releases/download/v4.2.0/rabbitmq_delayed_message_exchange-4.2.0.ez; \
    rabbitmq-plugins enable --offline rabbitmq_delayed_message_exchange