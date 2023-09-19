#.PHONY: build

#build:
#	sam build

.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get_weather src/get_weather/get_weather.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

build-sam:
	sam build

debug-with-sam: build-sam
	sam local --no-event | sam local invoke -d 9999 \
	--debug-args="-delveAPI=2" \
	--debugger-path $${HOME}/go/bin/linux_amd64  MySQSQueueFunction --event -

debug:
	sam local invoke GetWeather -d 8099 --debugger-path ./delve/ --debug-args "-delveAPI=2" --skip-pull-image