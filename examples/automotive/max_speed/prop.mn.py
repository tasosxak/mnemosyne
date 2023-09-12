import re
def ite(condition, b1, b2): 
	return b1 if condition else b2

global_MaxSpeed =  0 
PAST_global_MaxSpeed = None
PAST_event_1maker = None
PAST_event_1speed = None
PAST_event_1Maker = None
event_1Maker = None
PAST_event_1NewRecord = None
event_1NewRecord = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'recorded,\w+,-?\d+$',inp):
			event_1maker = str(list_inp[1])
			event_1speed = int(list_inp[2])
			event_1NewRecord = bool(ite(PAST_global_MaxSpeed!= None,PAST_global_MaxSpeed,global_MaxSpeed) < event_1speed)
			PAST_global_MaxSpeed = global_MaxSpeed
			global_MaxSpeed = int(ite(ite(PAST_global_MaxSpeed!= None,PAST_global_MaxSpeed,global_MaxSpeed) < event_1speed, event_1speed, ite(PAST_global_MaxSpeed!= None,PAST_global_MaxSpeed,global_MaxSpeed)))
			event_1Maker = str(event_1maker)
			PAST_event_1maker = event_1maker
			PAST_event_1speed = event_1speed
			PAST_event_1Maker = event_1Maker
			PAST_event_1NewRecord = event_1NewRecord
			print('fast' + ',' +  str(event_1Maker) + ',' +  str(event_1NewRecord))

