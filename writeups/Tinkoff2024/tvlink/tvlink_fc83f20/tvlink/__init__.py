__version__ = "20990131T070000-infected"


def run():
    from tvlink.gunicorn import TvlinkApplication
    TvlinkApplication().run()
