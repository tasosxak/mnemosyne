import re
def ite(condition, b1, b2): 
	return b1 if condition else b2

global_E = ""
PAST_global_E = None

if __name__ == '__main__':
	while True:
		try:
			inp = input()
			list_inp = inp.split(',')
		except EOFError:
			exit()
		if re.match(r'q$',inp):
			PAST_global_E = global_E
			global_E = str("q")
			if ite(PAST_global_E!= None,PAST_global_E,global_E) != "q" :
				print('q' + '')


		if re.match(r'p$',inp):
			PAST_global_E = global_E
			global_E = str("p" + "s")
			print('p' + '')

