.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get_weather src/get_weather/get_weather.go

deploy: clean build
	sls deploy --verbose

build-sam:
	sam build

debug-with-sam: build-sam
	sam local invoke GetWeather -d 8099 \
	--debug-args="-delveAPI=2" \
	--debugger-path ./delve

debug:
	sam local invoke GetWeather -d 8099 --debugger-path ./delve/ --debug-args "-delveAPI=2"

debug-skip:
	sam local invoke GetWeather -d 8099 --debugger-path ./delve/ --debug-args "-delveAPI=2" --skip-pull-image

debug-event:
	echo "{"city":"Cape Town"}" | \
	sam local invoke GetWeather -d 8099 --debugger-path ./delve/ --debug-args "-delveAPI=2" --event - --skip-pull-image