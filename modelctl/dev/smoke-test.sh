#!/bin/bash

# Instrucions:
#
# This script runs smokes tests for a model snap installed from store.
# It might not work for a local installation unless that all needed components are installed.

set -Eeuo pipefail

# =============================================================================
# ERROR HANDLING
# =============================================================================

error_handler() {
  local exit_code=$?
  local line_no=$1
  local bash_lineno=$2
  local last_command="$3"
  local func_name="${4:-main}"

  log_error "Script failed with exit code $exit_code"
  log_error "Error occurred in function: $func_name"
  log_error "Failed command: $last_command"
  log_error "Line number: $line_no"
  log_error "Bash line number: $bash_lineno"

  # Print call stack
  log_error "Call stack:"
  local frame=0
  while caller $frame >/dev/null 2>&1; do
    local caller_info
    caller_info=$(caller $frame)
    log_error "  [$frame] $caller_info"
    ((frame++))
  done

  exit "$exit_code"
}

# Set up trap for ERR signal
# ${LINENO} - line number where error occurred
# ${BASH_LINENO[0]} - line number in the calling function
# ${BASH_COMMAND} - command that caused the error
# ${FUNCNAME[1]} - name of the function where error occurred
trap 'error_handler ${LINENO} ${BASH_LINENO[0]} "$BASH_COMMAND" "${FUNCNAME[1]}"' ERR

# =============================================================================
# CONFIGURATION AND GLOBALS
# =============================================================================

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Features that can be validated during smoke tests
# To add:
# - openai_chat_image_recognition
# - openai_chat_image_ocr
# - openai_text_embeddings
# - openai_transcription
# - openai_realtime_transcription
# - openai_realtime_transcription_logprops
# - systemone_image_decision
SUPPORTED_FEATURES=(
  openai_models
  openai_chat_text
  systemone_decision
)

# =============================================================================
# UTILITY FUNCTIONS
# =============================================================================

log_info() {
  echo -e "${GREEN}[INFO]${NC} $1"
}

log_warning() {
  echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
  echo -e "${RED}[ERROR]${NC} $1" >&2
}

log_section() {
  echo -e "\n${BLUE}=== $1 ===${NC}"
}

log_debugging_info() {
  log_section "Snap logs"
  snap logs -n all "$SNAP_NAME" || log_warning "Could not retrieve snap logs."

  log_section "Machine info"
  "$SNAP_NAME" machine || log_warning "Could not retrieve machine info."
}

exit_error() {
  log_error "$1"

  if [[ -n "${GITHUB_ACTIONS:-}" && -n "${SNAP_NAME:-}" ]]; then
    echo "::group:: Debugging Information"
    log_debugging_info
    echo "::endgroup::"
  fi

  exit 1
}

usage() {
  local supported_features="${SUPPORTED_FEATURES[*]}"
  echo "Usage: $0 <inference-snap-name> <engine> [--features=x,y,z]"
  echo "Runs smoke tests for a specific engine against an inference snap."
  echo
  echo "--features: Comma separated list of features to validate."
  echo "            Optional; when omitted, no feature tests are run."
  echo "            Supported values: ${supported_features// /, }"
  echo
  echo "Example:"
  echo "./$(basename "$0") gemma4 cpu"
}

# =============================================================================
# VALIDATION FUNCTIONS
# =============================================================================

check_root_privileges() {
  if ((EUID != 0)); then
    exit_error "This script must be run as root or with sudo."
  fi
}

check_for_curl() {
  if ! command -v curl &>/dev/null; then
    exit_error "curl is not installed. Please install curl to continue."
  fi
}

check_for_jq() {
  if ! command -v jq &>/dev/null; then
    exit_error "jq is not available. Please install it and try again."
  fi
}

validate_arguments() {
  if [[ -z "$SNAP_NAME" || -z "$ENGINE" ]]; then
    usage
    exit_error "Engine and snap name are required."
  fi
}

validate_features() {
  local features="$1"
  local -a feature_list
  IFS=', ' read -r -a feature_list <<<"$features"

  for feature in "${feature_list[@]}"; do
    [[ -z "$feature" ]] && continue
    local supported
    local found=false
    for supported in "${SUPPORTED_FEATURES[@]}"; do
      if [[ "$feature" == "$supported" ]]; then
        found=true
        break
      fi
    done
    if [[ "$found" != true ]]; then
      usage
      exit_error "Unknown feature: '$feature'. Supported values: ${SUPPORTED_FEATURES[*]}."
    fi
  done
}

