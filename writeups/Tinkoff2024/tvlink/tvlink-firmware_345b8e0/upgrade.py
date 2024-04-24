import sys
import time

print("""░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░▒▓███████▓▒░░▒▓█▓▒░░▒▓█▓▒░ 
   ░▒▓█▓▒░   ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
   ░▒▓█▓▒░    ░▒▓█▓▒▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
   ░▒▓█▓▒░    ░▒▓█▓▒▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓███████▓▒░  
   ░▒▓█▓▒░     ░▒▓█▓▓█▓▒░ ░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
   ░▒▓█▓▒░     ░▒▓█▓▓█▓▒░ ░▒▓█▓▒░      ░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
   ░▒▓█▓▒░      ░▒▓██▓▒░  ░▒▓████████▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
""", file=sys.stderr)

print("CERTIFIED FIRWARE UPGRADE SELF-SERVICE PACK", file=sys.stderr)
print("Version 20240420\n", file=sys.stderr, flush=True)
time.sleep(1)

print("""DISCLAIMER
      
Manufacturer disclaims all warranties, express or implied, including without
limitation, the implied warranties of merchantability and fitness for a
particular purpose. Manufacturer shall not be liable for errors contained
herein or for incidental or consequential damages in connection with the
furnishing, performance, or use of this material.
      
ALWAYS CHECK THE SOURCE OF FIRMWARE BEFORE UPGRADING. USE ONLY OUR OFFICIAL
WEBSITE FOR DOWNLOADS. DO NOT TRUST THIRD PARTY WEBSITES.
      
By proceeding with the upgrade, you agree to the terms and conditions set by
the manufacturer and to processing of your personal data in accordance with our
Privacy Policy.
      
Upgrade process will continue in 5 seconds.
""", flush=True)

time.sleep(5)

print("Starting upgrade process...", flush=True)
print("Please do not disconnect supply from the TV box", flush=True)

time.sleep(1)

try:
    import tvlink
except ImportError:
    print("Upgrade failed: current firmware not found", file=sys.stderr, flush=True)
    print("Please contact our support center", file=sys.stderr, flush=True)
    sys.exit(1)

current_version = tvlink.__version__
if current_version >= "20231207":
    print("You have already installed this or newer version: downgrade is not possible", file=sys.stderr, flush=True)
    sys.exit(1)

print("Check success, continuing with the upgrade...")
print("Please do not forget to remember your license key in case of failure")

print(" |-------------------------------------------------------------|")
print(" | ", end="")
with open("/license/key.txt") as key:
    print(key.read().strip(), end="")
print(" | ")
print(" |-------------------------------------------------------------|", flush=True)

import subprocess

PATCHES = ["""--- tvlink/app.py	2023-12-06 13:37:00.000000000 +1100
+++ tvlink-20240420/app.py	2024-04-20 13:37:00.000000000 +1100
@@ -110,7 +110,7 @@
     if "auth" not in request.session:
         return RedirectResponse("/LOGIN.XHTML", status_code=303)
 
-    tv_enabled = False  # TODO
+    tv_enabled = True
     return templates.TemplateResponse(
         request=request,
         name="tv.html",
"""]

for idx, patch in enumerate(PATCHES, start=1):
    time.sleep(1)
    print(f"Applying patch {idx} of {len(PATCHES)}...", file=sys.stderr)
    try:
        subprocess.run(["patch", "-p0"], input=patch, capture_output=True, text=True, check=True, cwd="/app")
    except subprocess.CalledProcessError as exc:
        print(f"Patch {idx} failed with error {exc.returncode}", file=sys.stderr)
        print("SYSTEM MAY BE DAMAGED, PLEASE CONTACT SUPPORT CENTER FOR RECOVERY", file=sys.stderr, flush=True)
        sys.exit(1)

print(f"UPGRADE SUCCESSFUL, REBOOTING THE DEVICE", flush=True)
subprocess.run(["reboot"])
