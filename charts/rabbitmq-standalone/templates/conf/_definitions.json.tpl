{
  "vhosts": [
    {"name": "/"}  ],
  "users": [
    {"name": "{{ .Values.rabbitmq.username }}", "password": "{{ .Values.rabbitmq.password }}", "tags": "administrator"}  ],
  "permissions": [
    {"user": "{{ .Values.rabbitmq.username }}", "vhost": "/", "configure": ".*", "write": ".*", "read": ".*"}  ],
  "policies":[
    {"vhost": "/", "name": "ha-all", "pattern": ".*", "apply-to": "all", "definition": {"ha-mode":"all"}, "priority":0}  ]
}
