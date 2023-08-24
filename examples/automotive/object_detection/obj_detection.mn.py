import re

def ite(condition, b1, b2): 
	return b1 if condition else b2

global_SysErr =  0 
PAST_global_SysErr = None
global_X =  -1 
PAST_global_X = None
global_Y =  -1 
PAST_global_Y = None
global_Z =  -1 
PAST_global_Z = None
global_EPO =  10 
PAST_global_EPO = None
global_EMAX =  20 
PAST_global_EMAX = None
global_Predicted =  False 
PAST_global_Predicted = None
PAST_event_1ru = None
PAST_event_2ru = None
PAST_event_2x = None
PAST_event_2y = None
PAST_event_2z = None
PAST_event_3ru = None
PAST_event_3x = None
PAST_event_3y = None
PAST_event_3z = None
PAST_event_4ru = None
PAST_event_4lerr = None
PAST_event_1outRu = None
event_1outRu = None
PAST_event_2outRu = None
event_2outRu = None
PAST_event_3error = None
event_3error = None
PAST_event_3LErr = None
event_3LErr = None
PAST_event_3outRu = None
event_3outRu = None
PAST_event_4outRu = None
event_4outRu = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'entry,\w+$',inp) : 
			event_1ru = str(list_inp[1])
			event_1outRu = str(event_1ru)
			PAST_event_1ru = event_1ru
			PAST_event_1outRu = event_1outRu
			print('entry' + ',' +  str(event_1outRu))


		if re.match(r'mk_prediction,\w+,-?\d+,-?\d+,-?\d+$',inp) : 
			event_2ru = str(list_inp[1])
			event_2x = int(list_inp[2])
			event_2y = int(list_inp[3])
			event_2z = int(list_inp[4])
			PAST_global_X = global_X
			global_X = int(event_2x)
			PAST_global_Y = global_Y
			global_Y = int(event_2y)
			PAST_global_Z = global_Z
			global_Z = int(event_2z)
			PAST_global_Predicted = global_Predicted
			global_Predicted = bool( True )
			event_2outRu = str(event_2ru)
			PAST_event_2ru = event_2ru
			PAST_event_2x = event_2x
			PAST_event_2y = event_2y
			PAST_event_2z = event_2z
			PAST_event_2outRu = event_2outRu
			print('predicted' + ',' +  str(event_2outRu))


		if re.match(r'obstacle,\w+,-?\d+,-?\d+,-?\d+$',inp) : 
			event_3ru = str(list_inp[1])
			event_3x = int(list_inp[2])
			event_3y = int(list_inp[3])
			event_3z = int(list_inp[4])
			event_3LErr = int(ite(global_Predicted, abs(event_3x - global_X) + abs(event_3y - global_Y) + abs(event_3z - global_Z),  0 ))
			PAST_global_SysErr = global_SysErr
			global_SysErr = int(global_SysErr + event_3LErr)
			PAST_global_Predicted = global_Predicted
			global_Predicted = bool( False )
			event_3error = bool((event_3LErr > global_EPO) or (global_SysErr > global_EMAX))
			event_3outRu = str(event_3ru)
			PAST_event_3ru = event_3ru
			PAST_event_3x = event_3x
			PAST_event_3y = event_3y
			PAST_event_3z = event_3z
			PAST_event_3error = event_3error
			PAST_event_3LErr = event_3LErr
			PAST_event_3outRu = event_3outRu
			print('valid' + ',' +  str(event_3outRu) + ',' +  str(event_3error))


		if re.match(r'exit,\w+,-?\d+$',inp) : 
			event_4ru = str(list_inp[1])
			event_4lerr = int(list_inp[2])
			event_4outRu = str(event_4ru)
			PAST_global_SysErr = global_SysErr
			global_SysErr = int(global_SysErr - event_4lerr)
			PAST_event_4ru = event_4ru
			PAST_event_4lerr = event_4lerr
			PAST_event_4outRu = event_4outRu
			print('exit' + ',' +  str(event_4ru))

