def ite(condition, b1, b2): 
	return b1 if condition else b2
PAST_event_0bis = None
PAST_event_0bisTarget = None
PAST_event_0diffFromTarget = None
event_0diffFromTarget = None
PAST_event_0stepChange = None
event_0stepChange = None
while True:
	inp = input()
	inp = inp.split(',') 
	if inp[0] == 'instance' : 
		event_0bis = int(int(inp[1])) 
		event_0bisTarget = int(int(inp[2])) 
		event_0diffFromTarget = abs(event_0bis - event_0bisTarget)
		event_0stepChange = ite(event_0bisTarget != ite(PAST_event_0bisTarget!= None,PAST_event_0bisTarget,event_0bisTarget),  True ,  False )
		PAST_event_0bis = event_0bis
		PAST_event_0bisTarget = event_0bisTarget
		PAST_event_0diffFromTarget = event_0diffFromTarget
		PAST_event_0stepChange = event_0stepChange
		print('instance' + ',' + str(event_0diffFromTarget) + ',' + str(event_0stepChange)) 
