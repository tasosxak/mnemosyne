import re

from kafka import KafkaProducer, KafkaConsumer

class KafkaEventHandler:
    def __init__(self, i_bootstrap_servers, i_out_topic, i_in_topic, i_group_id):
        self._m_bootstrap_servers = i_bootstrap_servers
        self._m_out_topic = i_out_topic
        self._m_in_topic = i_in_topic
        self._m_producer = self.__initialize_producer()
        self._m_consumer = self.__initialize_consumer(i_group_id)

    def __initialize_producer(self):
        return KafkaProducer(bootstrap_servers=self._m_bootstrap_servers, linger_ms=10)

    def __initialize_consumer(self, group_id):
        # Create a Kafka consumer
        consumer = KafkaConsumer(
            bootstrap_servers=self._m_bootstrap_servers,
            group_id=group_id,
            auto_offset_reset='latest'
        )

        # Subscribe to the incoming topic
        consumer.subscribe(topics=[self._m_in_topic])

        return consumer

    def send_event(self, event: str):
        # Publish a message to the output topic
        self._m_producer.send(self._m_out_topic, value=event.encode()).add_callback(
            lambda r: print(f'Event ({event}) published successfully!')).add_errback(
            lambda e: print(f'Error publishing event: {str(e)}'))

        # Flush the producer to ensure the message is sent
        self._m_producer.flush()

    def receive_feedback(self):
        # Continuously consume messages from the input topic
        for message in self._m_consumer:
            print(f'Received feedback message: {message.value.decode()}')
            return

    def producer_terminate(self):
        self._m_producer.close()

    def consumer_terminate(self):
        self._m_consumer.close()

    def shutdown(self):
        self.producer_terminate()
        self.consumer_terminate()

def ite(condition, b1, b2): 
	return b1 if condition else b2

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
	kafka_handler = KafkaEventHandler('localhost:9092','event_topic', 'feedback_topic','dejavu_group')
	while True:
		inp = input()
		list_inp = inp.split(',') 
		if re.match(r'open,\d+$',inp) : 
			event_1file = int(int(list_inp[1])) 
			event_1id = event_1file
			PAST_event_1file = event_1file
			PAST_event_1id = event_1id
			kafka_handler.send_event('open' + ',' +  str(event_1id))


		if re.match(r'close,\d+$',inp) : 
			event_2file = int(int(list_inp[1])) 
			event_2id = event_2file
			PAST_event_2file = event_2file
			PAST_event_2id = event_2id
			kafka_handler.send_event('close' + ',' +  str(event_2id))


		if re.match(r'write,\d+,\d+$',inp) : 
			event_3file = int(int(list_inp[1])) 
			event_3user = int(int(list_inp[2])) 
			event_3fileId = event_3file
			event_3userId = event_3user
			PAST_event_3file = event_3file
			PAST_event_3user = event_3user
			PAST_event_3fileId = event_3fileId
			PAST_event_3userId = event_3userId
			kafka_handler.send_event('userId' + ',' +  str(event_3fileId) + ',' +  str(event_3userId))

