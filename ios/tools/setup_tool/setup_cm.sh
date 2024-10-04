#!/bin/bash

# Set strict mode
#set -euo pipefail

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colorful messages
print_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

# Welcome message
echo "======================================="
print_message "$BLUE" "🚀 Welcome to Critical Moments Setup 🚀"
echo "======================================="
print_message "$GREEN" "This script will guide you through setting up your iOS project with the Critical Moments SDK."
echo

# Function to check if a file is a valid Info.plist
is_valid_info_plist() {
    local file=$1
    [[ -f "$file" && "$(basename "$file")" == "Info.plist" ]]
}

# Check for command-line argument
if [[ $# -eq 1 ]]; then
    if is_valid_info_plist "$1"; then
        info_plist="$1"
        print_message "$GREEN" "Using provided Info.plist: $info_plist"
    else
        print_message "$RED" "Error: The provided path is not a valid Info.plist file."
        exit 1
    fi
else
    # Find Info.plist files
    info_plist_files=()
    
    # Check current directory
    if [[ -f "Info.plist" ]]; then
        info_plist_files+=("$(pwd)/Info.plist")
    fi
    
    # Check subdirectories (excluding those ending with "Test")
    for dir in */; do
        if [[ ! "$dir" =~ Test/$ ]]; then
            if [[ -f "${dir}Info.plist" ]]; then
                info_plist_files+=("$(pwd)/${dir}Info.plist")
            fi
        fi
    done
    
    # Process found Info.plist files
    case ${#info_plist_files[@]} in
        0)
            print_message "$RED" "Error: No Info.plist files found in this directory or its immediate subdirectories."
            print_message "$RED" "Please run the script again in your project directory, or with the Info.plist path as an argument."
            exit 1
            ;;
        1)
            print_message "$YELLOW" "We think we found your app's Info.plist:"
            print_message "$BLUE" "${info_plist_files[0]}"
            read -p "Is this the correct Info.plist file for your app? (Y/N): " confirm
            if [[ $confirm =~ ^[Yy]$ ]]; then
                info_plist="${info_plist_files[0]}"
            else
                print_message "$RED" "Please run the script again in your project directory, or with the Info.plist path as an argument."
                exit 1
            fi
            ;;
        *)
            print_message "$YELLOW" "Multiple Info.plist files found. Please select the correct one:"
            for i in "${!info_plist_files[@]}"; do
                echo "[$((i+1))] ${info_plist_files[$i]}"
            done
            while true; do
                read -p "Enter the number of the correct Info.plist: " selection
                if [[ "$selection" =~ ^[0-9]+$ ]] && [ "$selection" -ge 1 ] && [ "$selection" -le "${#info_plist_files[@]}" ]; then
                    info_plist="${info_plist_files[$((selection-1))]}"
                    break
                else
                    print_message "$RED" "Invalid selection. Please try again."
                fi
            done
            ;;
    esac
fi

print_message "$GREEN" "Using Info.plist: $info_plist"

# Function to check if an array contains a value
array_contains() {
    local array="$1"
    local value="$2"
    echo "$array" | jq -e --arg v "$value" 'contains([$v])' >/dev/null
}

# Function to safely extract plist values
safe_extract_plist() {
    local key=$1
    local type=$2
    local value
    local exit_code

    # Disable ALL strict mode settings
    set +e
    set +u
    set +o pipefail

    value=$(plutil -extract "$key" $type -o - "$info_plist" 2>/dev/null)
    exit_code=$?

    # Restore previous options
    set +euo pipefail

    if [ $exit_code -ne 0 ]; then
        echo ""
        return 1
    else
        echo "$value"
        return 0
    fi
}

# Initialize an array to store needed changes
changes_needed=()

# Check UIBackgroundModes
background_modes=$(safe_extract_plist "UIBackgroundModes" json)
if [ $? -ne 0 ] || ! (array_contains "$background_modes" "fetch" && array_contains "$background_modes" "processing"); then
    changes_needed+=("Add 'fetch' and 'processing' to UIBackgroundModes")
fi

