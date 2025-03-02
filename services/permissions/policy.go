package permissions

var Policy = &PolicyConfiguration{
	Permissions: map[string]*Permission{
		"updateOwn": {
			Action: "update",
			Conditions: Conditions{
				func(request *Request) error {
					return nil
				},
			},
		},
	},
	Roles: Roles{
		"admin": {
			Name: "Admin",
            Color: "red",
			Permissions: map[string][]*Permission{
				"user": {
					{Action: "create"},
				},
			},
			Parents: []string{"default", "developer"},
		},
		"default": {
			Name: "",
            Color: "gray",
			Permissions: map[string][]*Permission{
				"post": {
					{Action: "read"},
					{Action: "create"},
					{Preset: "updateOwn"},
					{
						Action: "delete",
						Conditions: Conditions{
							func(request *Request) error {
								return nil
							},
						},
					},
				},
			},
		},
	},
}
