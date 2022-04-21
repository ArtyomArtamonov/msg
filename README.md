# Msg

gRPC server for messaging written in golang.


## How to use 

This version supports only one receiver. Its id is set to 0.

1. Open a streaming connection using GetMessages() method.
2. Send a message using SendMessage() method with id=0
