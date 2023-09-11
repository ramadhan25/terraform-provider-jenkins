# jenkins_role Resource

Manages a role assignment within Jenkins.

~> The Jenkins installation that uses this resource is expected to have the [Role-based Authorization Strategy Plugin](https://plugins.jenkins.io/role-strategy) installed in their system.

## Example Usage

```hcl
resource "jenkins_role" "example" {
  user_id = "user_id@example.com"
  role {
    global = ["global_role1","global_role2"]
    item = ["item_role1","item_role2"]
    node = ["node_role1","node_role2"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required) The user id to assign to the role.
* `role` - (Optional) Define the roles that will be assigned, there are 3 roles that can be filled:
  
  - global
  - item
  - node

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Same as user id.
* `user_id` - The user id to assign to the role.
* `role` - Role that has been assigned to the user

## Import

Role assignment may be imported by their user id, e.g.

```sh
$ terraform import jenkins_role.example user_id@example.com
```
