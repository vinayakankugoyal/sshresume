#!/bin/bash
set -e

KEY_DIR="./temp_keys"
mkdir -p "$KEY_DIR"

echo "Generating Ed25519 host key..."
# Generate key pair, overwrite if exists, no passphrase
ssh-keygen -t ed25519 -f "$KEY_DIR/id_ed25519" -N "" -C "ssh-resume-host-key" -q

echo "Creating Kubernetes Secret manifest..."
# Create the secret yaml
kubectl create secret generic ssh-host-key \
  --from-file=id_ed25519="$KEY_DIR/id_ed25519" \
  --dry-run=client -o yaml > k8s-secret-ssh-host-key.yaml

echo "Secret manifest written to k8s-secret-ssh-host-key.yaml"
echo "Apply it to your cluster with: kubectl apply -f k8s-secret-ssh-host-key.yaml"
echo "Then apply the deployment: kubectl apply -f deployment.yaml"

# Cleanup private key from disk (it's in the yaml now)
rm -rf "$KEY_DIR"
