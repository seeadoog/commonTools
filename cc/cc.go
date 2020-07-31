package main
/*

#include <stdio.h>
#include<stdlib.h>

void show(char* str){
	printf("%s",str);
}

typedef struct stu{
	int a;
	char* b;
};

void ss(stu* s){
	printf("%d,%s",s->a,s->b);
}

 */
import "C"
import "unsafe"

func main()  {
	str:=C.CString("hello world")
	C.show(str)
	C.free(unsafe.Pointer(str))
}