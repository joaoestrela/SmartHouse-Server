APP = gateway
USER = ubuntu
KEY = ../pub/ce_key
REMOTE = /home/ubuntu

release: linux64 vars
	ssh -i $(KEY) $(HOST) "rm -f $(REMOTE)/$(APP)"
	scp -i $(KEY) $(APP) $(HOST):$(REMOTE)
	
	# Redirect stdout and stderr so it doesn't hang
	# https://stackoverflow.com/a/29172
	ssh -i $(KEY) $(HOST) "nohup ./$(APP) > st.out 2> st.err < /dev/null &"

vars:
	$(eval PUBLIC_DNS = $(shell terraform output public_dns))
	$(eval HOST = $(USER)@$(PUBLIC_DNS))

linux64:
	env GOOS=linux GOARCH=amd64 go build -o $(APP) -i ../.

build:
	env GOOS=linux GOARCH=amd64 go build -o bin/$(APP)-linux -i .
	env GOOS=darwin GOARCH=amd64 go build -o bin/$(APP)-darwin -i .
	env GOOS=windows GOARCH=amd64 go build -o bin/gateway.exe -i .

deploy:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o bin/gateway-arm -i .; scp bin/gateway-arm pi@raspberrypi.local:~/

clean:
	rm $(APP)