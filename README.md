# Project Name

## Description
Provide a brief description of your project here.

## Installation
Describe how to install and set up the project.

## Usage
Explain how to use the project and provide examples.

docker run --rm --name jaeger \
  -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -p 14268:14268 \
  -p 14269:14269 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.54

go run .
## Contributing
Specify how others can contribute to the project.

## License
Indicate the license under which the project is distributed.

## Contact
Provide contact information for any inquiries or support.
