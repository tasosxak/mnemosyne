import re
def ite(condition, b1, b2): 
	return b1 if condition else b2

global_Trace = ""
PAST_global_Trace = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'\w+$',inp):
			_EVENT_NAME = str(list_inp[0])
			PAST_global_Trace = global_Trace
			global_Trace = str(global_Trace + _EVENT_NAME)
			if (not re.match(r"(.*(pqp)p+$)+", global_Trace)) :
				print(_EVENT_NAME + ',' +  str(global_Trace))