check_port_listening() {
  local port="$1"
  local timeout_seconds=300
  local retry_delay=10
  local start_time
  start_time=$(date +%s)

  while true; do
    local current_time
    current_time=$(date +%s)
    local elapsed=$((current_time - start_time))

    log_info "Checking whether port $port is listening (${elapsed}/${timeout_seconds}s)"

    if ss -tuln | grep -q ":$port "; then
      log_info "✓ Port $port is listening"
      return 0
    fi

    if [[ $elapsed -lt $timeout_seconds ]]; then
      log_warning "Port $port not listening; retrying in ${retry_delay}s"
      sleep "$retry_delay"
    else
      exit_error "✗ Time out after ${timeout_seconds}s waiting for port $port to listen."
    fi
  done
}

# =============================================================================
# HTTP API TESTING FUNCTIONS
# =============================================================================

test_openai_models() {
  local timeout_seconds=300  # 5 minutes
  local retry_delay=10
  local connection_timeout=60
  local start_time
  start_time=$(date +%s)

  local base_url=$("$SNAP_NAME" status --format=json | jq -r '.entrypoints.openai.url // empty')
  if [[ -z "$base_url" ]]; then
    exit_error "Could not determine OpenAI base URL from status output."
  fi
  local endpoint="$base_url/models"

  log_info "Testing OpenAI models endpoint."

  while true; do
    local current_time
    current_time=$(date +%s)
    local elapsed=$((current_time - start_time))

    log_info "Checking $endpoint ($elapsed/${timeout_seconds}s)"

    if curl --retry 0 --fail-with-body --write-out '\n' --connect-timeout $connection_timeout "$endpoint"; then
      log_info "✓ $endpoint: Pass"
      return 0
    fi

    # TODO: warn or error if the model name doesn't match the name in status

    current_time=$(date +%s)
    elapsed=$((current_time - start_time))

    if [[ $elapsed -lt $timeout_seconds ]]; then
      log_warning "Endpoint failed; retrying in ${retry_delay}s"
      sleep "$retry_delay"
    else
      exit_error "✗ $endpoint: Failed after $timeout_seconds seconds"
    fi
  done
}

test_openai_chat_text() {
  local max_retries=5
  local retry_delay=60
  local connection_timeout=60
  local attempt=1

  local base_url=$("$SNAP_NAME" status --format=json | jq -r '.entrypoints.openai.url // empty')
  if [[ -z "$base_url" ]]; then
    exit_error "Could not determine OpenAI base URL from status output."
  fi
  local model_name=$("$SNAP_NAME" status --format=json | jq -r '.model.name')
  local endpoint="$base_url/chat/completions"

  log_info "Testing OpenAI chat completions endpoints."

  local system_message="You are a helpful assistant."
  local prompt="Hello!"
  local json_body
  json_body=$(
    cat <<EOF
{
  "model": "$model_name",
  "messages": [
    {
      "role": "user",
      "content": "$prompt"
    }
  ],
  "temperature": 0,
  "max_tokens": 5
}
EOF
  )

  local compact_json_body
  compact_json_body=$(echo "$json_body" | jq -c .)

  local api_response
  while [[ $attempt -le $max_retries ]]; do
    log_info "Checking $endpoint ($attempt/$max_retries)"

    local curl_exit_code=0
    set -x # log the curl command for debugging
    api_response=$(
      curl -X POST "$endpoint" \
        -H "Content-Type: application/json" \
        -d "$compact_json_body" \
        --connect-timeout $connection_timeout \
        --max-time 600 \
        --retry 0 \
        --fail-with-body \
        --write-out '\n' \
        2>/dev/null
    ) || curl_exit_code=$? # '||' avoids triggering the ERR trap on expected failures
    set +x

    if [[ $curl_exit_code -eq 0 ]]; then
      if [[ -z "$api_response" ]]; then
        exit_error "Empty response from server"
      fi

      log_info "✓ $endpoint: Pass"
      return 0
    fi

    if [[ $attempt -lt $max_retries ]]; then
      log_warning "Chat completion failed; retrying in ${retry_delay}s"
      sleep "$retry_delay"
    fi

    ((attempt++))
  done

  exit_error "✗ $endpoint: Failed after $max_retries attempts"

}

