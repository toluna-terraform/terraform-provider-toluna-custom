package provider

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceStartCodebuild() *schema.Resource {
	return &schema.Resource{
		Create: resourceStartCodebuildCreate,
		Read:   resourceStartCodebuildRead,
		Update: resourceStartCodebuildUpdate,
		Delete: resourceStartCodebuildDelete,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"payload": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func startCodebuild(d *schema.ResourceData, m interface{}, action string) (str string, er error) {
	var profile = d.Get("aws_profile").(string)
	if profile == "" {
		profile = "default"
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
	}))

	client := codebuild.New(sess, &aws.Config{Region: aws.String(d.Get("region").(string))})
	envVar := []*codebuild.EnvironmentVariable{}
	param := &codebuild.EnvironmentVariable{
		Name:  aws.String("action"),
		Type:  aws.String("PLAINTEXT"),
		Value: aws.String(action)}
	envVar = append(envVar, param)
	input := &codebuild.StartBuildInput{
		ProjectName:                  aws.String(d.Get("project_name").(string)),
		EnvironmentVariablesOverride: envVar,
	}
	result, err := client.StartBuild(input)
	if err != nil {
		return "", fmt.Errorf("Error calling build project", err)
	}
	Ids := []*string{}
	Ids = append(Ids, result.Build.Id)
	for {
		buildstatus, err := client.BatchGetBuilds(&codebuild.BatchGetBuildsInput{Ids: Ids})
		if err != nil {
			break
		}
		if *buildstatus.Builds[0].BuildComplete {
			break
		}
	}

	buildresult, err := client.BatchGetBuilds(&codebuild.BatchGetBuildsInput{Ids: Ids})
	if err != nil {
		return "", fmt.Errorf("Error failed getting build:%s", *buildresult.Builds[0].Id)
	}

	if *buildresult.Builds[0].BuildStatus != "SUCCEEDED" {
		return "", fmt.Errorf("Error build status:%s,%s", *buildresult.Builds[0].BuildStatus, *buildresult.Builds[0].Id)
	}
	return *buildresult.Builds[0].Arn, nil
}

func resourceStartCodebuildCreate(d *schema.ResourceData, m interface{}) error {
	result, err := startCodebuild(d, m, "apply")
	if err != nil {
		d.Partial(true)
		return err
	}
	d.SetId(result)
	return resourceStartCodebuildRead(d, m)
}

func resourceStartCodebuildRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceStartCodebuildUpdate(d *schema.ResourceData, m interface{}) error {
	result, err := startCodebuild(d, m, "apply")
	if err != nil {
		d.Partial(true)
		return err
	}
	d.SetId(result)
	return resourceStartCodebuildRead(d, m)
}

func resourceStartCodebuildDelete(d *schema.ResourceData, m interface{}) error {
	result, err := startCodebuild(d, m, "destroy")
	if err != nil {
		d.Partial(true)
		return err
	}
	d.SetId(result)
	return resourceStartCodebuildRead(d, m)
}
