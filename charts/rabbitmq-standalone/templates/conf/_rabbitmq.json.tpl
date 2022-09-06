{
    "command": "/usr/sbin/rabbitmq-server",
    "config_files": [
        {
            "source": "/var/lib/kolla/config_files/rabbitmq-env.conf",
            "dest": "/etc/rabbitmq/rabbitmq-env.conf",
            "owner": "rabbitmq",
            "perm": "0600"
        },
        {
            "source": "/var/lib/kolla/config_files/rabbitmq.conf",
            "dest": "/etc/rabbitmq/rabbitmq.conf",
            "owner": "rabbitmq",
            "perm": "0600"
        },  
        {
            "source": "/var/lib/kolla/config_files/erl_inetrc",
            "dest": "/etc/rabbitmq/erl_inetrc",
            "owner": "rabbitmq",
            "perm": "0600"
        },
        {
            "source": "/var/lib/kolla/config_files/definitions.json",
            "dest": "/etc/rabbitmq/definitions.json",
            "owner": "rabbitmq",
            "perm": "0600"
        }
    ],
    "permissions": [
        {
            "path": "/var/lib/rabbitmq",
            "owner": "rabbitmq:rabbitmq",
            "recurse": true
        },
        {
            "path": "/var/log/kolla/rabbitmq",
            "owner": "rabbitmq:rabbitmq",
            "recurse": true
        }
    ]
}
