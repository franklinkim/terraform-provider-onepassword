package onepassword

import (
	"fmt"
	"strings"
	"testing"

	"github.com/1Password/connect-sdk-go/onepassword"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataSourceOnePasswordItemRead(t *testing.T) {
	meta := &testClient{}
	expectedItem := generateItem()

	DoGetItemFunc = func(uuid string, vaultUUID string) (*onepassword.Item, error) {
		return expectedItem, nil
	}

	dataSourceData := generateDataSource(t, expectedItem)

	err := dataSourceOnepasswordItemRead(dataSourceData, meta)
	if err != nil {
		t.Errorf("Unexpected error occured")
	}
	compareItemToSource(t, dataSourceData, expectedItem)
}

func compareItemToSource(t *testing.T, dataSourceData *schema.ResourceData, item *onepassword.Item) {
	if dataSourceData.Get("uuid") != item.ID {
		t.Errorf("Expected uuid to be %v got %v", item.ID, dataSourceData.Get("uuid"))
	}
	if dataSourceData.Get("vault") != item.Vault.ID {
		t.Errorf("Expected vault to be %v got %v", item.Vault.ID, dataSourceData.Get("vault"))
	}
	expectedCategory := strings.ToLower(fmt.Sprintf("%v", item.Category))
	if dataSourceData.Get("category") != expectedCategory {
		t.Errorf("Expected category to be %v got %v", expectedCategory, dataSourceData.Get("category"))
	}
	if dataSourceData.Get("title") != item.Title {
		t.Errorf("Expected title to be %v got %v", item.Title, dataSourceData.Get("title"))
	}
	if dataSourceData.Get("url") != item.URLs[0].URL {
		t.Errorf("Expected url to be %v got %v", item.URLs[0].URL, dataSourceData.Get("url"))
	}

	for _, f := range item.Fields {
		if dataSourceData.Get(f.Label) != f.Value {
			t.Errorf("Expected field %v to be %v got %v", f.Label, f.Value, dataSourceData.Get(f.Label))
		}
	}
}

func generateDataSource(t *testing.T, item *onepassword.Item) *schema.ResourceData {
	dataSourceData := schema.TestResourceDataRaw(t, dataSourceOnepasswordItem().Schema, nil)
	dataSourceData.Set("uuid", item.ID)
	dataSourceData.Set("vault", item.Vault.ID)
	dataSourceData.SetId(fmt.Sprintf("vaults/%s/items/%s", item.Vault.ID, item.ID))
	return dataSourceData
}

func generateItem() *onepassword.Item {
	item := onepassword.Item{}
	item.Fields = generateFields()
	item.ID = "79841a98-dd4a-4c34-8be5-32dca20a7328"
	item.Vault.ID = "df2e9643-45ad-4ff9-8b98-996f801afa75"
	item.Category = "USERNAME"
	item.Title = "test item"
	item.URLs = []onepassword.ItemURL{
		{
			Primary: true,
			URL:     "some_url.com",
		},
	}
	return &item
}

func generateFields() []*onepassword.ItemField {
	fields := []*onepassword.ItemField{
		{
			Label: "username",
			Value: "test_user",
		},
		{
			Label: "password",
			Value: "test_password",
		},
		{
			Label: "hostname",
			Value: "test_host",
		},
		{
			Label: "database",
			Value: "test_database",
		},
		{
			Label: "port",
			Value: "test_port",
		},
		{
			Label: "type",
			Value: "test_type",
		},
	}
	return fields
}
