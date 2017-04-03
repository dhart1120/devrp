# devrp #

## Summary ##
devrp is a really simple reverse proxy that can be used in development. 

## Usage ##
Run the binary with the `-p` argument to forward from `src:dest` 

This can be used to alias a port or for funneling traffic from a number of ports to a single destination port.

*Example:* `devrp -p 8080:80,8081:80,8082:80`

Will forward traffic from ports 8080, 8081, and 8082 to port 80
