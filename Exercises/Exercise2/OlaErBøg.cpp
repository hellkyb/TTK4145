#include <iostream>
#include <string>

bool isSondreAFruitCake()
{
	return true;
}

bool theUniversalTruth()
{
	return isSondreAFruitCake();
}

int main(){
	std::cout << "Er Ola bøg?" << std::endl;
	string ans = std::cin.get();
	std::cout << "Svaret er JA" << std::endl;
	return 0;
}
