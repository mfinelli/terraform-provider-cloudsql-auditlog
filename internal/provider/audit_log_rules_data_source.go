package provider

import (
	"context"
	"database/sql"
	"fmt"
	"terraform-provider-cloudsql-auditlog/db"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
  _ datasource.DataSource = &auditLogRulesDataSource{}
  _ datasource.DataSourceWithConfigure = &auditLogRulesDataSource{}
)

// NewCoffeesDataSource is a helper function to simplify the provider implementation.
func NewAuditLogRulesDataSource() datasource.DataSource {
  return &auditLogRulesDataSource{}
}

// coffeesDataSource is the data source implementation.
type auditLogRulesDataSource struct{
	client *sql.DB
}

// coffeesDataSourceModel maps the data source schema data.
type auditLogRulesDataSourceModel struct {
	AuditLogRules []auditLogRulesModel `tfsdk:"audit_log_rules"`
}

// coffeesModel maps coffees schema data.
type auditLogRulesModel struct {
	ID types.Int64 `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
	DbName types.String `tfsdk:"dbname"`
	Object types.String `tfsdk:"object"`
	Operation types.String `tfsdk:"operation"`
	OpResult types.String `tfsdk:"op_result"`
}

// Metadata returns the data source type name.
func (d *auditLogRulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_audit_log_rules"
}

// Schema defines the schema for the data source.
func (d *auditLogRulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
  resp.Schema = schema.Schema{
	  Attributes: map[string]schema.Attribute {
		  "audit_log_rules": schema.ListNestedAttribute{
			  Computed: true,
			  NestedObject: schema.NestedAttributeObject{
				  Attributes: map[string]schema.Attribute {
					  "id": schema.Int64Attribute{
						  Computed: true,
					  },
					  "username": schema.StringAttribute{
						  Computed: true,
					  },
					  "dbname": schema.StringAttribute{
						  Computed: true,
					  },
					  "object": schema.StringAttribute{
						  Computed: true,
					  },
					  "operation": schema.StringAttribute{
						  Computed: true,
					  },
					  "op_result": schema.StringAttribute{
						  Computed: true,
					  },
				  },
			  },
		  },
	  },
  }
}

// Read refreshes the Terraform state with the latest data.
func (d *auditLogRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state auditLogRulesDataSourceModel

	q := db.New(d.client)
	rules, err := q.GetAllAuditRules(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to query audit rules",
			err.Error(),
		)
		return
	}

	for _, rule := range rules {
		ruleState := auditLogRulesModel {
			ID: types.Int64Value(rule.ID),
			Username: types.StringValue(rule.Username),
			DbName: types.StringValue(rule.Dbname),
			Object: types.StringValue(rule.Object),
			Operation: types.StringValue(rule.Operation),
			OpResult: types.StringValue(rule.OpResult),
		}

		state.AuditLogRules = append(state.AuditLogRules, ruleState)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *auditLogRulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sql.DB)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sql.DB got %T.", req.ProviderData),
		)

		return
	}

	d.client = client
}