validate_systemone_route_choice_response() {
  local response="$1"

  # check if it is a valid json object
  if ! jq -e 'type == "object"' <<<"$response" >/dev/null 2>&1; then
    log_error "Response is not a valid JSON object"
    return 1
  fi

  # check if response contains answers
  if ! jq -e '.answers | type == "object"' <<<"$response" >/dev/null 2>&1; then
    log_error "Response is missing required fields: answers"
    return 1
  fi

  # check if it contains an answer about the "route" question
  if ! jq -e '.answers.route' <<<"$response" >/dev/null 2>&1; then
    log_error "Decision about \"route\" is missing from the decision response; required field is missing: answers.route"
    return 1
  fi

  # check if the "route" question is of the correct type
  if ! jq -e '.answers.route.type == "choice"' <<<"$response" >/dev/null 2>&1; then
    log_error "Expected answers.route.type to be \"choice\""
    return 1
  fi

  # check if the "route" question was answered correctly
  if ! jq -e '.answers.route.choice == "billing"' <<<"$response" >/dev/null 2>&1; then
    log_error "Model provided the wrong answer: answers.route.choice must be \"billing\""
    return 1
  fi

  # Validation passed
}

test_systemone_decision() {
  local max_retries=5
  local retry_delay=60
  local connection_timeout=60
  local attempt=1

  local model_status
  model_status=$("$SNAP_NAME" status --format=json)
  local url
  local model_name
  url=$(echo "$model_status" | jq -r '.entrypoints.systemone.url // empty')
  if [[ -z "$url" ]]; then
    exit_error "Could not determine SystemOne URL from status output."
  fi
  model_name=$(echo "$model_status" | jq -r '.model.name')

  log_info "Testing SystemOne decision endpoint."

  local request_body
  request_body=$(
    cat <<EOF
{
  "model": "$model_name",
  "state": "Customer message: I was charged twice for my monthly subscription this morning. I need this fixed immediately.",
  "questions": {
    "route": {
      "type": "choice",
      "instructions": "Which team should handle this?",
      "criteria": {
        "billing": "payments, charges, refunds, invoices",
        "shipping": "delivery, tracking, lost or late parcels",
        "technical": "bugs, errors, login problems"
      }
    }
  }
}
EOF
)
  local api_response

  while [[ $attempt -le $max_retries ]]; do
    log_info "Checking $url ($attempt/$max_retries)"

    local curl_exit_code=0
    set -x # log the curl command for debugging
    api_response=$(
      curl -X POST "$url" \
        -H "Content-Type: application/json" \
        -d "$request_body" \
        --connect-timeout $connection_timeout \
        --max-time 600 \
        --retry 0 \
        --fail-with-body \
        --write-out '\n' \
        2>/dev/null
    ) || curl_exit_code=$? # '||' avoids triggering the ERR trap on expected failures
    set +x

    if [[ $curl_exit_code -eq 0 ]]; then
      if [[ -z "$api_response" ]]; then
        exit_error "Empty response from server"
      fi

      if ! validate_systemone_route_choice_response "$api_response"; then
        log_error "Response: $api_response"
        exit_error "Invalid SystemOne response."
      fi

      log_info "✓ $url: Pass"
      return 0
    fi

    if [[ $attempt -lt $max_retries ]]; then
      log_warning "Decision endpoint failed; retrying in ${retry_delay}s"
      sleep "$retry_delay"
    fi

    ((attempt++))
  done

  exit_error "✗ $url: Failed after $max_retries attempts"
}

test_features() {
  local features="$1"

  log_section "Feature Tests"

  # Split comma-separated list into an array.
  local -a feature_list
  IFS=', ' read -r -a feature_list <<<"$features"

  if [[ ${#feature_list[@]} -eq 0 ]]; then
    log_warning "No features specified for testing."
    return 0
  fi

  for feature in "${feature_list[@]}"; do
    case "$feature" in
    openai_models)
      test_openai_models
      ;;
    openai_chat_text)
      test_openai_chat_text
      ;;
    systemone_decision)
      test_systemone_decision
      ;;
    "")
      # Ignore empty entries from consecutive separators.
      ;;
    *)
      exit_error "Unknown feature: '$feature'. Supported values: ${SUPPORTED_FEATURES[*]}."
      ;;
    esac
  done
}

# =============================================================================
# SNAP MANAGEMENT FUNCTIONS
# =============================================================================

test_snap_installation() {
  log_section "Snap Installation Test"
  log_info "Checking snap installation..."
  snap list "$SNAP_NAME"
}

