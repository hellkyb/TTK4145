#include <stdio.h>
#include <pthread.h>

void* someThread()
{
	printf("This is my design\n");
	return NULL;
}

int main()
{
	pthread_t someOtherThread;
	pthread_create(&someOtherThread, NULL, someThread, NULL);
	// Arguments

	pthread_join(someOtherThread, NULL);
	printf("Hello, from the other side\n");
	return 0;
}
