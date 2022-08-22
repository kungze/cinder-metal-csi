RABBITMQ_LOG_BASE=/var/log/kolla/rabbitmq
RABBITMQ_DIST_PORT=25672
RABBITMQ_PID_FILE=/var/lib/rabbitmq/mnesia/rabbitmq.pid
RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS="-kernel inetrc '/etc/rabbitmq/erl_inetrc' "
RABBITMQ_CTL_ERL_ARGS=""

export ERL_EPMD_PORT=4369
export ERL_INETRC=/etc/rabbitmq/erl_inetrc
