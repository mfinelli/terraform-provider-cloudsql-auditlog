// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &ScaffoldingProvider{}
var _ provider.ProviderWithFunctions = &ScaffoldingProvider{}
var _ provider.ProviderWithEphemeralResources = &ScaffoldingProvider{}

// ScaffoldingProvider defines the provider implementation.
type ScaffoldingProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type cloudsqlAuditlogProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Engine   types.String `tfsdk:"engine"`
}

type CloudSqlClientAndConfig struct {
	client *sql.DB
	engine string
}

func (p *ScaffoldingProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cloudsql-auditlog"
	resp.Version = p.version
}

func (p *ScaffoldingProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"username": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"password": schema.StringAttribute{
				Required:  false,
				Optional:  true, // empty password allowed for e.g., cloud-sql-proxy
				Sensitive: true,
			},
			"engine": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
		},
	}
}

func (p *ScaffoldingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data cloudsqlAuditlogProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.

	if data.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown endpoint",
			"Must set mysql endpoint",
		)
	}

	if data.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown username",
			"Must set mysql username",
		)
	}

	if data.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown password",
			"Must set mysql password",
		)
	}

	if data.Engine.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("engine"),
			"Unknown engine",
			"Must set engine type",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := ""
	username := ""
	password := ""

	if !data.Endpoint.IsNull() {
		endpoint = data.Endpoint.ValueString()
	}

	if !data.Username.IsNull() {
		username = data.Username.ValueString()
	}

	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing mysql endpoint",
			"Must set mysql endpoint",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing mysql username",
			"Must set mysql username",
		)
	}

	if data.Engine.ValueString() != "mysql" && data.Engine.ValueString() != "postgresql" {
		resp.Diagnostics.AddAttributeError(
			path.Root("engine"),
			"Invalid engine",
			fmt.Sprintf("Invalid engine type %q, allowed values: mysql, postgresql", data.Engine.ValueString()),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Engine.ValueString() == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/mysql", username, password, endpoint)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("unable to open mysql: %v", err),
				"Unable to call sql.open",
			)
			return
		}

		clientEngine := CloudSqlClientAndConfig{
			client: db,
			engine: data.Engine.ValueString(),
		}

		resp.DataSourceData = clientEngine
		resp.ResourceData = clientEngine
	} else {
		resp.Diagnostics.AddError(
			"TODO",
			"postgresql not implemented yet",
		)
		return
	}
}

func (p *ScaffoldingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAuditLogRuleResource,
	}
}

func (p *ScaffoldingProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *ScaffoldingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAuditLogRulesDataSource,
	}
}

func (p *ScaffoldingProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ScaffoldingProvider{
			version: version,
		}
	}
}
