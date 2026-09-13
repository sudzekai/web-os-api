package abstractions

type IExecutor interface {
	Execute(command string, args ...string)
	ExecuteInDirectory(directoryPath string, command string, args ...string)
}
