package main

type ACLUser struct {
	Username  string
	Passwords []string
}

type ACLManager struct {
	users map[string]*ACLUser
}

func NewAclManager() *ACLManager {
	acl := &ACLManager{
		users: make(map[string]*ACLUser),
	}

	acl.users["default"] = &ACLUser{
		Username:  "default",
		Passwords: []string{},
	}

	return acl
}

func (aclManager ACLManager) Authenticate(username, password string) *ACLUser {
	if username == "" && password == "" {
		return aclManager.users["default"]
	}

	panic("not implemented")
}
