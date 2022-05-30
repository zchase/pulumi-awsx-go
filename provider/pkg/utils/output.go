package utils

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

func ApplyAny[A any, O pulumi.Output](output pulumi.AnyOutput, block func(a A) O) O {
	return output.ApplyT(func(x interface{}) O {
		v := x.(A)
		return block(v)
	}).(O)
}

func ApplyAnyError[A any, O pulumi.Output](output pulumi.AnyOutput, block func(a A) (O, error)) O {
	return output.ApplyT(func(x interface{}) (O, error) {
		v := x.(A)
		return block(v)
	}).(O)
}
