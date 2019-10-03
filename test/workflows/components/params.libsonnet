{
  global: {
    // User-defined global parameters; accessible to all component and environments, Ex:
    // replicas: 4,
  },
  components: {
    builddocker: {
      bucket: "kubeflow-ci_temp",
      cluster: "kubeflow-testing",
      dockerfile: "Dockerfile",
      dockerfileDir: "kubeflow/xgboost-operator/",
      extra_args: "null",
      extra_repos: "kubeflow/testing@HEAD",
      gcpCredentialsSecretName: "gcp-credentials",
      image: "xgboost-operator",
      name: "xgboost-operator",
      namespace: "kubeflow-ci",
      nfsVolumeClaim: "nfs-external",
      project: "kubeflow-ci",
      prow_env: "REPO_OWNER=kubeflow,REPO_NAME=xgboost-operator,PULL_BASE_SHA=master",
      registry: "gcr.io/kubeflow-images-public",
      testing_image: "gcr.io/kubeflow-ci/test-worker/test-worker:v20190116-b7abb8d-e3b0c4",
      versionTag: "v1.0",
      zone: "us-central1-a",
    },
  },
}
