$ErrorActionPreference = "Stop"

$repo = "Alia5/steaminputdb.com"
$apiUrl = "https://api.github.com/repos/$repo/releases/latest"

Write-Host "Fetching latest SteamInputDB Buddy release..."
$releaseData = Invoke-RestMethod -Uri $apiUrl -ErrorAction Stop
$version = $releaseData.tag_name

if (-not $version) {
    Write-Host "Error: Could not fetch release info" -ForegroundColor Red
    exit 1
}

Write-Host "Version: $version"

if (-not [Environment]::Is64BitOperatingSystem) {
    Write-Host "Error: Only 64-bit Windows is supported" -ForegroundColor Red
    exit 1
}

$binaryName = "steaminputdb-buddy-windows-amd64.exe"
$downloadUrl = "https://github.com/$repo/releases/download/$version/$binaryName"

$installDir = Join-Path $env:LOCALAPPDATA "SteamInputDB"
$installPath = Join-Path $installDir "steaminputdb-buddy.exe"

$isUpdate = Test-Path $installPath
$skipInstall = $false

function Get-BuddyVersion($path) {
    try {
        $help = & $path --help 2>$null
        $match = ($help | Select-String -Pattern 'SteamInputDB Buddy - v([^\s]+)' | Select-Object -First 1)
        if ($match) {
            return $match.Matches[0].Groups[1].Value
        }
    }
    catch { }
    return $null
}

$oldVersion = "unknown"
if ($isUpdate) {
    Write-Host "Existing installation detected at $installPath"
    $oldVersionRaw = Get-BuddyVersion $installPath
    if ($oldVersionRaw) { $oldVersion = $oldVersionRaw }
    Write-Host "Installed version: $oldVersion"
}

Write-Host "Downloading from: $downloadUrl"
$tempDir = New-TemporaryFile | ForEach-Object { Remove-Item $_; New-Item -ItemType Directory -Path $_ }

try {
    $tempBin = Join-Path $tempDir "steaminputdb-buddy.exe"
    Invoke-WebRequest -Uri $downloadUrl -OutFile $tempBin -ErrorAction Stop

    $newVersion = "unknown"
    $newVersionRaw = Get-BuddyVersion $tempBin
    if ($newVersionRaw) { $newVersion = $newVersionRaw }
    Write-Host "Downloaded version: $newVersion"

    if ($isUpdate) {
        if ($newVersion -eq $oldVersion -and $newVersion -ne "unknown") {
            Write-Host "Already at latest version. Skipping."
            $skipInstall = $true
        }
    }

    if (-not $skipInstall) {
        New-Item -ItemType Directory -Path $installDir -Force | Out-Null

        if ($isUpdate) {
            Write-Host "Stopping running instance(s)..."
            $procs = Get-CimInstance Win32_Process -ErrorAction SilentlyContinue |
                Where-Object { $_.ExecutablePath -eq $installPath }
            if ($procs) {
                foreach ($p in $procs) {
                    try {
                        Stop-Process -Id $p.ProcessId -Force -ErrorAction SilentlyContinue
                    }
                    catch { }
                }
                Start-Sleep -Milliseconds 500
            }
        }

        Write-Host "Installing to $installPath..."
        Copy-Item $tempBin $installPath -Force

        Write-Host "Running install..."
        Start-Process -FilePath $installPath -ArgumentList "install", "--in-place", "--show-ui" -Wait -NoNewWindow
    }

    Write-Host ""
    if ($skipInstall) {
        Write-Host "SteamInputDB Buddy is already up to date." -ForegroundColor Green
    }
    elseif ($isUpdate) {
        Write-Host "SteamInputDB Buddy updated successfully!" -ForegroundColor Green
    }
    else {
        Write-Host "SteamInputDB Buddy installed successfully!" -ForegroundColor Green
    }
    Write-Host "Binary: $installPath"
}
finally {
    Remove-Item -Recurse -Force $tempDir -ErrorAction SilentlyContinue
}
