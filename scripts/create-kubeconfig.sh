#!/bin/bash -e

display_usage() {
  echo -e "\nCopies settings from a kubeconfig context on the current machine to a new"
  echo "kubeconfig file using a designated service account for credentials."
  echo -e "\nUsage:"
  echo "  $0 [options] <serviceaccount>"
  echo "  $0 -h | --help"
  echo -e "\nOptions:"
  echo -e "  -h --help  Show this help text"
  echo -e "  -o <file>, --output <file>  File to output to [default: ./kubeconfig]"
  echo -e "  -c <context>, --context <context>  Name of the context to copy [default: current context]"
  echo -e "  -n <namespace>, --namespace <namespace>  Kubernetes namespace [default: kube-system]"
  echo ""
}

remove_temp() {
  rm -rf ${TEMP_DIR}
}

while [[ $# -gt 0 ]]
do
  KEY="$1"
  case ${KEY} in
    -h|--help)
      display_usage
      exit 0
      ;;
    -c|--context)
      CONTEXT_NAME="$2"
      shift
      shift
      ;;
    -o|--output)
      KUBECONFIG="$2"
      shift
      shift
      ;;
    -n|--namespace)
      NAMESPACE="$2"
      shift
      shift
      ;;
    *)
      SA_NAME=${KEY}
      shift
      ;;
  esac
done

if [[ -z ${SA_NAME} ]]; then
  echo "Error: A service account must be specified"
  display_usage
  exit 1
fi

if [[ -z ${KUBECONFIG} ]]; then
  KUBECONFIG="./kubeconfig"
fi

if [[ -z ${NAMESPACE} ]]; then
  NAMESPACE="kube-system"
fi

if [[ -z ${CONTEXT_NAME} ]]; then
  CONTEXT_NAME=$(kubectl config current-context)
fi

TEMP_DIR=$(mktemp -d)

trap remove_temp EXIT

CLUSTER_NAME=$(kubectl config view -o jsonpath="{.contexts[?(@.name == '${CONTEXT_NAME}')].context.cluster}")
CLUSTER_URL=$(kubectl config view -o jsonpath="{.clusters[?(@.name == '${CLUSTER_NAME}')].cluster.server}")
SA_SECRET=$(kubectl get sa -n ${NAMESPACE} ${SA_NAME} -o jsonpath="{.secrets[0].name}")
SA_TOKEN=$(kubectl get secrets -n ${NAMESPACE} ${SA_SECRET} -o jsonpath="{.data.token}" | base64 -d)

kubectl get secrets -n ${NAMESPACE} ${SA_SECRET} -o jsonpath="{.data.ca\.crt}" | base64 -d > "${TEMP_DIR}/ca.crt"

kubectl config --kubeconfig=${KUBECONFIG} set-cluster ${CLUSTER_URL} \
  --server=${CLUSTER_URL} \
  --certificate-authority="${TEMP_DIR}/ca.crt" \
  --embed-certs=true

kubectl config --kubeconfig=${KUBECONFIG} set-credentials ${CLUSTER_NAME}_${SA_NAME} \
  --token=${SA_TOKEN}

kubectl config --kubeconfig=${KUBECONFIG} set-context ${CONTEXT_NAME} \
  --cluster=${CLUSTER_URL} \
  --user=${CLUSTER_NAME}_${SA_NAME}