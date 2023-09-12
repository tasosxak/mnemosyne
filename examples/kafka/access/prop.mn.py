import re
from confluent_kafka import Consumer, Producer, KafkaError
import atexit
from confluent_kafka.admin import AdminClient
# Define the Kafka broker address
broker_address = 'localhost:9092'

# Create an AdminClient instance
admin_client = AdminClient({'bootstrap.servers': broker_address})

# Kafka configuration for the consumer
consumer_config = {
    'bootstrap.servers': 'localhost',
    'group.id': 'consumer_group',
    'auto.offset.reset': 'earliest'
    }

producer_config = {
    'bootstrap.servers': 'localhost',
}
def ite(condition, b1, b2): 
	return b1 if condition else b2

def cleanup_topics(): admin_client.delete_topics(['topic_1',], operation_timeout=30)

atexit.register(cleanup_topics)

PAST_event_1file = None
PAST_event_2file = None
PAST_event_3file = None
PAST_event_3user = None
PAST_event_1id = None
event_1id = None
PAST_event_2id = None
event_2id = None
PAST_event_3fileId = None
event_3fileId = None
PAST_event_3userId = None
event_3userId = None

if __name__ == '__main__':
	consumer = Consumer(consumer_config)
	consumer.subscribe(["topic_1", ])
	producer = Producer(producer_config)

	while True:
		inp = consumer.poll()
		if inp is None or inp.error():
			continue
		list_inp = str(inp.value().decode('utf-8')).split(',')
		if inp.topic() == "topic_1" and re.match(r'open,\w+$',inp.value().decode('utf-8')):
			event_1file = str(list_inp[1])
			event_1id = str(event_1file)
			PAST_event_1file = event_1file
			PAST_event_1id = event_1id
			print('open' + ',' +  str(event_1id))


		if inp.topic() == "topic_1" and re.match(r'close,\w+$',inp.value().decode('utf-8')):
			event_2file = str(list_inp[1])
			event_2id = str(event_2file)
			PAST_event_2file = event_2file
			PAST_event_2id = event_2id
			print('close' + ',' +  str(event_2id))


		if inp.topic() == "topic_1" and re.match(r'write,\w+,\w+$',inp.value().decode('utf-8')):
			event_3file = str(list_inp[1])
			event_3user = str(list_inp[2])
			event_3fileId = str(event_3file)
			event_3userId = str(event_3user)
			PAST_event_3file = event_3file
			PAST_event_3user = event_3user
			PAST_event_3fileId = event_3fileId
			PAST_event_3userId = event_3userId
			print('userId' + ',' +  str(event_3fileId) + ',' +  str(event_3userId))

