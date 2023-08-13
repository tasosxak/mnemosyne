PAST_event_0speed = None
PAST_event_0speedLimit = None
PAST_event_1objectDetected = None
PAST_event_0deltaSpeed = None
event_0deltaSpeed = None
PAST_event_0underLimit = None
event_0underLimit = None
PAST_event_1speed = None
event_1speed = None
PAST_event_1detected = None
event_1detected = None
while True:
	inp = input()
	inp = inp.split(',') 
	if inp[0] == 'car' : 
		event_0speed = int(int(inp[1])) 
		event_0speedLimit = int(int(inp[2])) 
		event_0deltaSpeed = event_0speed - (PAST_event_0speed if PAST_event_0speed!= None else event_0speed)
		event_0underLimit = event_0speed <= event_0speedLimit
		PAST_event_0speed = event_0speed
		PAST_event_0speedLimit = event_0speedLimit
		PAST_event_0deltaSpeed = event_0deltaSpeed
		PAST_event_0underLimit = event_0underLimit
		print('car' + ',' + str(event_0deltaSpeed) + ',' + str(event_0underLimit)) 

	if inp[0] == 'detect' : 
		event_1objectDetected = bool(int(inp[1])) 
		event_1speed =  0  if event_1objectDetected and (PAST_event_1speed if PAST_event_1speed!= None else  10 ) <  5  else  5 
		event_1detected = (PAST_event_1objectDetected if PAST_event_1objectDetected!= None else event_1objectDetected)
		PAST_event_1objectDetected = event_1objectDetected
		PAST_event_1speed = event_1speed
		PAST_event_1detected = event_1detected
		print('detect' + ',' + str(event_1speed) + ',' + str(event_1detected)) 
