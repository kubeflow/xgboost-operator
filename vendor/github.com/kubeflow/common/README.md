# common

[![Build Status](https://travis-ci.com/kubeflow/common.svg?branch=master)](https://travis-ci.com/kubeflow/common/)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubeflow/common)](https://goreportcard.com/report/github.com/kubeflow/common)

Common APIs and libraries shared by other Kubeflow operator repositories.

This repo is currently under construction. The overall design can be found at https://github.com/kubeflow/tf-operator/issues/960.

This repo should not be dependent on any other Kubeflow repositories. Instead, other operators should import this one as a common
library.
