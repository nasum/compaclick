[Setup]
AppName=CompaClick
AppVersion=1.0
AppPublisher=Nasum
DefaultDirName={autopf}\CompaClick
DefaultGroupName=CompaClick
OutputDir=Output
OutputBaseFilename=CompaClick_Setup
Compression=lzma
SolidCompression=yes
PrivilegesRequired=admin

[Files]
Source: "CompaClick.exe"; DestDir: "{app}"; Flags: ignoreversion

[Run]
Filename: "{app}\CompaClick.exe"; Parameters: "install"; Flags: runhidden

[UninstallRun]
Filename: "{app}\CompaClick.exe"; Parameters: "uninstall"; Flags: runhidden
