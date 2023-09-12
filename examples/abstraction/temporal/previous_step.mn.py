import re
def ite(condition, b1, b2): 
	return b1 if condition else b2

global_Counter =  0 
PAST_global_Counter = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'p$',inp):
			PAST_global_Counter = global_Counter
			global_Counter = int(global_Counter +  1 )
			if global_Counter <=  2  :
				print('p' + '')


		if re.match(r'q$',inp):
			PAST_global_Counter = global_Counter
			global_Counter = int( 0 )
			print('q' + '')

