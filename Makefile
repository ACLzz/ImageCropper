BIN_FOLDER=./bin
CONFIG_FOLDER=$(HOME)/.imcr

define shutdown_server
pkill -f image_cropper_server
endef

build:
	go build -o $(BIN_FOLDER)/image_cropper_server ./src/main

move_config:
	mkdir -p $(CONFIG_FOLDER)
	cp -rf extra/* $(CONFIG_FOLDER)

setup: move_config build

remove:
	rm -rf $(CONFIG_FOLDER)

tests: build
	cd bin;\
    	./image_cropper_server &
	go test -count=1 \
		$(shell find -name "*_test.go" -exec dirname {} \; | uniq) || $(call shutdown_server)

	$(call shutdown_server)