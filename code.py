PAST_event_0speed = None
PAST_event_0speedLimit = None
PAST_event_1objectDetected = None
PAST_event_0deltaSpeed = None
event_0deltaSpeed = None
PAST_event_0underLimit = None
event_0underLimit = None
PAST_event_1speed = None
event_1speed = None
while True:
	inp = input()
	inp = inp.split(',') 
	if inp[0] == 'car' : 
		event_0speed = int(inp[1]) 
		event_0speedLimit = int(inp[2]) 
		event_0deltaSpeed = event_0speed - (PAST_event_0speed if PAST_event_0speed else event_0speed)
		event_0underLimit = event_0speed <= event_0speedLimit
		PAST_event_0speed = event_0speed
		PAST_event_0speedLimit = event_0speedLimit
		PAST_event_0deltaSpeed = event_0deltaSpeed
		PAST_event_0underLimit = event_0underLimit
		print('car' + ',' + str(event_0deltaSpeed) + ',' + str(event_0underLimit)) 

	if inp[0] == 'detect' : 
		event_1objectDetected = bool(inp[1]) 
		event_1speed =  0  if event_1objectDetected and event_1speed <  5  else  5 
		PAST_event_1objectDetected = event_1objectDetected
		PAST_event_1speed = event_1speed
		print('detect' + ',' + str(event_1speed)) 
