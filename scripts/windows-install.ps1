# Windows Install Script
param ($token)

# Constants
$Constants = @{
    ServiceName           = "n2x-node"
    InstallationFolder    = "C:\Program Files\n2x"
    N2xNodeBinary         = "n2x-node.exe"
    N2xNodeBinaryChecksum = "n2x-node.exe_checksum.sha256"
    WintunBinary          = "wintun.dll"
    WintunVersion         = "wintun-0.14.1"
    ConfigFile            = "n2x-node.yml"
    UriN2x                = "https://dl.n2x.io/binaries/stable/latest/windows/amd64/n2x-node.exe"
    UriN2xChecksum        = "https://dl.n2x.io/binaries/stable/latest/windows/amd64/n2x-node.exe_checksum.sha256"
    UriWintun             = "https://www.wintun.net/builds/wintun-0.14.1.zip"
}

# Functions

## Get-TimeStamp Function
function Get-TimeStamp {
    return "{0:yyyy/MM/dd} {0:HH:mm:ss.fff}" -f (Get-Date)
}

## Write-Log Function
function Write-Log {
    param (
        [Parameter(Mandatory=$true, Position=0)]
        [string] $LogLevel,
        [Parameter(Mandatory=$true, Position=1)]
        [string] $Message
    )

    $colors = @{
        info = "White"
        warn = "Yellow"
        error = "Red"
    }

    if ($colors.ContainsKey($LogLevel)) {
        Write-Host "[ $LogLevel] " -ForegroundColor $colors[$LogLevel] -NoNewLine
    }

    Write-Host "$(Get-TimeStamp) $Message" -ForegroundColor White
}

## Download Function
function Download {
    param (
        [Parameter(Mandatory=$true, Position=0)]
        [string] $Uri,
        [Parameter(Mandatory=$true, Position=1)]
        [string] $OutFile
    )

    try {
        $ProgressPreference = 'SilentlyContinue'
        Invoke-WebRequest -Uri $Uri -OutFile $OutFile
    } catch {
        Write-Log -LogLevel "error" -Message ("Failed to download " + $Uri + ": " + $_)
        exit 1
    } finally {
        $ProgressPreference = 'Continue'
    }
}

## Test-Hash Function
function Test-Hash {
    param(
        [Parameter(Mandatory=$true, Position=0)]
        [string] $file,
        [Parameter(Mandatory=$true, Position=1)]
        [string] $hash
    )

    if ((Get-FileHash $file -Algorithm SHA256).Hash.ToUpper() -eq $hash.ToUpper()) {
        Write-Log -LogLevel "info" -Message "Checksum validation succeeded for $file"
    } else {
        Write-Log -LogLevel "error" -Message "Checksum validation failed for $file"
        exit 1
    }
}

## Test-Administrator Function
function Test-Administrator {
    [Security.Principal.WindowsPrincipal]$user = [Security.Principal.WindowsIdentity]::GetCurrent()
    return $user.IsInRole([Security.Principal.WindowsBuiltinRole]::Administrator)
}

## DownloadAndValidate-Binary Function
function DownloadAndValidate-Binary {
    param (
        [string]$binaryPath,
        [string]$checksumPath,
        [string]$binaryUri,
        [string]$checksumUri
    )

    # Download n2x-node binary and checksum
    Write-Log -LogLevel "info" -Message "Downloading n2x-node binary and checksum."
    Download -Uri $binaryUri -OutFile $binaryPath
    Download -Uri $checksumUri -OutFile $checksumPath

    $downloadedHash = (Get-Content $checksumPath).Split(" ")[0]
    Test-Hash -file $binaryPath -hash $downloadedHash
    Write-Log -LogLevel "info" -Message "n2x-node binary validated successfully."
}

# Main Script

## Ensure Administrator Rights
if (-not (Test-Administrator)) {
    Write-Log -LogLevel "error" -Message "This script must be executed as Administrator!"
    exit 1
}

# Ensure token is provided
if (-not $token) {
    Write-Log -LogLevel "error" -Message "Token parameter is required."
    exit 1
}

# Create installation folder if it doesn't exist
if (-not (Test-Path $Constants.InstallationFolder -PathType Container)) {
    Write-Log -LogLevel "info" -Message "Creating installation folder."
    New-Item -ItemType Directory -Force -Path $Constants.InstallationFolder | Out-Null
}

# Create or replace configuration file
$configFilePath = "$($Constants.InstallationFolder)\$($Constants.ConfigFile)"
if (Test-Path $configFilePath) {
    Move-Item -Path $configFilePath -Destination "$($Constants.InstallationFolder)\old-$($Constants.ConfigFile)" -Force
}
New-Item -Path $Constants.InstallationFolder -Name $Constants.ConfigFile -ItemType "file" -Value "Token: $token" | Out-Null
Write-Log -LogLevel "info" -Message "Configuration file created."

