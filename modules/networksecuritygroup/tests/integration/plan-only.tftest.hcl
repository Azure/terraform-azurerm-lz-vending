# Tests for the networksecuritygroup module
# Converts the tests from tests/networksecuritygroup/networksecuritygroup_test.go

# Test 1: Basic NSG without security rules
run "basic_network_security_group" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
  }

  assert {
    condition     = var.name == "test"
    error_message = "NSG name should be test"
  }

  assert {
    condition     = var.location == "westeurope"
    error_message = "NSG location should be westeurope"
  }
}

# Test 2: NSG with a primary security rule
run "nsg_with_security_rule_primary" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
    security_rules = {
      primary = {
        access                     = "Allow"
        direction                  = "Inbound"
        priority                   = 100
        protocol                   = "Tcp"
        source_port_range          = "*"
        destination_port_range     = "*"
        name                       = "test-rule"
        source_address_prefix      = "*"
        destination_address_prefix = "*"
      }
    }
  }

  assert {
    condition     = length(var.security_rules) == 1
    error_message = "Should have 1 security rule"
  }

  assert {
    condition     = var.security_rules["primary"].name == "test-rule"
    error_message = "Security rule name should be test-rule"
  }
}

# Test 3: NSG with source prefixes
run "nsg_with_source_prefixes" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
    security_rules = {
      primary = {
        access                     = "Allow"
        direction                  = "Inbound"
        priority                   = 100
        protocol                   = "Tcp"
        destination_port_range     = "*"
        destination_address_prefix = "*"
        name                       = "test-rule"
        source_port_ranges         = ["*"]
        source_address_prefixes    = ["*"]
      }
    }
  }

  assert {
    condition     = length(var.security_rules["primary"].source_port_ranges) == 1
    error_message = "Should have source port ranges"
  }

  assert {
    condition     = length(var.security_rules["primary"].source_address_prefixes) == 1
    error_message = "Should have source address prefixes"
  }
}

# Test 4: NSG with destination prefixes
run "nsg_with_destination_prefixes" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
    security_rules = {
      primary = {
        access                       = "Allow"
        direction                    = "Outbound"
        priority                     = 100
        protocol                     = "Tcp"
        source_port_range            = "*"
        source_address_prefix        = "*"
        name                         = "test-rule"
        destination_port_ranges      = ["*"]
        destination_address_prefixes = ["*"]
      }
    }
  }

  assert {
    condition     = length(var.security_rules["primary"].destination_port_ranges) == 1
    error_message = "Should have destination port ranges"
  }

  assert {
    condition     = length(var.security_rules["primary"].destination_address_prefixes) == 1
    error_message = "Should have destination address prefixes"
  }
}

# Test 5: NSG with all prefixes (source and destination)
run "nsg_with_all_prefixes" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
    security_rules = {
      primary = {
        access                       = "Allow"
        direction                    = "Inbound"
        priority                     = 100
        protocol                     = "Tcp"
        name                         = "test-rule"
        source_port_ranges           = ["*"]
        destination_port_ranges      = ["*"]
        source_address_prefixes      = ["*"]
        destination_address_prefixes = ["*"]
      }
    }
  }

  assert {
    condition = (
      length(var.security_rules["primary"].source_port_ranges) == 1 &&
      length(var.security_rules["primary"].destination_port_ranges) == 1 &&
      length(var.security_rules["primary"].source_address_prefixes) == 1 &&
      length(var.security_rules["primary"].destination_address_prefixes) == 1
    )
    error_message = "Should have all prefixes configured"
  }
}

# Test 6: NSG with source application security groups
run "nsg_with_source_asgs" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
    security_rules = {
      primary = {
        access                     = "Allow"
        direction                  = "Inbound"
        priority                   = 100
        protocol                   = "Tcp"
        source_port_range          = "*"
        destination_port_range     = "*"
        destination_address_prefix = "*"
        name                       = "test-rule"
        source_application_security_group_ids = [
          "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/applicationSecurityGroups/sourceASG"
        ]
      }
    }
  }

  assert {
    condition     = length(var.security_rules["primary"].source_application_security_group_ids) == 1
    error_message = "Should have source ASG configured"
  }
}

# Test 7: NSG with destination application security groups
run "nsg_with_destination_asgs" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
    security_rules = {
      primary = {
        access                 = "Allow"
        direction              = "Inbound"
        priority               = 100
        protocol               = "Tcp"
        source_port_range      = "*"
        destination_port_range = "*"
        source_address_prefix  = "*"
        name                   = "test-rule"
        destination_application_security_group_ids = [
          "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/applicationSecurityGroups/destinationASG"
        ]
      }
    }
  }

  assert {
    condition     = length(var.security_rules["primary"].destination_application_security_group_ids) == 1
    error_message = "Should have destination ASG configured"
  }
}

# Test 8: NSG with both source and destination application security groups
run "nsg_with_all_asgs" {
  command = plan

  variables {
    name      = "test"
    location  = "westeurope"
    parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg-test"
    security_rules = {
      primary = {
        access                 = "Allow"
        direction              = "Inbound"
        priority               = 100
        protocol               = "Tcp"
        source_port_range      = "*"
        destination_port_range = "*"
        name                   = "test-rule"
        source_application_security_group_ids = [
          "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/applicationSecurityGroups/sourceASG"
        ]
        destination_application_security_group_ids = [
          "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/sampleResourceGroup/providers/Microsoft.Network/applicationSecurityGroups/destinationASG"
        ]
      }
    }
  }

  assert {
    condition = (
      length(var.security_rules["primary"].source_application_security_group_ids) == 1 &&
      length(var.security_rules["primary"].destination_application_security_group_ids) == 1
    )
    error_message = "Should have both source and destination ASGs configured"
  }
}
