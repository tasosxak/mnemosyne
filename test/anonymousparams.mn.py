import re
def ite(condition, b1, b2): 
	return b1 if condition else b2

PAST_event_1y = None
event_1y = None
PAST_event_1t = None
event_1t = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'\w+$',inp):
			_EVENT_NAME = str(list_inp[0])
			_EVENT_PARAMS = list_inp[1:]
			event_1t = int("Gg" +  5 )
			PAST_event_1y = event_1y
			PAST_event_1t = event_1t
			print('x' + ',' +  str(event_1y))

