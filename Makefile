.PHONY: api

# generate protobuf api go code
api:
	cd conf && buf generate