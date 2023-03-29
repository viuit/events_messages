import asyncio
import aioamqp
from loguru import logger
import time
from read_json import get_message

message = get_message()


def url():
    port = "5672"
    user = "guest"
    password = "guest"
    v_host = "new"
    host = "localhost"
    return f"amqp://{user}:{password}@{host}:{port}/"


async def connect():
    try:
        transport, protocol = await aioamqp.from_url(url=url())
        logger.success("RabbitMQ connected!")
        return transport, protocol
    except aioamqp.AmqpClosedConnection as _err:
        logger.error(f"RabbitMQ connection error: {_err}")


async def send_message(protocol: aioamqp.AmqpProtocol):
    channel = await protocol.channel()
    while True:
        # message = get_message()
        # message = time.ctime(time.time())
        await channel.publish(
            payload=message.encode(),
            exchange_name='amq.topic',
            routing_key='test'
        )
        logger.success(f"Отправлено сообщение: {message}")
        time.sleep(5)


async def send_another(protocol):
    channel = await protocol.channel()
    await channel.queue('task_queue', durable=True)
    await channel.basic_publish(
    payload="message".encode(),
    exchange_name='',
    routing_key='task_queue',
    properties={
        'delivery_mode': 2,
    },)
    logger.success("Отправлено сообщение")

        

async def close(transport, protocol):
    await protocol.close()
    transport.close()
    logger.success("RabbitMQ connect closed!")


async def main():
    transport, protocol = await connect()
    await send_message(protocol)
    # await close(transport, protocol)


def worker_loop():
    try:
        loop = asyncio.get_event_loop()
        loop.run_until_complete(main())
        loop.run_forever()
    except:
        logger.error(f"Ошибка")
        time.sleep(10)
        worker_loop()


if __name__ == "__main__":
    worker_loop()

        


