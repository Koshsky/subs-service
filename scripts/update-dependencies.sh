#!/bin/bash

set -e

# Update dependencies for each service
services=("auth-service" "core-service" "notification-service")

# Check for go.work file
if [ ! -f "go.work" ]; then
	echo "go.work file not found, creating from example..."
	if [ -f "go.work.example" ]; then
		cp go.work.example go.work
		echo "go.work file created from example"
	else
		echo "go.work.example not found"
		exit 1
	fi
fi

# Update dependencies for each service

for service in "${services[@]}"; do
	if [ -d "$service" ]; then
		echo "Updating dependencies for $service..."
		cd "$service"

		# Check for go.mod
		if [ ! -f "go.mod" ]; then
			echo "go.mod not found in $service"
			exit 1
		fi

		# Update dependencies
		go mod tidy
		go mod download

		echo "Dependencies for $service updated"
		cd ..
	else
		echo "Directory $service not found"
	fi
done

# Update workspace
echo "Updating workspace dependencies..."
go work sync
echo "Workspace dependencies updated"