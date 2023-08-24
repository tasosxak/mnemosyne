import re

def ite(condition, b1, b2): 
	return b1 if condition else b2

global_Temp =  0 
PAST_global_Temp = None
PAST_event_1Car = None
PAST_event_1temp = None
PAST_event_2Car = None
PAST_event_2temp = None
PAST_event_1outCar = None
event_1outCar = None
PAST_event_2Warming = None
event_2Warming = None
PAST_event_2outCar = None
event_2outCar = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'StartMeasure,\w+,\d+$',inp) : 
			event_1Car = str(list_inp[1])
			event_1temp = int(list_inp[2])
			event_1outCar = event_1Car
			global_Temp = event_1temp
			PAST_event_1Car = event_1Car
			PAST_event_1temp = event_1temp
			PAST_event_1outCar = event_1outCar
			PAST_global_Temp = global_Temp
			print('test' + ',' +  str(event_1outCar))


		if re.match(r'EndMeasure,\w+,\d+$',inp) : 
			event_2Car = str(list_inp[1])
			event_2temp = int(list_inp[2])
			global_Temp = event_2temp
			event_2Warming = (global_Temp - ite(PAST_global_Temp!= None,PAST_global_Temp,global_Temp)) >  5 
			event_2outCar = event_2Car
			PAST_event_2Car = event_2Car
			PAST_event_2temp = event_2temp
			PAST_event_2Warming = event_2Warming
			PAST_event_2outCar = event_2outCar
			PAST_global_Temp = global_Temp
			print('increase' + ',' +  str(event_2outCar) + ',' +  str(event_2Warming))

