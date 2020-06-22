# Micro Gen
Cli tool for go micro, used to generate empty micro service.

## Pkger
```
pkger -include /cmd/new/templates -o cmd/new/pkger
```

## Installation
```
go get github.com/x-punch/micro-gen
```

## Usage
```
micro-gen new --name demo-service --namespace github.com/demo --path src
```