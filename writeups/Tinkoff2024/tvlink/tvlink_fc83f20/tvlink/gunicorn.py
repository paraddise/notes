import multiprocessing

from gunicorn.app.base import BaseApplication

from tvlink.app import app


class TvlinkApplication(BaseApplication):
    def __init__(self):
        self.application = app
        super().__init__()


    def load_config(self):
        config = {
            "bind": "0.0.0.0:8000",
            "workers": min(multiprocessing.cpu_count() * 2, 8),
            "timeout": 180,
            "worker_class": "uvicorn.workers.UvicornWorker",
            "proc_name": "gunicorn: tvlink"
        }

        for key, value in config.items():
            if key in self.cfg.settings:
                self.cfg.set(key.lower(), value)


    def load(self):
        return self.application
