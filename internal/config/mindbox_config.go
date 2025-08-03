package config

type MindboxConfig struct {
	Url             string `env:"MINDBOX_URL" required:"true"`
	OperationPrefix string `env:"MINDBOX_OPERATION_PREFIX" required:"true"`
	Android         struct {
		Auth       string `env:"MINDBOX_ANDROID_AUTHORIZATION" required:"true"`
		EndpointID string `env:"MINDBOX_ANDROID_ENDPOINT_ID" required:"true"`
	}
	IOS struct {
		Auth       string `env:"MINDBOX_IOS_AUTHORIZATION" required:"true"`
		EndpointID string `env:"MINDBOX_IOS_ENDPOINT_ID" required:"true"`
	}
	Web struct {
		Auth       string `env:"MINDBOX_WEB_AUTHORIZATION" required:"true"`
		EndpointID string `env:"MINDBOX_WEB_ENDPOINT_ID" required:"true"`
	}
}