# Check BGTaskSchedulerPermittedIdentifiers
bg_identifiers=$(safe_extract_plist "BGTaskSchedulerPermittedIdentifiers" json)
if [ $? -ne 0 ] || ! (array_contains "$bg_identifiers" "io.criticalmoments.bg_fetch" && array_contains "$bg_identifiers" "io.criticalmoments.bg_process"); then
    changes_needed+=("Add 'io.criticalmoments.bg_fetch' and 'io.criticalmoments.bg_process' to BGTaskSchedulerPermittedIdentifiers")
fi

# Check NSBluetoothAlwaysUsageDescription
bluetooth_description=$(safe_extract_plist "NSBluetoothAlwaysUsageDescription" raw)
if [ $? -ne 0 ] || [ -z "$bluetooth_description" ]; then
    changes_needed+=("Add NSBluetoothAlwaysUsageDescription")
fi

# If no changes are needed, print a friendly message and exit
if [ ${#changes_needed[@]} -eq 0 ]; then
    print_message "$GREEN" "Great news! Your Info.plist is already set up correctly for Critical Moments. No changes are needed.\nContinue with the next step in the setup guide at docs.criticalmoments.io"
    exit 0
fi


# List changes and get user consent
print_message "$YELLOW" "\n\n\nThe following changes are needed in your Info.plist:"
echo

for change in "${changes_needed[@]}"; do
    case "$change" in
        "Add 'fetch' and 'processing' to UIBackgroundModes")
            print_message "$BLUE" "- $change"
            print_message "$GREEN" "  Reason: Critical Moments uses background processing, making it possible to deliver smart notifications at the right moment."
            ;;
        "Add 'io.criticalmoments.bg_fetch' and 'io.criticalmoments.bg_process' to BGTaskSchedulerPermittedIdentifiers")
            print_message "$BLUE" "- $change"
            print_message "$GREEN" "  Reason: Critical Moments uses background processing, making it possible to deliver smart notifications at the right moment. These task IDs must be registered for our workers to run."
            ;;
        "Add NSBluetoothAlwaysUsageDescription")
            print_message "$BLUE" "- $change"
            print_message "$GREEN" "  Reason: We link bluetooth APIs to allow conditions like checking if the user has BT headphones paired. Apple will detect this, and require a reason string. We'll never call these unless you use BT properties in a condition. The user will never be prompted for bluetooth access as the APIs we use do not require it."
            ;;
    esac
    echo
done

read -p "Do you want to make these changes? (Y/N): " make_changes

if [[ ! $make_changes =~ ^[Yy]$ ]]; then
    print_message "$RED" "Changes not applied. Exiting."
    exit 1
fi

# Make the changes
for change in "${changes_needed[@]}"; do
    case "$change" in
        "Add 'fetch' and 'processing' to UIBackgroundModes")
            current_modes=$(plutil -extract UIBackgroundModes json -o - "$info_plist" 2>/dev/null)
            if [ $? -ne 0 ] || [ "$current_modes" == "null" ]; then
                plutil -replace UIBackgroundModes -json '["fetch", "processing"]' "$info_plist"
            else
                new_modes=$(echo "$current_modes" | jq '. + ["fetch", "processing"] | unique')
                plutil -replace UIBackgroundModes -json "$new_modes" "$info_plist"
            fi
            ;;
        "Add 'io.criticalmoments.bg_fetch' and 'io.criticalmoments.bg_process' to BGTaskSchedulerPermittedIdentifiers")
            current_identifiers=$(plutil -extract BGTaskSchedulerPermittedIdentifiers json -o - "$info_plist" 2>/dev/null)
            if [ $? -ne 0 ] || [ "$current_identifiers" == "null" ]; then
                plutil -replace BGTaskSchedulerPermittedIdentifiers -json '["io.criticalmoments.bg_fetch", "io.criticalmoments.bg_process"]' "$info_plist"
            else
                new_identifiers=$(echo "$current_identifiers" | jq '. + ["io.criticalmoments.bg_fetch", "io.criticalmoments.bg_process"] | unique')
                plutil -replace BGTaskSchedulerPermittedIdentifiers -json "$new_identifiers" "$info_plist"
            fi
            ;;
        "Add NSBluetoothAlwaysUsageDescription")
            plutil -replace NSBluetoothAlwaysUsageDescription -string "Used to show messages when not using peripherals." "$info_plist"
            ;;
    esac
done

print_message "$GREEN" "All changes have been applied successfully to your Info.plist file."

