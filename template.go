package main

const (
	DefaultTemplate = `
{{range .}}
server {
        listen 80;
        server_name {{.Config.Hostname}};
        proxy_set_header Host {{.Config.Hostname}};


        location / {
                proxy_pass http://localhost:{{index .NetworkSettings.PortMapping "80"}};
        }

}
{{end}}
`
)
