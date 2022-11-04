!define PRODUCT_VERSION "0.3.3.0"
!define PUBLISHER "Forest Soft"
!define VERSION "0.3.3.0"
!define PRODUCT_URL "https://github.com/Forest33/warthog"

VIProductVersion "${PRODUCT_VERSION}"
VIFileVersion "${VERSION}"
VIAddVersionKey FileVersion "${VERSION}"
VIAddVersionKey FileDescription "Warthog cross platform gRPC client"
VIAddVersionKey LegalCopyright ""

VIAddVersionKey ProductName "Warthog"
VIAddVersionKey Comments "Installs Warthog."
VIAddVersionKey CompanyName "${PUBLISHER}"
VIAddVersionKey ProductVersion "${PRODUCT_VERSION}"
VIAddVersionKey InternalName "Warthog"
;VIAddVersionKey LegalTrademarks " "
;VIAddVersionKey PrivateBuild ""
;VIAddVersionKey SpecialBuild ""

!include "MUI.nsh"
!define MUI_ICON "app.ico"
!define MUI_UNICON "app.ico"

# The name of the installer
Name "Warthog ${VERSION}"

# The file to write
OutFile "setup.exe"

; Build Unicode installer
Unicode True

# The default installation directory
InstallDir $PROGRAMFILES\Warthog\

; -------
; Registry key to check for directory (so if you install again, it will
; overwrite the old one automatically)
InstallDirRegKey HKLM "Software\Warthog" "Install_Dir"
; -------

# The text to prompt the user to enter a directory
DirText "This will install Warthog on your computer. Choose a directory"

!insertmacro MUI_LANGUAGE "English"

#--------------------------------

# The stuff to install
Section "" #No components page, name is not important

# Set output path to the installation directory.
SetOutPath $INSTDIR

# Put a file there
File Warthog.exe
File app.ico

# Tell the compiler to write an uninstaller and to look for a "Uninstall" section
WriteUninstaller $INSTDIR\Uninstall.exe

CreateDirectory "$SMPROGRAMS\Warthog"
CreateShortCut "$SMPROGRAMS\Warthog\Warthog.lnk" "$INSTDIR\Warthog.exe"
CreateShortCut "$SMPROGRAMS\Warthog\Uninstall.lnk" "$INSTDIR\Uninstall.exe"

; -------
; Write the installation path into the registry
WriteRegStr HKLM SOFTWARE\Warthog "Install_Dir" "$INSTDIR"

; Write the uninstall keys for Windows
WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Warthog" "DisplayName" "Warthog"
WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Warthog" "UninstallString" '"$INSTDIR\uninstall.exe"'
WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Warthog" "Version" "${PRODUCT_VERSION}"
WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Warthog" "DisplayVersion" "${PRODUCT_VERSION}"
WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Warthog" "Publisher" "${PUBLISHER}"
WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Warthog" "DisplayIcon" "$INSTDIR\app.ico"
WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Warthog" "HelpLink" "${PRODUCT_URL}"
WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Warthog" "URLUpdateInfo" "${PRODUCT_URL}"
WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Warthog" "NoModify" 1
WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Warthog" "NoRepair" 1

WriteUninstaller "$INSTDIR\uninstall.exe"
; -------

SectionEnd # end the section

# The uninstall section
Section "Uninstall"

; -------
; Remove registry keys
DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Warthog"
DeleteRegKey HKLM SOFTWARE\Warthog
; -------

RMDir /r "$SMPROGRAMS\Warthog"
RMDir /r "$PROFILE\Warthog"
RMDir /r "$APPDATA\Warthog"
RMDir /r $INSTDIR

SectionEnd
