import re
def ite(condition, b1, b2): 
	return b1 if condition else b2

PAST_event_1ac = None
PAST_event_1temp = None
PAST_event_2ac = None
PAST_event_3ac = None
PAST_event_1InBound = None
event_1InBound = None
PAST_event_1Ac = None
event_1Ac = None
PAST_event_1Temp = None
event_1Temp = None
PAST_event_2Ac = None
event_2Ac = None
PAST_event_3Ac = None
event_3Ac = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'set,\w+,-?\d+$',inp):
			_EVENT_NAME = str(list_inp[0])
			_EVENT_PARAMS = list_inp[1:]
			event_1ac = str(list_inp[1])
			event_1temp = int(list_inp[2])
			event_1InBound = bool(event_1temp >=  17  and event_1temp <=  26 )
			event_1Ac = str(event_1ac)
			event_1Temp = int(event_1temp)
			PAST_event_1ac = event_1ac
			PAST_event_1temp = event_1temp
			PAST_event_1InBound = event_1InBound
			PAST_event_1Ac = event_1Ac
			PAST_event_1Temp = event_1Temp
			print('set' + ',' +  str(event_1Ac) + ',' +  str(event_1Temp) + ',' +  str(event_1InBound))


		if re.match(r'turn_on,\w+$',inp):
			_EVENT_NAME = str(list_inp[0])
			_EVENT_PARAMS = list_inp[1:]
			event_2ac = str(list_inp[1])
			event_2Ac = str(event_2ac)
			PAST_event_2ac = event_2ac
			PAST_event_2Ac = event_2Ac
			print('turn_on' + ',' +  str(event_2Ac))


		if re.match(r'turn_off,\w+$',inp):
			_EVENT_NAME = str(list_inp[0])
			_EVENT_PARAMS = list_inp[1:]
			event_3ac = str(list_inp[1])
			event_3Ac = str(event_3ac)
			PAST_event_3ac = event_3ac
			PAST_event_3Ac = event_3Ac
			print('turn_off' + ',' +  str(event_3Ac))