test_configuration_management() {
  local default_port

  log_section "Configuration Tests"

  log_info "Print internal configs (snap get $SNAP_NAME)..."
  snap get "$SNAP_NAME" -d

  log_info "Print configs ($SNAP_NAME get)..."
  "$SNAP_NAME" get

  log_info "Getting specific config..."
  default_port=$("$SNAP_NAME" get http.port)
  echo "$default_port"

  log_info "Testing config change..."
  "$SNAP_NAME" set http.port=9999 --assume-yes

  # Verify config change persisted
  local port
  port=$("$SNAP_NAME" get http.port)
  if (("$port" != 9999)); then
    exit_error "Config change did not persist."
  fi
  log_info "✓ Config change persisted successfully"

  log_info "Reverting config change..."
  "$SNAP_NAME" set http.port="$default_port" --assume-yes
}

# =============================================================================
# ENGINE MANAGEMENT FUNCTIONS
# =============================================================================

test_engine_listing() {
  log_section "Engine Listing Tests"

  log_info "Comparing available vs declared engines..."

  mapfile -t avail_engines < <("$SNAP_NAME" engines --format=json | jq -r '.engines[].name' | sort)
  echo -e "Available engines:\n${avail_engines[*]}"

  mapfile -t src_engines < <(find "/snap/$SNAP_NAME/current/engines/" -maxdepth 1 -mindepth 1 -type d -printf '%f\n' | sort)
  echo -e "Declared engines:\n${src_engines[*]}"

  if [[ ${#avail_engines[@]} -ne ${#src_engines[@]} ]]; then
    exit_error "Number of engines reported by CLI does not match number of engine directories."
  fi

  # The items in both arrays should perfectly match since they are sorted
  if [[ "${avail_engines[*]}" != "${src_engines[*]}" ]]; then
    exit_error "Available engines do not match declared engines."
  fi

  log_info "✓ Engine lists match successfully"

  log_info "Querying individual engines..."
  for engine in "${src_engines[@]}"; do
    log_info "Querying engine: $engine"
    "$SNAP_NAME" engine "$engine" >/dev/null
  done
}

get_curr_engine() {
  "$SNAP_NAME" status --format=json | jq -r '.engine'
}

test_engine_switching() {
  local target_engine="$1"

  log_section "Engine Switching Tests"

  log_info "Checking status..."
  "$SNAP_NAME" status

  log_info "Showing current engine..."
  "$SNAP_NAME" engine

  log_info "Testing engine switch..."
  if ! "$SNAP_NAME" use-engine "$target_engine" --assume-yes; then
    exit_error "Failed to switch to engine: $target_engine"
  fi

  log_info "Verifying engine switch via status command..."
  local curr_engine
  curr_engine=$(get_curr_engine)
  if [[ "$curr_engine" != "$target_engine" ]]; then
    exit_error "Current engine from status command ($curr_engine) does not match expected engine ($target_engine)."
  fi
  log_info "✓ Engine switch verified via status command"
}

test_automatic_engine_selection() {
  log_section "Automatic Engine Selection Test"

  log_info "Running: $SNAP_NAME use-engine --auto"
  if ! "$SNAP_NAME" use-engine --auto --assume-yes; then
    exit_error "Automatic engine selection command failed."
  fi
  log_info "Automatic engine selection command completed successfully."
}

# =============================================================================
# MAIN EXECUTION FUNCTION
# =============================================================================

main() {
  local target_engine="$1"
  local features="${2:-}"

  log_section "Starting Smoke Tests"
  log_info "Running tests against snap: $SNAP_NAME"
  log_info "Selected engine: $target_engine"
  log_info "Features: ${features:-<none>}"

  # Pre-flight checks
  local server_port
  server_port=$("$SNAP_NAME" get http.port)
  check_port_listening "$server_port"

  # Run all test suites
  test_snap_installation
  test_configuration_management
  test_engine_listing
  test_automatic_engine_selection
  test_engine_switching "$target_engine"
  test_features "$features"

  log_section "All Smoke Tests Completed Successfully!"
}

# =============================================================================
# SCRIPT ENTRY POINT
# =============================================================================

# Validation
check_root_privileges
check_for_curl
check_for_jq

# Parse arguments
SNAP_NAME=""
ENGINE=""
FEATURES=""
positional_args=()

for arg in "$@"; do
  case "$arg" in
  --features=*)
    FEATURES="${arg#*=}"
    ;;
  -*)
    usage
    exit_error "Unknown option: $arg"
    ;;
  *)
    positional_args+=("$arg")
    ;;
  esac
done

SNAP_NAME="${positional_args[0]:-}"
ENGINE="${positional_args[1]:-}"

validate_arguments
validate_features "$FEATURES"

# Run main function
main "$ENGINE" "$FEATURES"
