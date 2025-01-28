package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type MyIamSampleStackProps struct {
	awscdk.StackProps
}

func NewMyIamSampleStack(scope constructs.Construct, id string, props *MyIamSampleStackProps) awscdk.Stack {

	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Add Tags for Stack
	awscdk.Tags_Of(stack).Add(jsii.String("Environment"), jsii.String("dev"), &awscdk.TagProps{})
	awscdk.Tags_Of(stack).Add(jsii.String("Owner"), jsii.String("dannyboy"), &awscdk.TagProps{})

	// Create policy document
	policyDocument := awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions:   &[]*string{jsii.String("s3:ListBucket")},
				Resources: &[]*string{jsii.String("arn:aws:s3:::example-bucket")},
			}),
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions:   &[]*string{jsii.String("s3:GetObject")},
				Resources: &[]*string{jsii.String("arn:aws:s3:::example-bucket/*")},
			}),
		},
	})

	// Create policy and reference policy document
	myPolicy := iam.NewManagedPolicy(stack, jsii.String("my-cdk-policy"), &iam.ManagedPolicyProps{
		Document:          policyDocument,
		ManagedPolicyName: jsii.String("my-cdk-managed-policy"),
	})

	// Create role set trust to root of the current account, and add inline policy
	myrole := iam.NewRole(stack, jsii.String("my-cdk-role"), &iam.RoleProps{
		RoleName:       jsii.String("my-cdk-role"),
		AssumedBy:      iam.NewAccountPrincipal(jsii.String(*iam.NewAccountRootPrincipal().PrincipalAccount())),
		InlinePolicies: &map[string]iam.PolicyDocument{"MyPolicy": policyDocument},
	})

	//Add policies to new role (Managed and in-line)
	myrole.AddManagedPolicy(awsiam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonS3ReadOnlyAccess")))
	myrole.AddToPolicy(iam.NewPolicyStatement(&iam.PolicyStatementProps{ //add "in-line policy" NO option to name the given policy
		Actions: &[]*string{
			jsii.String("s3:ListBucket"),
			jsii.String("s3:GetObject"),
		},
		Resources: &[]*string{
			jsii.String("arn:aws:s3:::my-bucket"),
			jsii.String("arn:aws:s3:::my-bucket/*"),
		},
		Effect: awsiam.Effect_ALLOW,
	}))
	myrole.AddManagedPolicy(myPolicy)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewMyIamSampleStack(app, "MyIamSampleStack", &MyIamSampleStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {

	return nil

}
