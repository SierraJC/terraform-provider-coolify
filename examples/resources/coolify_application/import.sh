#!/bin/bash
# Import script for coolify_application

# Usage: ./import.sh <application-uuid>
# Example: ./import.sh rg8ks8c

if [ -z "$1" ]; then
  echo "Usage: $0 <application-uuid>"
  exit 1
fi

terraform import coolify_application.example "$1"

