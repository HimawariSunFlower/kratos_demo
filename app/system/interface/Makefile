include ../../../app_makefile



.PHONY: config
# generate internal proto
config:
	protoc --proto_path=. \
 	       --go_out=paths=source_relative:. \
	       $(INTERNAL_PROTO_FILES)