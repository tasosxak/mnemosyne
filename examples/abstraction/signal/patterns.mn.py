import re
def ite(condition, b1, b2): 
	return b1 if condition else b2

global_Signal =  0 
PAST_global_Signal = None
global_Deff1 =  True 
PAST_global_Deff1 = None
global_Deff2 =  False 
PAST_global_Deff2 = None
global_LocalMax =  False 
PAST_global_LocalMax = None
global_LocalMin =  False 
PAST_global_LocalMin = None
global_Change1 =  False 
PAST_global_Change1 = None
global_Change2 =  False 
PAST_global_Change2 = None
global_TimeChange1 =  -1 
PAST_global_TimeChange1 = None
global_TimeChange2 =  -1 
PAST_global_TimeChange2 = None
PAST_event_1val = None
PAST_event_1timestamp = None
PAST_event_2val = None
PAST_event_2timestamp = None
PAST_event_3val = None
PAST_event_3timestamp = None
PAST_event_2prevTime = None
event_2prevTime = None
PAST_event_3prevTime = None
event_3prevTime = None

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
			PAST_global_Signal = global_Signal
			global_Signal = float(event_1val)
			PAST_global_Deff1 = global_Deff1
			global_Deff1 = bool(ite(((global_Signal - ite(PAST_global_Signal!= None,PAST_global_Signal,global_Signal)) >=  0 ),  True ,  False ))
			PAST_global_Deff2 = global_Deff2
			global_Deff2 = bool(ite(((global_Signal - ite(PAST_global_Signal!= None,PAST_global_Signal,global_Signal)) >=  0 ),  True ,  False ))
			PAST_global_Change1 = global_Change1
			global_Change1 = bool(global_Deff1 != ite(PAST_global_Deff1!= None,PAST_global_Deff1,global_Deff1))
			PAST_global_Change2 = global_Change2
			global_Change2 = bool(global_Deff2 != ite(PAST_global_Deff2!= None,PAST_global_Deff2,global_Deff2))
			PAST_global_TimeChange1 = global_TimeChange1
			global_TimeChange1 = int(ite(global_Change1 and global_TimeChange1 ==  -1 , ite(PAST_event_1timestamp!= None,PAST_event_1timestamp,event_1timestamp), global_TimeChange1))
			PAST_global_TimeChange2 = global_TimeChange2
			global_TimeChange2 = int(ite(global_Change2 and global_TimeChange2 ==  -1 , ite(PAST_event_1timestamp!= None,PAST_event_1timestamp,event_1timestamp), global_TimeChange2))
			PAST_event_1val = event_1val
			PAST_event_1timestamp = event_1timestamp


		if re.match(r'\w+,-?\d+(\.\d*)?,-?\d+$',inp):
			_EVENT_NAME = str(list_inp[0])
			event_2val = float(list_inp[1])
			event_2timestamp = int(list_inp[2])
			PAST_global_TimeChange1 = global_TimeChange1
			global_TimeChange1 = int(ite(global_LocalMax,  -1 , global_TimeChange1))
			PAST_global_Change1 = global_Change1
			global_Change1 = bool(ite(global_LocalMax,  False , global_Change1))
			PAST_global_LocalMax = global_LocalMax
			global_LocalMax = bool(ite(global_LocalMax,  False , global_LocalMax))
			PAST_global_LocalMax = global_LocalMax
			global_LocalMax = bool(global_Change1 and global_Deff1)
			event_2prevTime = int(ite(global_LocalMax, ite(PAST_event_2timestamp!= None,PAST_event_2timestamp,event_2timestamp), event_2timestamp))
			PAST_event_2val = event_2val
			PAST_event_2timestamp = event_2timestamp
			PAST_event_2prevTime = event_2prevTime
			if global_LocalMax :
				print('peak' + ',' +  str(global_TimeChange1) + ',' +  str(event_2prevTime))


		if re.match(r'\w+,-?\d+(\.\d*)?,-?\d+$',inp):
			_EVENT_NAME = str(list_inp[0])
			event_3val = float(list_inp[1])
			event_3timestamp = int(list_inp[2])
			PAST_global_TimeChange2 = global_TimeChange2
			global_TimeChange2 = int(ite(global_LocalMin,  -1 , global_TimeChange2))
			PAST_global_Change2 = global_Change2
			global_Change2 = bool(ite(global_LocalMin,  False , global_Change2))
			PAST_global_LocalMin = global_LocalMin
			global_LocalMin = bool(ite(global_LocalMin,  False , global_LocalMin))
			PAST_global_LocalMin = global_LocalMin
			global_LocalMin = bool(global_Change2 and not global_Deff2)
			event_3prevTime = int(ite(global_LocalMin, ite(PAST_event_3timestamp!= None,PAST_event_3timestamp,event_3timestamp), event_3timestamp))
			PAST_event_3val = event_3val
			PAST_event_3timestamp = event_3timestamp
			PAST_event_3prevTime = event_3prevTime
			if global_LocalMin :
				print('fall' + ',' +  str(global_TimeChange2) + ',' +  str(event_3prevTime))

