import hashlib
import os
import pathlib
import signal
import subprocess
import sys
import tempfile
import typing
import zipfile

from fastapi import FastAPI, Form, HTTPException, Request, Response, UploadFile
from fastapi.middleware import Middleware
from fastapi.responses import RedirectResponse, StreamingResponse
from fastapi.staticfiles import StaticFiles
from fastapi.templating import Jinja2Templates
from starlette.middleware.sessions import SessionMiddleware
from starlette.exceptions import HTTPException as StarletteHTTPException

import tvlink

app = FastAPI(
    docs_url=None,
    redoc_url=None,
    middleware=[
        Middleware(SessionMiddleware, secret_key=os.environ["SECRET_KEY"])
    ]
)
app.mount("/static", StaticFiles(directory="tvlink/static"), name="static")
templates = Jinja2Templates(directory="tvlink/templates")



@app.exception_handler(StarletteHTTPException)
async def http_exception_handler(request: Request, exc: StarletteHTTPException) -> Response:
    return templates.TemplateResponse(
        request=request,
        name="error.html",
        context={"exception": str(exc.detail)},
        status_code=exc.status_code
    )


@app.get("/")
async def frame(request: Request) -> Response:
    if "auth" not in request.session:
        return RedirectResponse("/LOGIN.XHTML", status_code=303)

    return templates.TemplateResponse(
        request=request,
        name="base.html"
    )


@app.get("/HEADER.XHTML")
async def header(request: Request) -> Response:
    return templates.TemplateResponse(
        request=request,
        name="header.html"
    )


@app.get("/NAVI.XHTML")
async def navi(request: Request) -> Response:
    return templates.TemplateResponse(
        request=request,
        name="navi.html"
    )


@app.get("/HELP.XHTML")
async def help(request: Request) -> Response:
    return templates.TemplateResponse(
        request=request,
        name="help.html"
    )


@app.get("/LOGIN.XHTML")
async def login_page(request: Request) -> Response:
    if "auth" in request.session:
        return RedirectResponse("/", status_code=303)

    return templates.TemplateResponse(
        request=request,
        name="login.html",
        context={}
    )


@app.post("/LOGIN.XHTML")
async def login(request: Request, login: typing.Annotated[str, Form()], password: typing.Annotated[str, Form()]) -> Response:
    if "auth" in request.session:
        return RedirectResponse("/", status_code=303)

    if login != "admin" or password != "admin":
        raise HTTPException(403, "Invalid credentials")

    request.session["auth"] = True
    return RedirectResponse("/", status_code=303)


@app.get("/LOGOUT.DO")
async def logout(request: Request) -> Response:
    del request.session["auth"]
    return RedirectResponse("/LOGIN.XHTML", status_code=303)


@app.get("/TV.XHTML")
async def tv(request: Request) -> Response:
    if "auth" not in request.session:
        return RedirectResponse("/LOGIN.XHTML", status_code=303)

    tv_enabled = False  # TODO
    return templates.TemplateResponse(
        request=request,
        name="tv.html",
        context={"tv_enabled": tv_enabled}
    )


@app.get("/UPGRADE.XHTML")
async def upgrade_page(request: Request) -> Response:
    if "auth" not in request.session:
        return RedirectResponse("/LOGIN.XHTML", status_code=303)

    return templates.TemplateResponse(
        request=request,
        name="upgrade.html",
        context={"version": tvlink.__version__}
    )


@app.post("/UPGRADE.XHTML")
async def upgrade(request: Request, firmware: UploadFile) -> Response:
    if "auth" not in request.session:
        return RedirectResponse("/LOGIN.XHTML", status_code=303)

    if not firmware.size:
        raise HTTPException(400, "Invalid firmware file")

    if firmware.size >= 128 * 1024:
        raise HTTPException(422, f"Not enough disk space")

    try:
        directory = tempfile.TemporaryDirectory()

        with zipfile.ZipFile(firmware.file) as archive:
            if archive.testzip() is not None:
                raise HTTPException(400, "Corrupt firmware")

            entries = archive.infolist()
            if len(entries) != 2:
                raise HTTPException(400, "Invalid firmware content")

            if not any(item.filename == "signature.txt" for item in entries):
                raise HTTPException(400, "Invalid firmware content")

            if any(item.file_size > 128 * 1024 for item in entries):
                raise HTTPException(422, "Not enough disk space")
            
            with archive.open("signature.txt") as signature_file:
                signature = signature_file.read()
                try:
                    signature = signature.decode()
                except UnicodeDecodeError:
                    raise HTTPException(400, "Invalid signature")
            
            fw_entry = next(item for item in entries if item.filename != "signature.txt")

            archive.extract(fw_entry, directory.name)

            firmware_script = pathlib.Path(directory.name) / fw_entry.filename

        hasher = hashlib.md5(os.environ["FIRMWARE_SECRET"].encode())
        hasher.update(firmware_script.read_bytes())
        expected_signature = hasher.hexdigest()
        if expected_signature != signature.strip():
            raise HTTPException(400, "Invalid signature")

        def log_upgrade_progress():
            yield templates.get_template("upgrade-status.html").render()
            try:
                upgrader = subprocess.Popen(
                    ["python", str(firmware_script)],
                    env={"PYTHONUNBUFFERED": "1", "PYTHONPATH": ":".join(sys.path), "PATH": os.environ["PATH"]},
                    bufsize=0,
                    stdout=subprocess.PIPE,
                    stderr=subprocess.STDOUT,
                    universal_newlines=True,
                    preexec_fn=lambda: signal.alarm(10)
                )

                if upgrader.stdout is not None:
                    yield from upgrader.stdout
                
                yield "</CODE></PRE><P>"
                
                if upgrader.wait(timeout=1) != 0:
                    yield f"Upgrade failed with code {upgrader.returncode}.</P><P><A href=\"/UPGRADE.XHTML\">Back</A></P>"
                else:
                    yield f"Upgrade succeeeded, reboot the device"
            except:
                upgrader.wait(timeout=1)
            finally:
                directory.cleanup()
        
        return StreamingResponse(
            log_upgrade_progress(),
            media_type="text/html",
            headers={
                "X-Accel-Buffering": "no",
                "X-Content-Type-Options": "nosniff"
            }
        )
    except Exception as exc:
        directory.cleanup()
        raise HTTPException(500, str(exc))


def debug():
    import uvicorn
    uvicorn.run("tvlink.app:app", reload=True, port=8000)
