load('ext://restart_process', 'docker_build_with_restart')
allow_k8s_contexts('aee')

image = 'asia.gcr.io/atomicloud/local-dev'

local_resource(
    'compile',
    'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/demo-app .',
    deps=['./main.go', './store.go']
)
docker_build_with_restart(
    image,
    '.',
    entrypoint='/app/demo-app',
    only=['./build'],
    live_update=[sync('./build', '/app')]
)

k8s_yaml(helm(
    './helm',
    set=['image=%s' % image]
))
k8s_resource('skaffold-demo',
    port_forwards=8080,
    resource_deps=['compile']
)

