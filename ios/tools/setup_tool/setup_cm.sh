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
echo
echo
echo "======================================="
print_message "$BLUE" "🚀 Welcome to Critical Moments Setup 🚀"
echo "======================================="
echo
echo
print_message "$GREEN" "This script will guide you through setting up your iOS project with the Critical Moments SDK."
echo
echo
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
    [[ "$array" == *"\"$value\""* ]]
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

# Define change strings as variables
CHANGE_BG_MODES="Update UIBackgroundModes to include 'fetch' and 'processing'"
CHANGE_BG_IDENTIFIERS="Update BGTaskSchedulerPermittedIdentifiers to include 'io.criticalmoments.bg_fetch' and 'io.criticalmoments.bg_process'"
CHANGE_BT_DESCRIPTION="Update NSBluetoothAlwaysUsageDescription, with a reason for using bluetooth"

# Check UIBackgroundModes
background_modes=$(safe_extract_plist "UIBackgroundModes" json)
if [ $? -ne 0 ] || ! (array_contains "$background_modes" "fetch" && array_contains "$background_modes" "processing"); then
    changes_needed+=("$CHANGE_BG_MODES")
fi

# Check BGTaskSchedulerPermittedIdentifiers
bg_identifiers=$(safe_extract_plist "BGTaskSchedulerPermittedIdentifiers" json)
if [ $? -ne 0 ] || ! (array_contains "$bg_identifiers" "io.criticalmoments.bg_fetch" && array_contains "$bg_identifiers" "io.criticalmoments.bg_process"); then
    changes_needed+=("$CHANGE_BG_IDENTIFIERS")
fi

# Check NSBluetoothAlwaysUsageDescription
bluetooth_description=$(safe_extract_plist "NSBluetoothAlwaysUsageDescription" raw)
if [ $? -ne 0 ] || [ -z "$bluetooth_description" ]; then
    changes_needed+=("$CHANGE_BT_DESCRIPTION")
fi

# If no changes are needed, print a friendly message and exit
if [ ${#changes_needed[@]} -eq 0 ]; then
    print_message "$GREEN" "\n\nGreat news! Your Info.plist is already set up correctly for Critical Moments. No changes are needed.\nContinue with the next step in the setup guide at docs.criticalmoments.io"
    print_message "$GREEN" "Happy coding! 🚀"
    exit 0
fi


# List changes and get user consent
print_message "$YELLOW" "\n\n\nThe following $([ ${#changes_needed[@]} -eq 1 ] && echo "change is" || echo "${#changes_needed[@]} changes are") needed in your Info.plist:"
echo

for change in "${changes_needed[@]}"; do
    case "$change" in
        "$CHANGE_BG_MODES")
            print_message "$BLUE" "$change"
            print_message "$GREEN" "↪ Reason: Critical Moments uses background processing, making it possible to deliver smart notifications at the right moment. These modes must be enabled for our workers to run."
            ;;
        "$CHANGE_BG_IDENTIFIERS")
            print_message "$BLUE" "$change"
            print_message "$GREEN" "↪ Reason: Critical Moments uses background processing, making it possible to deliver smart notifications at the right moment. These task IDs must be registered for our workers to run."
            ;;
        "$CHANGE_BT_DESCRIPTION")
            print_message "$BLUE" "$change"
            print_message "$GREEN" "↪ Reason: The SDK links bluetooth APIs to allow checking conditions, such as if the user has BT headphones paired. Apple will detect this, and require a reason string. Note: we never call these APIs unless you use BT properties in your config. The user will not be prompted for bluetooth access, as the APIs we use do not require it. You can change this string after if desired."
            ;;
    esac
    echo
done

read -p "Do you want to make these changes? (Y/N): " make_changes

if [[ ! $make_changes =~ ^[Yy]$ ]]; then
    print_message "$RED" "Changes not applied. Exiting."
    exit 1
fi

# Helper function to add items to an array in plist
add_to_plist_array() {
    local key=$1
    local items=("${@:2}")
    local current_array
    
    current_array=$(safe_extract_plist "$key" json)
    if [ $? -ne 0 ] || [ -z "$current_array" ]; then
        plutil -replace "$key" -json "$(printf '"%s",' "${items[@]}" | sed 's/,$//' | awk '{print "["$0"]"}')" "$info_plist"
    else
        local new_array=$(echo "$current_array" | sed 's/\[//;s/\]//;s/"//g' | tr ',' '\n' | sort -u)
        for item in "${items[@]}"; do
            if ! echo "$new_array" | grep -q "^$item$"; then
                new_array="$new_array"$'\n'"$item"
            fi
        done
        new_array=$(echo "$new_array" | sort -u | sed 's/^/"/;s/$/"/' | tr '\n' ',' | sed 's/,$//')
        new_array="[$new_array]"
        plutil -replace "$key" -json "$new_array" "$info_plist"
    fi
}

# Make the changes
for change in "${changes_needed[@]}"; do
    case "$change" in
        "$CHANGE_BG_MODES")
            add_to_plist_array "UIBackgroundModes" "fetch" "processing"
            ;;
        "$CHANGE_BG_IDENTIFIERS")
            add_to_plist_array "BGTaskSchedulerPermittedIdentifiers" "io.criticalmoments.bg_fetch" "io.criticalmoments.bg_process"
            ;;
        "$CHANGE_BT_DESCRIPTION")
            plutil -replace NSBluetoothAlwaysUsageDescription -string "Used to avoid showing messages when using peripherals." "$info_plist"
            ;;
    esac
done

print_message "$GREEN" "\n\nAll changes have been successfully applied to your Info.plist file.\n\n"

print_message "$BLUE" "🎉 Congratulations! Your Info.plist has been successfully updated for Critical Moments. 🎉"
print_message "$GREEN" "Happy coding! 🚀"
echo