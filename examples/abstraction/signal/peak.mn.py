import re
def ite(condition, b1, b2): 
	return b1 if condition else b2

global_Signal =  0 
PAST_global_Signal = None
global_Deff =  True 
PAST_global_Deff = None
global_Peak =  False 
PAST_global_Peak = None
global_Half =  False 
PAST_global_Half = None
global_TimeHalf =  -1 
PAST_global_TimeHalf = None
PAST_event_1val = None
PAST_event_1timestamp = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'\w+,-?\d+(\.\d*)?,-?\d+$',inp):
			_EVENT_NAME = str(list_inp[0])
			event_1val = float(list_inp[1])
			event_1timestamp = int(list_inp[2])
			PAST_global_TimeHalf = global_TimeHalf
			global_TimeHalf = int(ite(global_Peak,  -1 , global_TimeHalf))
			PAST_global_Deff = global_Deff
			global_Deff = bool(ite(global_Peak,  True , global_Deff))
			PAST_global_Half = global_Half
			global_Half = bool(ite(global_Peak,  False , global_Half))
			PAST_global_Peak = global_Peak
			global_Peak = bool(ite(global_Peak,  False , global_Peak))
			PAST_global_Signal = global_Signal
			global_Signal = float(event_1val)
			PAST_global_Deff = global_Deff
			global_Deff = bool(ite(((global_Signal - ite(PAST_global_Signal!= None,PAST_global_Signal,global_Signal)) >=  0 ),  True ,  False ))
			PAST_global_Half = global_Half
			global_Half = bool(global_Deff != ite(PAST_global_Deff!= None,PAST_global_Deff,global_Deff))
			PAST_global_TimeHalf = global_TimeHalf
			global_TimeHalf = int(ite(global_Half and global_TimeHalf ==  -1 , ite(PAST_event_1timestamp!= None,PAST_event_1timestamp,event_1timestamp), global_TimeHalf))
			PAST_global_Peak = global_Peak
			global_Peak = bool(global_Half and global_Deff)
			PAST_event_1val = event_1val
			PAST_event_1timestamp = event_1timestamp
			if global_Peak :
				print('peak' + ',' +  str(global_TimeHalf) + ',' +  str(event_1timestamp))

