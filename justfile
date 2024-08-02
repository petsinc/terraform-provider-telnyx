provider_dir := "provider/provider/internal/provider"

test:
    #!/usr/bin/env bash
    set -eou pipefail

    cd {{provider_dir}}

    export TF_ACC_TERRAFORM_PATH=$(asdf which terraform)
    TF_ACC=true go test