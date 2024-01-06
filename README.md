To create a project, you run go mod init and a name for a project, for example, “my-project”:

   go mod init my-project

If you want to update the go.mod and go.sum files to include the latest versions of the dependencies, you can use the following command:
  go get -u ./...
