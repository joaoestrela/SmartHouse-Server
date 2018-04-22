APP = gateway
USER = ubuntu
KEY = ../pub/ce_key
REMOTE = /home/ubuntu

release: linux64 vars
	ssh -i $(KEY) $(HOST) "rm -f $(REMOTE)/$(APP)"
	scp -i $(KEY) $(APP) $(HOST):$(REMOTE)
	ssh -i $(KEY) $(HOST) "nohup ./$(APP) &"

vars:
	$(eval PUBLIC_DNS = $(shell terraform output public_dns))
	$(eval HOST = $(USER)@$(PUBLIC_DNS))

linux64:
	env GOOS=linux GOARCH=amd64 go build -o $(APP) -i ../.

clean:
	rm $(APP)