dotnet.exe build

$Processes = @(1, 2, 3, 4, 5, 6) | ForEach-Object {
    Start-Process "dotnet.exe" -ArgumentList @("run", "--", "run", "--id", $_, "--port", "1000$_") -NoNewWindow
}

try {
    $Routes = @(
        @{ Source = 1; Dest = 2; Cost = 4 },
        @{ Source = 2; Dest = 1; Cost = 4 },
        @{ Source = 1; Dest = 5; Cost = 1 },
        @{ Source = 5; Dest = 1; Cost = 1 },
        @{ Source = 2; Dest = 3; Cost = 3 },
        @{ Source = 3; Dest = 2; Cost = 3 },
        @{ Source = 2; Dest = 5; Cost = 2 },
        @{ Source = 3; Dest = 2; Cost = 2 },
        @{ Source = 2; Dest = 6; Cost = 1 },
        @{ Source = 6; Dest = 2; Cost = 1 },
        @{ Source = 3; Dest = 4; Cost = 4 },
        @{ Source = 4; Dest = 3; Cost = 4 },
        @{ Source = 3; Dest = 6; Cost = 1 },
        @{ Source = 6; Dest = 3; Cost = 1 },
        @{ Source = 4; Dest = 6; Cost = 3 },
        @{ Source = 6; Dest = 4; Cost = 3 },
        @{ Source = 5; Dest = 6; Cost = 3 },
        @{ Source = 6; Dest = 5; Cost = 3 }
    )

    $Routes | ForEach-Object {
        $Source = $_.Source
        $Dest = $_.Dest
        $Cost = $_.Cost
        Write-Host "Configuring route from $Source to $Dest with cost $Cost"
        dotnet.exe run -- configure --id $Source --port "1000$Source" --src $Source --dest $Dest --dport "1000$Dest" --cost $Cost
    }

    Start-Sleep -Seconds 16

    Start-Sleep -Seconds 2
}
finally {
    $Processes | ForEach-Object {
        Stop-Process -Id $_.Id -Force
        Write-Host $_.StandardOutput
        Wait-Process -Id $_.Id
    }
}