package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

func resourceJenkinsRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceJenkinsRoleCreate,
		Read:   resourceJenkinsRoleRead,
		Update: resourceJenkinsRoleUpdate,
		Delete: resourceJenkinsRoleDelete,
		Importer: &schema.ResourceImporter{
			State: resourceJenkinsRoleImport,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: "A structured role objects",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"global": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"item": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"node": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func initVariables(d *schema.ResourceData, m interface{}) (jenkins_url, jenkins_username, jenkins_api_token, user_id string, role_set *schema.Set) {
	config := m.(MyProviderConfig)

	jenkins_url = config.jenkins_url
	jenkins_username = config.jenkins_username
	jenkins_api_token = config.jenkins_api_token
	user_id = d.Get("user_id").(string)
	role_set = d.Get("role").(*schema.Set)

	return jenkins_url, jenkins_username, jenkins_api_token, user_id, role_set
}

func sendRequest(endpoint, method, jenkins_url, jenkins_username, jenkins_api_token, user_id, role_name, role_type string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("type", role_type)
	writer.WriteField("roleName", role_name)
	writer.WriteField("sid", user_id)
	writer.Close()

	var req *http.Request
	var err error

	// Set authentication
	req, err = http.NewRequest(method, jenkins_url+endpoint, body)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	req.SetBasicAuth(jenkins_username, jenkins_api_token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}

	// Do HTTP request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request was successful")
	} else {
		fmt.Println("Error:", resp.Status)
	}

	return nil
}

func resourceJenkinsRoleCreate(d *schema.ResourceData, m interface{}) error {
	jenkins_url, jenkins_username, jenkins_api_token, user_id, role_set := initVariables(d, m)

	roleList := role_set.List()
	endpoint := "/role-strategy/strategy/assignRole"
	method := "POST"

	for _, role := range roleList {
		roleData := role.(map[string]interface{})

		// Execute for Global Roles
		if global, ok := roleData["global"]; ok {
			globalRoles := global.([]interface{})
			role_type := "globalRoles"
			for _, globalRole := range globalRoles {
				role_name := globalRole.(string)
				if err := sendRequest(endpoint, method, jenkins_url, jenkins_username, jenkins_api_token, user_id, role_name, role_type); err != nil {
					return err
				}
			}
		}

		// Execute for Item Roles
		if item, ok := roleData["item"]; ok {
			itemRoles := item.([]interface{})
			role_type := "projectRoles"
			for _, itemRole := range itemRoles {
				role_name := itemRole.(string)
				if err := sendRequest(endpoint, method, jenkins_url, jenkins_username, jenkins_api_token, user_id, role_name, role_type); err != nil {
					return err
				}
			}
		}

		// Execute for Node Roles
		if node, ok := roleData["node"]; ok {
			nodeRoles := node.([]interface{})
			role_type := "slaveRoles"
			for _, nodeRole := range nodeRoles {
				role_name := nodeRole.(string)
				if err := sendRequest(endpoint, method, jenkins_url, jenkins_username, jenkins_api_token, user_id, role_name, role_type); err != nil {
					return err
				}
			}
		}
	}

	// Set Resource ID
	d.SetId(user_id)

	return nil
}

func resourceJenkinsRoleRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceJenkinsRoleUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceJenkinsRoleDelete(d *schema.ResourceData, m interface{}) error {
	jenkins_url, jenkins_username, jenkins_api_token, user_id, role_set := initVariables(d, m)

	roleList := role_set.List()
	endpoint := "/role-strategy/strategy/unassignRole"
	method := "POST"

	for _, role := range roleList {
		roleData := role.(map[string]interface{})

		// Execute for Global Roles
		if global, ok := roleData["global"]; ok {
			globalRoles := global.([]interface{})
			role_type := "globalRoles"
			for _, globalRole := range globalRoles {
				role_name := globalRole.(string)
				if err := sendRequest(endpoint, method, jenkins_url, jenkins_username, jenkins_api_token, user_id, role_name, role_type); err != nil {
					return err
				}
			}
		}

		// Execute for Item Roles
		if item, ok := roleData["item"]; ok {
			itemRoles := item.([]interface{})
			role_type := "projectRoles"
			for _, itemRole := range itemRoles {
				role_name := itemRole.(string)
				if err := sendRequest(endpoint, method, jenkins_url, jenkins_username, jenkins_api_token, user_id, role_name, role_type); err != nil {
					return err
				}
			}
		}

		// Execute for Node Roles
		if node, ok := roleData["node"]; ok {
			nodeRoles := node.([]interface{})
			role_type := "slaveRoles"
			for _, nodeRole := range nodeRoles {
				role_name := nodeRole.(string)
				if err := sendRequest(endpoint, method, jenkins_url, jenkins_username, jenkins_api_token, user_id, role_name, role_type); err != nil {
					return err
				}
			}
		}
	}

	// Set Resource ID
	d.SetId("")

	return nil
}

func convertToSet(data map[string][]string) *schema.Set {
	// Create *schema.Set
	set := schema.NewSet(schema.HashString, []interface{}{})

	// Convert
	for key, values := range data {
		for _, value := range values {
			item := map[string]interface{}{
				key: value,
			}
			set.Add(item)
		}
	}

	return set
}

func resourceJenkinsRoleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	user_id := d.Id()
	config := m.(MyProviderConfig)

	jenkins_url := config.jenkins_url
	jenkins_username := config.jenkins_username
	jenkins_api_token := config.jenkins_api_token

	roleTypes := []string{"globalRoles", "projectRoles", "slaveRoles"}
	roleData := map[string][]string{}
	var role_type string

	for _, roleType := range roleTypes {
		if roleType == "globalRoles" {
			role_type = "global"
		} else if roleType == "projectRoles" {
			role_type = "item"
		} else if roleType == "slaveRoles" {
			role_type = "node"
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("type", roleType)
		writer.Close()

		// Create request
		client := &http.Client{}
		req, err := http.NewRequest("GET", jenkins_url+"/role-strategy/strategy/getAllRoles", body)
		if err != nil {
			fmt.Println("Error creating request:", err)
		}
		req.SetBasicAuth(jenkins_username, jenkins_api_token)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Do HTTP request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
		}
		defer resp.Body.Close()

		// Read json output
		if resp.StatusCode == http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
			}

			output := string(body)

			var data map[string][]string
			err = json.Unmarshal([]byte(output), &data)
			if err != nil {
				fmt.Println("Error:", err)
			}

			for role, users := range data {
				for _, user := range users {
					if user == user_id {
						roleData[role_type] = append(roleData[role_type], role)
					}
				}
			}

		} else {
			fmt.Println("Error:", resp.Status)
		}
	}

	d.SetId(user_id)
	d.Set("jenkins_url", jenkins_url)
	d.Set("jenkins_username", jenkins_username)
	d.Set("jenkins_api_token", jenkins_api_token)
	d.Set("user_id", user_id)
	d.Set("role", []map[string]interface{}{
		{
			"global": roleData["global"],
			"item":   roleData["item"],
			"node":   roleData["node"],
		},
	})

	return []*schema.ResourceData{d}, nil
}