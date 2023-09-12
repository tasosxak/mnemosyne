import re
def ite(condition, b1, b2): 
	return b1 if condition else b2

PAST_event_1carA = None
PAST_event_1carB = None
PAST_event_1distance = None
PAST_event_1SafeDistance = None
event_1SafeDistance = None
PAST_event_1CarA = None
event_1CarA = None
PAST_event_1CarB = None
event_1CarB = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'report,\w+,\w+,-?\d+$',inp):
			event_1carA = str(list_inp[1])
			event_1carB = str(list_inp[2])
			event_1distance = int(list_inp[3])
			event_1SafeDistance = bool(event_1distance <  10  and event_1distance >  5 )
			event_1CarA = str(event_1carA)
			event_1CarB = str(event_1carB)
			PAST_event_1carA = event_1carA
			PAST_event_1carB = event_1carB
			PAST_event_1distance = event_1distance
			PAST_event_1SafeDistance = event_1SafeDistance
			PAST_event_1CarA = event_1CarA
			PAST_event_1CarB = event_1CarB
			print('distance' + ',' +  str(event_1CarA) + ',' +  str(event_1CarB) + ',' +  str(event_1SafeDistance))

