allow_k8s_contexts('aee')

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
    disable_push=True,
)

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
