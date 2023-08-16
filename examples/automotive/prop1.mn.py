def ite(condition, b1, b2): 
	return b1 if condition else b2
PAST_event_0vDetected = None
PAST_event_0vGroundTruth = None
PAST_event_0steering = None
PAST_event_0speed = None
PAST_event_0correctDetection = None
event_0correctDetection = None
PAST_event_0deltaSteering = None
event_0deltaSteering = None
PAST_event_0deltaSpeed = None
event_0deltaSpeed = None
while True:
	inp = input()
	inp = inp.split(',') 
	if inp[0] == 'instance' : 
		event_0vDetected = int(int(inp[1])) 
		event_0vGroundTruth = int(int(inp[2])) 
		event_0steering = int(int(inp[3])) 
		event_0speed = int(int(inp[4])) 
		event_0correctDetection = event_0vDetected == event_0vGroundTruth
		event_0deltaSteering = abs(event_0steering - ite(PAST_event_0steering!= None,PAST_event_0steering,event_0steering))
		event_0deltaSpeed = event_0speed - ite(PAST_event_0speed!= None,PAST_event_0speed,event_0speed)
		PAST_event_0vDetected = event_0vDetected
		PAST_event_0vGroundTruth = event_0vGroundTruth
		PAST_event_0steering = event_0steering
		PAST_event_0speed = event_0speed
		PAST_event_0correctDetection = event_0correctDetection
		PAST_event_0deltaSteering = event_0deltaSteering
		PAST_event_0deltaSpeed = event_0deltaSpeed
		print('instance' + ',' + str(event_0correctDetection) + ',' + str(event_0deltaSteering) + ',' + str(event_0deltaSpeed)) 
