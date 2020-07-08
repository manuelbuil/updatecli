source:
  kind: helmChart
  spec:
    url: https://grafana.github.io/loki/charts
    name: loki

conditions:
  isKeyExistWithCorrectValue:
    kind: yaml
    spec:
      file: "updateCli.d/jenkins.yaml"
      key: "targets.imageTag.kind"
      value: "yaml"
  isKeyExistWithWrongValue:
    kind: yaml
    spec:
      file: "updateCli.d/jenkins.yaml"
      key: "targets.imageTag.kind"
      value: "wrongValueTest"
  isKeyExistDoesntExist:
    kind: yaml
    spec:
      file: "updateCli.d/jenkins.yaml"
      key: "targets.imageTag.kindkindkind"
      value: "doesntExistKeyCheck"
