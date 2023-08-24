import re

def ite(condition, b1, b2): 
	return b1 if condition else b2

PAST_event_1speed = None
PAST_event_1speedLimit = None
PAST_event_2objectDetected = None
PAST_event_3id = None
PAST_event_1deltaSpeed = None
event_1deltaSpeed = None
PAST_event_1underLimit = None
event_1underLimit = None
PAST_event_2speed = None
event_2speed = None
PAST_event_2detected = None
event_2detected = None
PAST_event_3newId = None
event_3newId = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'car,-?\d+,-?\d+$',inp) : 
			event_1speed = int(list_inp[1])
			event_1speedLimit = int(list_inp[2])
			event_1deltaSpeed = int(event_1speed - ite(PAST_event_1speed!= None,PAST_event_1speed,event_1speed))
			event_1underLimit = bool(event_1speed <= event_1speedLimit)
			PAST_event_1speed = event_1speed
			PAST_event_1speedLimit = event_1speedLimit
			PAST_event_1deltaSpeed = event_1deltaSpeed
			PAST_event_1underLimit = event_1underLimit
			print('car' + ',' +  str(event_1deltaSpeed) + ',' +  str(event_1underLimit))


		if re.match(r'detect,\d+$',inp) : 
			event_2objectDetected = bool(int(list_inp[1]))
			event_2speed = int(ite(event_2objectDetected and ite(PAST_event_2speed!= None,PAST_event_2speed, 10 ) <  5 ,  0 ,  5 ))
			event_2detected = bool(ite(PAST_event_2objectDetected!= None,PAST_event_2objectDetected,event_2objectDetected))
			PAST_event_2objectDetected = event_2objectDetected
			PAST_event_2speed = event_2speed
			PAST_event_2detected = event_2detected
			print('detect' + ',' +  str(event_2speed) + ',' +  str(event_2detected))


		if re.match(r'example,\w+$',inp) : 
			event_3id = str(list_inp[1])
			event_3newId = str(event_3id)
			PAST_event_3id = event_3id
			PAST_event_3newId = event_3newId
			print('newExample' + ',' +  str(event_3newId))

