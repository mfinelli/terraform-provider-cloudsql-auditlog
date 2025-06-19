package provider

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"terraform-provider-cloudsql-auditlog/db"

	// "time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &auditLogRuleResource{}
	_ resource.ResourceWithConfigure = &auditLogRuleResource{}
	_ resource.ResourceWithImportState = &auditLogRuleResource{}
)

func NewAuditLogRuleResource() resource.Resource {
	return &auditLogRuleResource{}
}

type auditLogRuleResource struct{
	client *sql.DB
}

type auditLogRuleResourceModel struct {
	ID types.String `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
	DbName types.String `tfsdk:"dbname"`
	Object types.String `tfsdk:"object"`
	Operation types.String `tfsdk:"operation"`
	OpResult types.String `tfsdk:"op_result"`
	// LastUpdated types.String `tfsdk:"last_updated"`
}

func (r *auditLogRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_audit_log_rule"
}

func (r *auditLogRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Required: true,
			},
			"dbname": schema.StringAttribute{
				Required: true,
			},
			"object": schema.StringAttribute{
				Required: true,
			},
			"operation": schema.StringAttribute{
				Required: true,
			},
			"op_result": schema.StringAttribute{
				Required: true,
			},
			// "last_updated": schema.StringAttribute{
			// 	Computed: true,
			// },
		},
	}
}

func (r *auditLogRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// retrieve values from plan
	var plan auditLogRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	q := db.New(r.client)
	ruleIdCheck, err := q.ReadAuditRuleIDAfterCreate(ctx,
		db.ReadAuditRuleIDAfterCreateParams{
			Username: plan.Username.ValueString(),
			Dbname: plan.DbName.ValueString(),
			Object: plan.Object.ValueString(),
			Operation: plan.Operation.ValueString(),
			OpResult: plan.OpResult.ValueString(),
		})
	if err == nil {
		resp.Diagnostics.AddError(
			"Rule already exists",
			fmt.Errorf("existing ID: %d", ruleIdCheck).Error(),
		)
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		resp.Diagnostics.AddError(
			"Unable to check rule existence",
			err.Error(),
		)
		return
	}
	
	err = q.CreateAuditRule(ctx, db.CreateAuditRuleParams{
		Username: plan.Username.ValueString(),
		Dbname: plan.DbName.ValueString(),
		Object: plan.Object.ValueString(),
		Operation: plan.Operation.ValueString(),
		OpResult: plan.OpResult.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to call audit rule create",
			err.Error(),
		)
		return
	}

	ruleID, err := q.ReadAuditRuleIDAfterCreate(ctx,
		db.ReadAuditRuleIDAfterCreateParams{
			Username: plan.Username.ValueString(),
			Dbname: plan.DbName.ValueString(),
			Object: plan.Object.ValueString(),
			Operation: plan.Operation.ValueString(),
			OpResult: plan.OpResult.ValueString(),
		})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to call read after create",
			err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(strconv.FormatInt(ruleID, 10))
	// plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *auditLogRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state auditLogRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	q := db.New(r.client)
	ruleID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting id to int",
			fmt.Sprintf("Could not convert rule with id %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	rule, err := q.ReadAuditLogRuleByID(ctx, int64(ruleID))
	if err != nil  && !errors.Is(err, sql.ErrNoRows) {
		resp.Diagnostics.AddError(
			"Error reading audit log rule",
			fmt.Sprintf("Could not read rule with id %d: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	state.Username = types.StringValue(rule.Username)
	state.DbName = types.StringValue(rule.Dbname)
	state.Object = types.StringValue(rule.Object)
	state.OpResult = types.StringValue(rule.OpResult)
	state.Operation = types.StringValue(rule.Operation)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *auditLogRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// retrieve values from plan
	var plan auditLogRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	q := db.New(r.client)
	err := q.UpdatedAuditRuleByID(ctx, db.UpdatedAuditRuleByIDParams{
		ID: plan.ID.ValueString(),
		Username: plan.Username.ValueString(),
		Dbname: plan.DbName.ValueString(),
		Object: plan.Object.ValueString(),
		Operation: plan.Operation.ValueString(),
		OpResult: plan.OpResult.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to call audit rule update",
			err.Error(),
		)
		return
	}

	// plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *auditLogRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state auditLogRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	q := db.New(r.client)
	err := q.DeleteAuditRuleByID(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to call audit rule delete",
			err.Error(),
		)
		return
	}
}

func (r *auditLogRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sql.DB)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *sql.DB got %T.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *auditLogRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
