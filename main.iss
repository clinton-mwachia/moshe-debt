[Setup]
AppName=Moshe Debt
AppVersion=1.0.6
DefaultDirName={pf}\Moshe Debt
DefaultGroupName=moshe-debt
OutputDir=output
OutputBaseFilename=moshe-debt
Compression=lzma
SolidCompression=yes
DisableDirPage=no
LicenseFile="LICENSE.txt"

[Files]
Source: "moshe-debt.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "assets\*"; DestDir: "{app}\assets"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "README.txt"; DestDir: "{app}"; Flags: isreadme

[Icons]
Name: "{group}\Moshe Debt"; Filename: "{app}\moshe-debt.exe"
Name: "{group}\Uninstall Moshe Debt"; Filename: "{uninstallexe}"
Name: "{commondesktop}\Moshe Debt"; Filename: "{app}\moshe-debt.exe"



