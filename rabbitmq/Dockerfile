FROM rabbitmq:3.10.7-management

RUN apt-get update && \
    apt-get install -y libpq5 wait-for-it && \
    apt-get clean

RUN rabbitmq-plugins enable rabbitmq_stream

RUN rabbitmq-plugins enable rabbitmq_stream_management

# CMD rabbitmqadmin declare binding source=amq.topic destination=new routing_key=test && rabbitmqadmin declare queue name=new queue_type=stream


# CMD sleep 5 && rabbitmqadmin declare queue name=new queue_type=stream
# && rabbitmqadmin declare binding source=amq.topic destination=new routing_key=test