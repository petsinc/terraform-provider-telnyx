provider_dir := "provider/provider/internal/provider"

test:
    #!/usr/bin/env bash
    set -eou pipefail
    cd {{provider_dir}}

    export TF_ACC_TERRAFORM_PATH=$(asdf which terraform)
    TF_ACC=true go test 2>&1 | tee -a {{justfile_directory()}}/last-test.log ; ( exit ${PIPESTATUS} )


format-hcl *FILES='.':
    #!/usr/bin/env bash
    # If files is ., then just run at root recursively
    if [ "{{FILES}}" = "." ]; then
        terraform fmt -recursive
        exit 0;
    fi

    # Specific files, run them one by one
    for file in {{FILES}}; do
        if [[ "$file" == *.tf ]]; then
            terraform fmt "$file"
        fi
    done

format-go *FILES=provider_dir:
    gofmt -w {{FILES}}

