allow_k8s_contexts('aee')


def helm_dep_up(path_to_helm_dir):
    update = False
    chart_yaml = read_yaml('%s/Chart.yaml' % path_to_helm_dir)
    for dep in chart_yaml['dependencies']:
        if not os.path.exists(path_to_helm_dir + '/charts/' + dep['name'] + '-' + dep['version'] + '.tgz'):
            update = True
            break

    if update:
        local('helm dep update %s' % path_to_helm_dir)


image = 'asia.gcr.io/atomicloud/local-dev'

custom_build(
    image,
    'docker buildx build --platform=linux/amd64 -t $EXPECTED_REF --push .',
    entrypoint='/tilt-restart-wrapper --watch_file="/.restart-process" /app/demo-app',
    deps=['.'],
    live_update=[
        sync('./go.mod', '/app/'),
        sync('./go.sum', '/app/'),
        sync('./main.go', '/app/'),
        sync('./store.go', '/app/'),
        run('cd /app && go mod download', trigger=['./go.mod', './go.sum']),
        run('cd /app && GOOS=linux GOARCH=amd64 go build -o demo-app .'),
        run('date > /.restart-process')
    ],
    skips_local_docker=True,
    disable_push=True
)

helm_dep_up('./helm')
k8s_yaml(helm(
    './helm',
    set=['image=%s' % image]
))
k8s_resource('tilt-demo', port_forwards=8080)

local_resource(
    'test',
    ['./test_endpoints.sh'],
    deps=['./test_endpoints.sh'],
    resource_deps=['tilt-demo']
)
