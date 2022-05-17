package connect

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Mirantis/terraform-provider-mirantis/mirantis/mke/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MKE_ENDPOINT", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MKE_USER", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("MKE_PASS", nil),
			},
			"unsafe_ssl_client": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				DefaultFunc: schema.EnvDefaultFunc("MKE_UNSAFE_CLIENT", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mirantis-mke-connect_clientbundle": ResourceClientBundle(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	endpoint := d.Get("endpoint").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	unsafeClient := d.Get("unsafe_ssl_client").(bool)

	if (username == "") || (password == "") || (endpoint == "") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create MKE client",
			Detail:   "Unable to create anonymous MKE client",
		})

		return nil, diags
	}

	var c client.Client
	var clientErr error

	if unsafeClient {
		c, clientErr = client.NewUnsafeSSLClient(endpoint, username, password)
	} else {
		c, clientErr = client.NewClientSimple(endpoint, username, password)
	}

	if clientErr != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create MKE client",
			Detail:   "Unable to authenticate user for authenticated MKE client",
		})

		return nil, diags
	}

	if err := c.ApiPing(ctx); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "MKE endpoint is not healthy",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return c, diags
}
