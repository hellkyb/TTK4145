#include <stdio.h>
#include <pthread.h>

int i = 0;

void* someThread()
{
	for(int j = 0; j < 10000; j++)
	{
		i++;
	}

	return NULL;
}

void* someThreadTwo()
{
	for(int j=0; j < 10000; j++)
	{
		i--;
	}
	return NULL;
}

int main()
{
	pthread_t someOtherThread;
	pthread_create(&someOtherThread, NULL, someThread, NULL);
	
	pthread_join(someOtherThread, NULL);
	
	pthread_t someOtherThread2;
	pthread_create(&someOtherThread2, NULL, someThreadTwo, NULL);
	
	pthread_join(someOtherThread2, NULL);
	printf("%d\n", i);

	return 0;
}
