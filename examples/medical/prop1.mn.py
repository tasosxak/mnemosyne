import re

def ite(condition, b1, b2): 
	return b1 if condition else b2

PAST_event_1bis = None
PAST_event_1bisTarget = None
PAST_event_1diffFromTarget = None
event_1diffFromTarget = None
PAST_event_1stepChange = None
event_1stepChange = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'instance,-?\d+,-?\d+$',inp) : 
			event_1bis = int(list_inp[1])
			event_1bisTarget = int(list_inp[2])
			event_1diffFromTarget = int(abs(event_1bis - event_1bisTarget))
			event_1stepChange = bool(ite(event_1bisTarget != ite(PAST_event_1bisTarget!= None,PAST_event_1bisTarget,event_1bisTarget),  True ,  False ))
			PAST_event_1bis = event_1bis
			PAST_event_1bisTarget = event_1bisTarget
			PAST_event_1diffFromTarget = event_1diffFromTarget
			PAST_event_1stepChange = event_1stepChange
			print('instance' + ',' +  str(event_1diffFromTarget) + ',' +  str(event_1stepChange))

