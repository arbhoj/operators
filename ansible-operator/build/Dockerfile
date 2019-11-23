FROM quay.io/operator-framework/ansible-operator:v0.11.0

COPY watches.yaml ${HOME}/watches.yaml

COPY roles/ ${HOME}/roles/

COPY takebackups.yaml ${HOME}/takebackups.yaml

COPY deletebackups.yaml ${HOME}/deletebackups.yaml
