import re

def ite(condition, b1, b2): 
	return b1 if condition else b2

PAST_event_1speed = None
PAST_event_1speedLimit = None
PAST_event_2objectDetected = None
PAST_event_1deltaSpeed = None
event_1deltaSpeed = None
PAST_event_1underLimit = None
event_1underLimit = None
PAST_event_2speed = None
event_2speed = None
PAST_event_2detected = None
event_2detected = None

if __name__ == '__main__':
	while True:
		inp = input()
		list_inp = inp.split(',') 
		if re.match(r'car,\d+,\d+$',inp) : 
			event_1speed = int(int(list_inp[1])) 
			event_1speedLimit = int(int(list_inp[2])) 
			event_1deltaSpeed = event_1speed - ite(PAST_event_1speed!= None,PAST_event_1speed,event_1speed)
			event_1underLimit = event_1speed <= event_1speedLimit
			PAST_event_1speed = event_1speed
			PAST_event_1speedLimit = event_1speedLimit
			PAST_event_1deltaSpeed = event_1deltaSpeed
			PAST_event_1underLimit = event_1underLimit
			print('car' + ',' +  str(event_1deltaSpeed) + ',' +  str(event_1underLimit))


		if re.match(r'detect,\d+$',inp) : 
			event_2objectDetected = bool(int(list_inp[1])) 
			event_2speed = ite(event_2objectDetected and ite(PAST_event_2speed!= None,PAST_event_2speed, 10 ) <  5 ,  0 ,  5 )
			event_2detected = ite(PAST_event_2objectDetected!= None,PAST_event_2objectDetected,event_2objectDetected)
			PAST_event_2objectDetected = event_2objectDetected
			PAST_event_2speed = event_2speed
			PAST_event_2detected = event_2detected
			print('detect' + ',' +  str(event_2speed) + ',' +  str(event_2detected))

