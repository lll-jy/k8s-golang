# Guide

## Start program

1. In the `kubernetes-trial` directory, run `go build`.
1. Run `./k8s-trial` to start running the program.

## Instructions after launching the program

In the line asking for `Task (view, create, or delete): `, type in a task. Valid tasks are `view`, `create`, `delete`, 
and `exit`. Then follow the tips as provided in the stdout to provide further input.

## Unit tests

In `kubernetes-trial` directory, run `go test`, and wait. It will take about 2.5 minutes to finish the tests. Do make 
sure that before running the tests, there is no deployment called "kubernetes-bootcamp" running (if so, run `kubectl 
delete deployment kubernetes-bootcamp` and wait for a while until `kubectl get pod` would show no pods running 
deployment `kubernetes-bootcamp`.