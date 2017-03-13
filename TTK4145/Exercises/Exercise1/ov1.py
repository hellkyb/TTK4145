from threading import Thread
i = 0

def threadFunction():
	global i
	print("Thread")
	for n in range(10000):
		i = i + 1

def threadFunctionTwo():
	global i
	print("Thread2")
	for n in range(10000):
		i = i - 1
	
def main():
	someThread = Thread(target = threadFunction, args = (),)
	someThread.start()
	someThread.join()

	someThreadTwo = Thread(target = threadFunctionTwo, args = (),)
	someThreadTwo.start()
	someThreadTwo.join()
	print("Main")
	print(i)

main()
