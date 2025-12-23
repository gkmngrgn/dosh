# PowerShell install script for DOSH on Windows
# Usage: .\install.ps1

$ErrorActionPreference = "Stop"

# Get system information
$os = "windows"
$architecture = if ([Environment]::Is64BitOperatingSystem) { "x86_64" } else { "386" }
$downloadUrl = "https://github.com/miratcan/dosh/releases/latest/download/dosh-$os-$architecture.exe"
$tempDir = [System.IO.Path]::GetTempPath() + [System.Guid]::NewGuid().ToString()
$localDir = "$env:USERPROFILE\.local"
$binDir = "$localDir\bin"
$binFile = "$binDir\dosh.exe"

Write-Host "Operating System: $os"
Write-Host "Architecture: $architecture"
Write-Host "Temporary directory: $tempDir"
Write-Host "Download URL: $downloadUrl"

Write-Host "`nSTEP 1: Downloading DOSH..."

# Create temporary directory
New-Item -ItemType Directory -Path $tempDir -Force | Out-Null
$tempFile = "$tempDir\dosh.exe"

try {
    # Download the file
    Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UseBasicParsing

    Write-Host "`nSTEP 2: Installing DOSH CLI..."

    # Check if binary already exists and backup if it does
    if (Test-Path $binFile) {
        $backupFile = "$tempDir\dosh.old.exe"
        Move-Item $binFile $backupFile
        Write-Host "Existing installation backed up to: $backupFile"
    } else {
        # Create the local bin directory if it doesn't exist
        if (-not (Test-Path $binDir)) {
            New-Item -ItemType Directory -Path $binDir -Force | Out-Null
        }
    }

    # Move the downloaded file to the bin directory
    Move-Item $tempFile $binFile

    # Add to PATH if not already there
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($currentPath -notlike "*$binDir*") {
        $newPath = "$currentPath;$binDir"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Host "Added $binDir to your PATH environment variable."
        Write-Host "Please restart your terminal or run 'refreshenv' to use the new PATH."
    }

    Write-Host "`nSTEP 3: Done! You can delete the temporary directory if you want:"
    Write-Host "$tempDir"
    Write-Host "`nDOSH has been installed to: $binFile"
    Write-Host "You may need to restart your terminal for the PATH changes to take effect."

} catch {
    Write-Error "Installation failed: $_"
    # Clean up temporary directory on error
    if (Test-Path $tempDir) {
        Remove-Item $tempDir -Recurse -Force
    }
    exit 1
} finally {
    # Note: We don't auto-delete the temp directory in case user wants to recover backup
}
