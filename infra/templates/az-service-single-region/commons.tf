module "provider" {
  source = "../../modules/providers/azure/provider"
}

resource "random_string" "workspace_scope" {
  keepers = {
    # Generate a new id each time we switch to a new workspace or app id
    ws_name = replace(trimspace(lower(terraform.workspace)), "_", "-")
    app_id  = replace(trimspace(lower(var.name)), "_", "-")
  }

  length  = max(1, var.randomization_level) // error for zero-length
  special = false
  upper   = false
}

locals {
  // sanitize names
  app_id  = random_string.workspace_scope.keepers.app_id
  region  = replace(trimspace(lower(var.resource_group_location)), "_", "-")
  ws_name = random_string.workspace_scope.keepers.ws_name
  suffix  = var.randomization_level > 0 ? "-${random_string.workspace_scope.result}" : ""

  // base name for resources, name constraints documented here: https://docs.microsoft.com/en-us/azure/architecture/best-practices/naming-conventions
  base_name    = length(local.app_id) > 0 ? "${local.ws_name}${local.suffix}-${local.app_id}" : "${local.ws_name}${local.suffix}"
  base_name_21 = length(local.base_name) < 22 ? local.base_name : "${substr(local.base_name, 0, 21 - length(local.suffix))}${local.suffix}"
  base_name_46 = length(local.base_name) < 47 ? local.base_name : "${substr(local.base_name, 0, 46 - length(local.suffix))}${local.suffix}"
  base_name_60 = length(local.base_name) < 61 ? local.base_name : "${substr(local.base_name, 0, 60 - length(local.suffix))}${local.suffix}"
  base_name_76 = length(local.base_name) < 77 ? local.base_name : "${substr(local.base_name, 0, 76 - length(local.suffix))}${local.suffix}"
  base_name_83 = length(local.base_name) < 84 ? local.base_name : "${substr(local.base_name, 0, 83 - length(local.suffix))}${local.suffix}"

  // Resource names
  admin_rg_name       = "${local.base_name_83}-adm-rg"               // resource group used for admin resources (max 90 chars)
  app_rg_name         = "${local.base_name_83}-app-rg"               // app resource group (max 90 chars)
  sp_name             = "${local.base_name}-sp"                      // service plan
  acr_name            = "${replace(local.base_name_46, "-", "")}acr" // container registry (max 50 chars, alphanumeric *only*)
  app_svc_name_prefix = local.base_name

  // Resolved TF Vars
  resolved_acr_rg_name = var.azure_container_resource_group == "" ? azurerm_resource_group.svcplan.name : var.azure_container_resource_group
  resolved_acr_name    = var.azure_container_resource_name == "" ? local.acr_name : var.azure_container_resource_name
}

