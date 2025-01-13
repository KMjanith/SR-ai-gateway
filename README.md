1. install dependencies
```
go mod tidy
```

2. Compile the protobuf file
   ```
   protoc --go_out=. --go_opt=paths=source_relative spec/apiMessages.proto
   ```
3. To use this project you need to follow the instructions in [here](https://github.com/KMjanith/SR-service-runner/blob/main/Readme.md).
4. This service is workign as a Api-Gateway to pass http request and RPC msgs between front-end, auth and sorting services.
