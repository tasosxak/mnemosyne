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

PAST_event_0file = None
PAST_event_1file = None
PAST_event_2file = None
PAST_event_2user = None
PAST_event_0id = None
event_0id = None
PAST_event_1id = None
event_1id = None
PAST_event_2fileId = None
event_2fileId = None
PAST_event_2userId = None
event_2userId = None

if __name__ == '__main__':
	kafka_handler = KafkaEventHandler('localhost:9092','event_topic', 'feedback_topic','dejavu_group')
	while True:
		inp = input()
		inp = inp.split(',') 
		if inp[0] == 'open' : 
			event_0file = int(int(inp[1])) 
			event_0id = event_0file
			PAST_event_0file = event_0file
			PAST_event_0id = event_0id
			kafka_handler.send_event('open' + ',' + str(event_0id))
			#kafka_handler.receive_feedback()


		if inp[0] == 'close' : 
			event_1file = int(int(inp[1])) 
			event_1id = event_1file
			PAST_event_1file = event_1file
			PAST_event_1id = event_1id
			kafka_handler.send_event('close' + ',' + str(event_1id))
			#kafka_handler.receive_feedback()


		if inp[0] == 'write' : 
			event_2file = int(int(inp[1])) 
			event_2user = int(int(inp[2])) 
			event_2fileId = event_2file
			event_2userId = event_2user
			PAST_event_2file = event_2file
			PAST_event_2user = event_2user
			PAST_event_2fileId = event_2fileId
			PAST_event_2userId = event_2userId
			kafka_handler.send_event('write' + ',' + str(event_2fileId) + ',' + str(event_2userId))
			kafka_handler.receive_feedback()

