#include <stdio.h>
#include <pthread.h>

int i = 0;
pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;

void* someThread()
{
	for(int j = 0; j < 1000000; j++)
	{
		pthread_mutex_lock(&mutex);
		i++;
		pthread_mutex_unlock(&mutex);
	}

	return NULL;
}

void* someThreadTwo()
{
	for(int j=0; j < 1000001; j++)
	{
		pthread_mutex_lock(&mutex);
		i--;
		pthread_mutex_unlock(&mutex);
	}
	return NULL;
}

int main()
{
	

	pthread_t someOtherThread;
	pthread_t someOtherThread2;
	pthread_create(&someOtherThread, NULL, someThread, NULL);
	pthread_create(&someOtherThread2, NULL, someThreadTwo, NULL);
	
	pthread_join(someOtherThread, NULL);
	pthread_join(someOtherThread2, NULL);
	printf("%d\n", i);

	return 0;
}

// Semaphores or Mutex ? - Mutex
