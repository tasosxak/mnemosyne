import re
def ite(condition, b1, b2): 
	return b1 if condition else b2

PAST_event_1vDetected = None
PAST_event_1vGroundTruth = None
PAST_event_1steering = None
PAST_event_1speed = None
PAST_event_1correctDetection = None
event_1correctDetection = None
PAST_event_1deltaSteering = None
event_1deltaSteering = None
PAST_event_1deltaSpeed = None
event_1deltaSpeed = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'instance,-?\d+,-?\d+,-?\d+,-?\d+$',inp):
			event_1vDetected = int(list_inp[1])
			event_1vGroundTruth = int(list_inp[2])
			event_1steering = int(list_inp[3])
			event_1speed = int(list_inp[4])
			event_1correctDetection = bool(event_1vDetected == event_1vGroundTruth)
			event_1deltaSteering = int(abs(event_1steering - ite(PAST_event_1steering!= None,PAST_event_1steering,event_1steering)))
			event_1deltaSpeed = int(event_1speed - ite(PAST_event_1speed!= None,PAST_event_1speed,event_1speed))
			PAST_event_1vDetected = event_1vDetected
			PAST_event_1vGroundTruth = event_1vGroundTruth
			PAST_event_1steering = event_1steering
			PAST_event_1speed = event_1speed
			PAST_event_1correctDetection = event_1correctDetection
			PAST_event_1deltaSteering = event_1deltaSteering
			PAST_event_1deltaSpeed = event_1deltaSpeed
			print('instance' + ',' +  str(event_1correctDetection) + ',' +  str(event_1deltaSteering) + ',' +  str(event_1deltaSpeed))