# Check if n2x-node binary exits
$n2xNodeBinaryPath = "$($Constants.InstallationFolder)\$($Constants.N2xNodeBinary)"
$checksumPath = "$($Constants.InstallationFolder)\$($Constants.N2xNodeBinaryChecksum)"

if (Test-Path $n2xNodeBinaryPath) {

    do {
        $response = (Read-Host "n2x-node binary already exists. Replace? (Y/N)").ToUpper()
        if ($response -notmatch "^(Y|N)$") {
            Write-Host "Invalid input. Please enter 'Y' or 'N'." -ForegroundColor Yellow
        }
    } while ($response -notmatch "^(Y|N)$")

    if ($response -eq "Y") {

        # Stop and Uninstall n2x-node service
        $arrService = Get-Service -Name $Constants.ServiceName -ErrorAction SilentlyContinue
        if ($arrService -and $arrService.Status -eq 'Running') {
            Write-Log -LogLevel "info" -Message "Stopping $($Constants.ServiceName) service."
            Stop-Service $Constants.ServiceName -Force
        }

        if ($arrService) {
            Write-Log -LogLevel "info" -Message "Uninstalling $($Constants.ServiceName) service."
            & "$n2xNodeBinaryPath" service-uninstall
        }

        DownloadAndValidate-Binary -binaryPath $n2xNodeBinaryPath -checksumPath $checksumPath -binaryUri $Constants.UriN2x -checksumUri $Constants.UriN2xChecksum

    } else {
        Write-Log -LogLevel "info" -Message "Binary replacement skipped."
        return
} else {
    Write-Log -LogLevel "info" -Message "Binary not found. Downloading n2x-node binary and checksum."
    DownloadAndValidate-Binary -binaryPath $n2xNodeBinaryPath -checksumPath $checksumPath -binaryUri $Constants.UriN2x -checksumUri $Constants.UriN2xChecksum
}

# Download and extract Wintun DLL
$wintunBinaryPath = "$($Constants.InstallationFolder)\$($Constants.WintunBinary)"
if (-not (Test-Path $wintunBinaryPath)) {
    Write-Log -LogLevel "info" -Message "Downloading Wintun DLL."
    $wintunZipPath = "$($Constants.InstallationFolder)\$($Constants.WintunVersion).zip"
    Download -Uri $Constants.UriWintun -OutFile $wintunZipPath

    Write-Log -LogLevel "info" -Message "Extracting Wintun DLL."
    Expand-Archive -Path $wintunZipPath -DestinationPath "$($Constants.InstallationFolder)\$($Constants.WintunVersion)"

    Move-Item -Path "$($Constants.InstallationFolder)\$($Constants.WintunVersion)\wintun\bin\amd64\$($Constants.WintunBinary)" -Destination $wintunBinaryPath -Force
    Remove-Item "$($Constants.InstallationFolder)\$($Constants.WintunVersion)" -Recurse -Force
    Remove-Item $wintunZipPath
}

# Install and Start n2x-node service

$service = Get-Service -Name $Constants.ServiceName -ErrorAction SilentlyContinue

# Check if the service exists
if (-not $service) {
    # Service is not installed
    Write-Log -LogLevel "info" -Message "$($Constants.ServiceName) service is not installed. Installing..."
    & "$n2xNodeBinaryPath" service-install
} else {
    Write-Log -LogLevel "info" -Message "$($Constants.ServiceName) service is already installed."
}

# Check if the service is running
if ($service -and $service.Status -eq 'Running') {
    Write-Log -LogLevel "info" -Message "$($Constants.ServiceName) service is already running."
} else {
    Write-Log -LogLevel "info" -Message "Starting $($Constants.ServiceName) service."
    Start-Service $Constants.ServiceName
}

# Check if the scheduled task already exists
$taskName = "n2xNodeStartupTask"
$existingTask = Get-ScheduledTask -TaskName $taskName -ErrorAction SilentlyContinue

if (-not $existingTask) {
    Write-Log -LogLevel "info" -Message "Scheduled task $taskName does not exist. Creating it."

   # Define the scheduled task action
    $action = New-ScheduledTaskAction -Execute "PowerShell.exe" -Argument "-File $n2xNodeBinaryPath"

    # Define the trigger to run at system startup
    $trigger = New-ScheduledTaskTrigger -AtStartup

    # Define task settings for reliability
    $settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable -RunOnlyIfNetworkAvailable -RestartCount 3 -RestartInterval (New-TimeSpan -Minutes 1)

    # Register the task
    Register-ScheduledTask -Action $action -Trigger $trigger -TaskName $taskName -Description "Runs n2x-node at startup" -Settings $settings -User "SYSTEM" -RunLevel Highest

    Write-Log -LogLevel "info" -Message "Scheduled task $taskName created successfully."
} else {
    Write-Log -LogLevel "info" -Message "Scheduled task $taskName already exists. Skipping creation."
}

Write-Log -LogLevel "info" -Message "Script completed successfully."