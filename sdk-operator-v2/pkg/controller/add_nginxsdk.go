package controller

import (
	"github.com/arbhoj/ansible-operator/pkg/controller/nginxsdk"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, nginxsdk.Add)
}
