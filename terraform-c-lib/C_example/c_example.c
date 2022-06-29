#include <stdio.h>
#include <stdlib.h>
#include <terraform-c-lib.h>


#ifndef MAX_BUF
#define MAX_BUF 200
#endif


int main() {
	struct LoadModule_return module;

	char path[MAX_BUF];
	getcwd(path, MAX_BUF);
	printf("Current working directory: %s\n", path);

	char filePath[] = "../test_tf_files/";

	GoString filePath_gostring = {p: filePath, n: sizeof(filePath)}; 

	module = LoadModule(filePath_gostring);
	if (module.r0 != 0) {
		printf("Error: %s\n", module.r2);
		exit(1);
	}

	printf("Module: %s\n", module.r1);
	return 0;
}
