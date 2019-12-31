# mqp setup

### Install rabbitmq

    brew install rabbitmq

    brew services start rabbitmq

### create user and vhost

    rabbitmqctl add_user mqp mqptest

    rabbitmqctl add_vhost mqp

    rabbitmqctl set_permissions --vhost mqp mqp ".*" ".*" ".*"
