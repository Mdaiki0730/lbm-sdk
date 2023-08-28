all: rollupd-env

rollupd-env: rollupd-rmi
	docker build --tag cosmossdk/rollupd -f rollupd-env/Dockerfile \
    		$(shell git rev-parse --show-toplevel)

rollupd-dlv: rollupd-rmi
	docker build --tag cosmossdk/rollupd -f rollupd-dlv/Dockerfile \
        	$(shell git rev-parse --show-toplevel)

rollupd-rmi:
	docker rmi cosmossdk/rollupd 2>/dev/null; true

.PHONY: all rollupd-env rollupd-dlv rollupd-rmi