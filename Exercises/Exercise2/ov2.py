from threading import Thread
i = 0

def threadFunction():
	global i
	for n in range(1000000):
		i = i + 1

def threadFunctionTwo():
	global i
	for n in range(1000001):
		i = i - 1
	
def main():
	someThread = Thread(target = threadFunction, args = (),)
	someThreadTwo = Thread(target = threadFunctionTwo, args = (),)
	someThread.start()
	someThreadTwo.start()
	someThread.join()
	someThreadTwo.join()
	print(i)

main()
