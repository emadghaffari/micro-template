# middleware
in this directory we check middlewares for http,grpc calls <br>
in your config file [service.router] you can use middlewares and add custom middleware to your endpoint

## Middleware UnaryInterceptor
in this method we check middlewares for unary requests 

## Middleware StreamInterceptor
in this method we check middlewares for streaming requests 

## example
* check middlewares for {METHOD_NAME} endpoint
- router:
    - method: {METHOD_NAME}
      description: "some description for method"
      maxAllowedAnomaly: 50
      middlewares:
        - CheckSome
        - checkSome
        - middleware3


### new method
1- create a method [http, unary, streaming] type <br>
2- use method in your config
