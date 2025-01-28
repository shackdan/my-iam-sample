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

	// The code that defines your stack goes here

	// example resource
	// queue := awssqs.NewQueue(stack, jsii.String("MyIamTestQueue"), &awssqs.QueueProps{
	// 	VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	// })

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

	myPolicy := iam.NewManagedPolicy(stack, jsii.String("my-cdk-policy"), &iam.ManagedPolicyProps{
		Document:          policyDocument,
		ManagedPolicyName: jsii.String("my-cdk-managed-policy"),
	})

	myrole := iam.NewRole(stack, jsii.String("my-cdk-role"), &iam.RoleProps{
		RoleName:       jsii.String("my-cdk-role"),
		AssumedBy:      iam.NewAccountPrincipal(jsii.String(*iam.NewAccountRootPrincipal().PrincipalAccount())),
		InlinePolicies: &map[string]iam.PolicyDocument{"MyPolicy": policyDocument},
	})

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

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
