package provider

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceInvokeLambda() *schema.Resource {
	return &schema.Resource{
		Create: resourceInvokeLambdaCreate,
		Read:   resourceInvokeLambdaRead,
		Update: resourceInvokeLambdaUpdate,
		Delete: resourceInvokeLambdaDelete,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_profile": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"function_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"payload": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type getItemsRequest struct {
	SortBy     string
	SortOrder  string
	ItemsToGet int
}

type getItemsResponseError struct {
	Message string `json:"message"`
}

type getItemsResponseData struct {
	Item string `json:"item"`
}

type getItemsResponseBody struct {
	Result string                 `json:"result"`
	Data   []getItemsResponseData `json:"data"`
	Error  getItemsResponseError  `json:"error"`
}

type getItemsResponseHeaders struct {
	ContentType string `json:"Content-Type"`
}

type getItemsResponse struct {
	StatusCode int                     `json:"statusCode"`
	Headers    getItemsResponseHeaders `json:"headers"`
	Body       getItemsResponseBody    `json:"body"`
}

func invokeLambda(d *schema.ResourceData, m interface{}, action string) (str string, er error) {
	var profile = d.Get("aws_profile").(string)
	if profile == "" {
		profile = "default"
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
	}))
	client := lambda.New(sess, &aws.Config{Region: aws.String(d.Get("region").(string))})

	var j map[string]interface{}
	var p = []byte(d.Get("payload").(string))
	err := json.Unmarshal(p, &j)
	j["action"] = action
	payload, err := json.Marshal(j)
	if err != nil {
		return "", fmt.Errorf("Error getting items, StatusCode: ", err)
	}
	result, err := client.Invoke(&lambda.InvokeInput{FunctionName: aws.String(d.Get("function_name").(string)), Payload: payload})
	if err != nil {
		return "", fmt.Errorf("Error calling lambda function", err)
	}
	var resp getItemsResponse
	err = json.Unmarshal(result.Payload, &resp)
	if err != nil {
		return "", fmt.Errorf("Error unmarshalling response", err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error getting items, StatusCode: " + strconv.Itoa(resp.StatusCode))
	}

	if resp.Body.Result == "failure" {
		return "", fmt.Errorf("Error Failed to get items", err)
	}

	if len(resp.Body.Data) > 0 {
		for i := range resp.Body.Data {
			fmt.Println(resp.Body.Data[i].Item)
		}
	} else {
		fmt.Println("There were no items")
	}
	return strconv.Itoa(resp.StatusCode), nil
}

func resourceInvokeLambdaCreate(d *schema.ResourceData, m interface{}) error {
	result, err := invokeLambda(d, m, "apply")
	if err != nil {
		d.Partial(true)
		return err
	}
	d.SetId(result)
	return resourceInvokeLambdaRead(d, m)
}

func resourceInvokeLambdaRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceInvokeLambdaUpdate(d *schema.ResourceData, m interface{}) error {
	result, err := invokeLambda(d, m, "apply")
	if err != nil {
		d.Partial(true)
		return err
	}
	d.SetId(result)
	return resourceInvokeLambdaRead(d, m)
}

func resourceInvokeLambdaDelete(d *schema.ResourceData, m interface{}) error {
	result, err := invokeLambda(d, m, "destroy")
	if err != nil {
		d.Partial(true)
		return err
	}
	d.SetId(result)
	return resourceInvokeLambdaRead(d, m)
}
