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

def cleanup_topics(): admin_client.delete_topics(['system_topic','rm_topic','rules_topic','rules_topic',], operation_timeout=30)

atexit.register(cleanup_topics)

global_Danger =  False 
PAST_global_Danger = None
global_Violations =  0 
PAST_global_Violations = None
PAST_event_1level = None
PAST_event_2verd = None
PAST_event_3danger = None
event_3danger = None
PAST_event_3violations = None
event_3violations = None

if __name__ == '__main__':
	consumer = Consumer(consumer_config)
	consumer.subscribe(["system_topic", "rm_topic", "rules_topic", ])
	producer = Producer(producer_config)

	while True:
		inp = consumer.poll()
		if inp is None or inp.error():
			continue
		list_inp = str(inp.value().decode('utf-8')).split(',')
		if inp.topic() == "system_topic" and re.match(r'system,-?\d+$',inp.value().decode('utf-8')):
			event_1level = int(list_inp[1])
			PAST_global_Danger = global_Danger
			global_Danger = bool(ite(event_1level - ite(PAST_event_1level!= None,PAST_event_1level,event_1level) >  0 ,  True ,  False ))
			PAST_event_1level = event_1level
			print('STATUS' + ',' +  str(global_Danger))


		if inp.topic() == "rm_topic" and re.match(r'verdict,\d+$',inp.value().decode('utf-8')):
			event_2verd = bool(int(list_inp[1]))
			PAST_global_Violations = global_Violations
			global_Violations = int(ite(event_2verd, global_Violations +  1 , global_Violations))
			PAST_event_2verd = event_2verd
			if (global_Violations >  2 ) :
				producer.produce(topic='rules_topic', value='shutdown' + '')


		if inp.topic() == "rules_topic" and re.match(r'shutdown$',inp.value().decode('utf-8')):
			PAST_global_Violations = global_Violations
			global_Violations = int( 0 )
			PAST_global_Danger = global_Danger
			global_Danger = bool( False )
			event_3violations = int(ite(PAST_global_Violations!= None,PAST_global_Violations,global_Violations))
			event_3danger = bool(ite(PAST_global_Danger!= None,PAST_global_Danger,global_Danger))
			PAST_event_3danger = event_3danger
			PAST_event_3violations = event_3violations
			print('STOP' + ',' +  str(event_3danger) + ',' +  str(event_3violations))

