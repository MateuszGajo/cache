package acl

type ACLUser struct {
	Username  string
	passwords []string
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
		passwords: []string{},
	}

	return acl
}

func (aclManager ACLManager) Authenticate(username, password string) *ACLUser {
	if username == "" && password == "" {
		return aclManager.users["default"]
	}

	panic("not implemented")
}

type UserRuleFlag string

const (
	flagAllKeys UserRuleFlag = "allkeys"
	flagNoPass  UserRuleFlag = "nopass"
)

type ACLUserRules struct {
	Flags []UserRuleFlag
}

func (aclManager ACLManager) findUser(username string) *ACLUser {
	for _, user := range aclManager.users {
		if username == user.Username {
			return user
		}
	}

	return nil
}

func (aclManager ACLManager) GetRules(username string) *ACLUserRules {
	user := aclManager.findUser(username)

	if user == nil {
		return nil
	}

	flags := []UserRuleFlag{}

	if len(user.passwords) == 0 {
		flags = append(flags, flagNoPass)
	}

	return &ACLUserRules{
		Flags: flags,
	}
}
