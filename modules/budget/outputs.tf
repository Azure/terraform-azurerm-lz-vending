output "budget_resource_id" {
  description = "The Azure resource id of the created budget."
  value       = azurerm_consumption_budget.this.id
}
