import re

def ite(condition, b1, b2): 
	return b1 if condition else b2

global_Distance =  -1 
PAST_global_Distance = None
PAST_event_1carA = None
PAST_event_1carB = None
PAST_event_1xA = None
PAST_event_1yA = None
PAST_event_1xB = None
PAST_event_1yB = None
PAST_event_2carA = None
PAST_event_2carB = None
PAST_event_2xA = None
PAST_event_2yA = None
PAST_event_2xB = None
PAST_event_2yB = None
PAST_event_3carA = None
PAST_event_3carB = None
PAST_event_4carA = None
PAST_event_4carB = None
PAST_event_1CarA = None
event_1CarA = None
PAST_event_1CarB = None
event_1CarB = None
PAST_event_2CarA = None
event_2CarA = None
PAST_event_2CarB = None
event_2CarB = None
PAST_event_2HighlyDecreasingDistance = None
event_2HighlyDecreasingDistance = None
PAST_event_3CarA = None
event_3CarA = None
PAST_event_3CarB = None
event_3CarB = None
PAST_event_4CarA = None
event_4CarA = None
PAST_event_4CarB = None
event_4CarB = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'startMeasure,\w+,\w+,-?\d+,-?\d+,-?\d+,-?\d+$',inp) : 
			event_1carA = str(list_inp[1])
			event_1carB = str(list_inp[2])
			event_1xA = int(list_inp[3])
			event_1yA = int(list_inp[4])
			event_1xB = int(list_inp[5])
			event_1yB = int(list_inp[6])
			PAST_global_Distance = global_Distance
			global_Distance = float((abs(event_1xA - event_1xB) **  2  + abs(event_1yA - event_1yB) **  2 ) **  0.5 )
			event_1CarA = str(event_1carA)
			event_1CarB = str(event_1carB)
			PAST_event_1carA = event_1carA
			PAST_event_1carB = event_1carB
			PAST_event_1xA = event_1xA
			PAST_event_1yA = event_1yA
			PAST_event_1xB = event_1xB
			PAST_event_1yB = event_1yB
			PAST_event_1CarA = event_1CarA
			PAST_event_1CarB = event_1CarB
			print('start' + ',' +  str(event_1CarA) + ',' +  str(event_1carB))


		if re.match(r'location,\w+,\w+,-?\d+,-?\d+,-?\d+,-?\d+$',inp) : 
			event_2carA = str(list_inp[1])
			event_2carB = str(list_inp[2])
			event_2xA = int(list_inp[3])
			event_2yA = int(list_inp[4])
			event_2xB = int(list_inp[5])
			event_2yB = int(list_inp[6])
			PAST_global_Distance = global_Distance
			global_Distance = float((abs(event_2xA - event_2xB) **  2  + abs(event_2yA - event_2yB) **  2 ) **  0.5 )
			event_2HighlyDecreasingDistance = bool((global_Distance - ite(PAST_global_Distance!= None,PAST_global_Distance,global_Distance)) <  -5 )
			event_2CarA = str(event_2carA)
			event_2CarB = str(event_2carB)
			PAST_event_2carA = event_2carA
			PAST_event_2carB = event_2carB
			PAST_event_2xA = event_2xA
			PAST_event_2yA = event_2yA
			PAST_event_2xB = event_2xB
			PAST_event_2yB = event_2yB
			PAST_event_2CarA = event_2CarA
			PAST_event_2CarB = event_2CarB
			PAST_event_2HighlyDecreasingDistance = event_2HighlyDecreasingDistance
			print('distance' + ',' +  str(event_2CarA) + ',' +  str(event_2CarB) + ',' +  str(event_2HighlyDecreasingDistance))


		if re.match(r'startMonitoring,\w+,\w+$',inp) : 
			event_3carA = str(list_inp[1])
			event_3carB = str(list_inp[2])
			event_3CarA = str(event_3carA)
			event_3CarB = str(event_3carB)
			PAST_event_3carA = event_3carA
			PAST_event_3carB = event_3carB
			PAST_event_3CarA = event_3CarA
			PAST_event_3CarB = event_3CarB
			print('startMonitoring' + ',' +  str(event_3CarA) + ',' +  str(event_3CarB))


		if re.match(r'stopMonitoring,\w+,\w+$',inp) : 
			event_4carA = str(list_inp[1])
			event_4carB = str(list_inp[2])
			event_4CarA = str(event_4carA)
			event_4CarB = str(event_4carB)
			PAST_event_4carA = event_4carA
			PAST_event_4carB = event_4carB
			PAST_event_4CarA = event_4CarA
			PAST_event_4CarB = event_4CarB
			print('stopMonitoring' + ',' +  str(event_4CarA) + ',' +  str(event_4CarB))

