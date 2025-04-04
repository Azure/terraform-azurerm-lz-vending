mock_provider "azapi" {
  override_data {
    target = data.azapi_resource_list.role_definitions
    values = {
      output = {
        results = [
          {
            "id" : "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635",
            "role_name" : "Owner"
          },
          {
            "id" : "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7",
            "role_name" : "Reader"
          },
        ]
      }
    }
  }
}

variables {
  role_assignment_principal_id = "00000000-0000-0000-0000-000000000000"
  role_assignment_scope        = "/subscriptions/00000000-0000-0000-0000-000000000000"
  role_assignment_definition   = "Owner"
}

run "simple_role_name_valid" {
  command = plan
  assert {
    error_message = "Definition id is not correct"
    condition     = azapi_resource.this.body.properties.roleDefinitionId == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635"
  }
}

run "simple_role_definition_id_valid" {
  command = plan

  variables {
    role_assignment_definition = "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635"
  }
  assert {
    error_message = "Definition id is not correct"
    condition     = azapi_resource.this.body.properties.roleDefinitionId == "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635"
  }
}

run "scope_invalid" {
  command = plan

  variables {
    role_assignment_scope = "/"
  }

  expect_failures = [var.role_assignment_scope]
}

run "condition_valid_v2" {
  command = plan

  variables {
    role_assignment_condition         = "(!(ActionMatches{'Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read'} AND NOT SubOperationMatches{'Blob.List'}))"
    role_assignment_condition_version = "2.0"
  }
}


run "condition_valid_v1" {
  command = plan

  variables {
    role_assignment_condition         = "(!(ActionMatches{'Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read'} AND NOT SubOperationMatches{'Blob.List'}))"
    role_assignment_condition_version = "1.0"
  }
}

run "condition_invalid" {
  command = plan

  variables {
    role_assignment_condition         = "(!(ActionMatches{'Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read'} AND NOT SubOperationMatches{'Blob.List'}))"
    role_assignment_condition_version = "2.2"
  }

  expect_failures = [var.role_assignment_condition_version]
}
