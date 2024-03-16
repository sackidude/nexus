# nexus

Data handling website for yeast project.

## About

This is for a yeast project I did at school analyzing carbon dioxid production under anaerobic circumstances. This website was built for data collection and displaying of the collected data in real time for my presentation.

The tech stack include: golang, htmx and mySQL server. Node.js was also use for some helper functions in data collecting and storing.

## Running and building

Check config/configuration.md for instructions

## Setup for rust

### Terminal 1 - To run the server.

cargo watch -q -c -w src/ -x "run"

### Terminal 2 - To run the tests.

cargo watch -q -c -w examples/ -x "run --example quick_dev"
