terraform-c-lib: terraform-c-lib.go
	go build -buildmode=c-archive -o terraform-c-lib.a terraform-c-lib.go

clean:
	rm terraform-c-lib.a terraform-c-lib.h
