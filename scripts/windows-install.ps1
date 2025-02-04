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
        info = "Blue"
        warn = "Yellow"
        error = "Red"
    }

    if ($colors.ContainsKey($LogLevel)) {
        Write-Host "[$LogLevel] " -ForegroundColor $colors[$LogLevel] -NoNewLine
    }

    Write-Host " $(Get-TimeStamp) $Message" -ForegroundColor DarkGray
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

## Ensure Administrator Rights
if (-not (Test-Administrator)) {
    Write-Log -LogLevel "error" -Message "This script must be executed as Administrator!"
    exit 1
}

# Main Script

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

# Download and validate n2x-node binary
$n2xNodeBinaryPath = "$($Constants.InstallationFolder)\$($Constants.N2xNodeBinary)"
if (-not (Test-Path $n2xNodeBinaryPath) -or (Read-Host "n2x-node binary already exists. Replace? (Y/N)" -eq "Y")) {
    Write-Log -LogLevel "info" -Message "Downloading n2x-node binary."
    Download -Uri $Constants.UriN2x -OutFile $n2xNodeBinaryPath

    Write-Log -LogLevel "info" -Message "Downloading n2x-node checksum."
    $checksumPath = "$($Constants.InstallationFolder)\$($Constants.N2xNodeBinaryChecksum)"
    Download -Uri $Constants.UriN2xChecksum -OutFile $checksumPath

    $downloadedHash = (Get-Content $checksumPath).Split(" ")[0]
    Test-Hash -file $n2xNodeBinaryPath -hash $downloadedHash
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

# Manage n2x-node service
$arrService = Get-Service -Name $Constants.ServiceName -ErrorAction SilentlyContinue
if ($arrService -and $arrService.Status -eq 'Running') {
    Write-Log -LogLevel "info" -Message "Stopping $($Constants.ServiceName) service."
    Stop-Service $Constants.ServiceName -Force
}

if ($arrService) {
    Write-Log -LogLevel "info" -Message "Uninstalling $($Constants.ServiceName) service."
    & "$n2xNodeBinaryPath" service-uninstall
}

Write-Log -LogLevel "info" -Message "Installing $($Constants.ServiceName) service."
& "$n2xNodeBinaryPath" service-install

Write-Log -LogLevel "info" -Message "Starting $($Constants.ServiceName) service."
Start-Service $Constants.ServiceName

# Schedule Task
$taskName = "n2xNodeStartupTask"
if (Get-ScheduledTask -TaskName $taskName -ErrorAction SilentlyContinue) {
    Unregister-ScheduledTask -TaskName $taskName -Confirm:$false
}

$action = New-ScheduledTaskAction -Execute "PowerShell.exe" -Argument "-File $n2xNodeBinaryPath"
$trigger = New-ScheduledTaskTrigger -AtStartup
$settings = New-ScheduledTaskSettingsSet -AllowStartIfOnBatteries -DontStopIfGoingOnBatteries -StartWhenAvailable

Register-ScheduledTask -Action $action -Trigger $trigger -TaskName $taskName -Description "Runs n2x-node at startup" -Settings $settings

Write-Log -LogLevel "info" -Message "Script completed successfully."