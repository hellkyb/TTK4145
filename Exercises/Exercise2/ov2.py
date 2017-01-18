from threading import Thread
import threading
i = 0
lock = threading.Lock()
def threadFunction():
	global i
	for n in range(100000):
		lock.acquire()
		i = i + 1
		lock.release()


def threadFunctionTwo():
	global i
	for n in range(100001):
		lock.acquire()
		i = i - 1
		lock.release()
	
def main():
	someThread = Thread(target = threadFunction, args = (),)
	someThreadTwo = Thread(target = threadFunctionTwo, args = (),)
	someThread.start()
	someThreadTwo.start()
	someThread.join()
	someThreadTwo.join()
	print(i)

main()
