// ---- General Configuration ----

variable "name" {
  description = "An identifier used to construct the names of all resources in this template."
  type        = string
}

variable "randomization_level" {
  description = "Number of additional random characters to include in resource names to insulate against unexpected resource name collisions."
  type        = number
  default     = 8
}

variable "resource_group_location" {
  description = "The Azure region where all resources in this template should be created."
  type        = string
}



// ---- App Service Configuration ----

variable "application_type" {
  description = "Type of the App Insights Application.  Valid values are ios for iOS, java for Java web, MobileCenter for App Center, Node.JS for Node.js, other for General, phone for Windows Phone, store for Windows Store and web for ASP.NET."
  default     = "Web"
}

variable "acr_build_git_source_url" {
  description = "The URL to a git repository (e.g., 'https://github.com/Azure-Samples/acr-build-helloworld-node.git') containing the docker build manifests."
  type        = string
}

variable "acr_build_docker_file" {
  description = "The relative path of the the docker file to the source code root folder. Default to 'Dockerfile'."
  type        = string
  default     = "Dockerfile"
}

variable "deployment_targets" {
  description = "Metadata about apps to deploy, such as image metadata."
  type = list(object({
    app_name                 = string
    image_name               = string
    image_release_tag_prefix = string
  }))
}

variable "azure_container_resource_group" {
  description = "The resource group name for the azure container registry instance."
  type        = string
  default     = ""
}

variable "azure_container_resource_name" {
  description = "The resource name for the azure container registry instance."
  type        = string
  default     = ""
}

variable "azure_container_tags" {
  description = "A mapping of tags to assign to the resource.."
  type        = map(string)
  default     = {}
}
