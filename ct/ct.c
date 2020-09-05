
#include<stdio.h>
#define bool char
#define true 1
#define false 0
typedef struct Info{
    int vm;
    float pfi[512];
    int nspkId;
    int nspknum;
    int nb;
    int ntb ;
    char spkname[128];
    bool update;

}INFO;


int main(){
    printf("%d",sizeof(INFO));
    printf("\n%d,%d",true,false);
}