# -*- coding: utf-8 -*-

from flask_cors import CORS
from flask import Flask
from flask_cors import *
from api.v1 import apiv1


def _access_control(response):
    # response.headers['Access-Control-Allow-Origin'] = '*'
    response.headers['Access-Control-Allow-Methods'] = 'GET,HEAD,PUT,PATCH,POST,DELETE'
    response.headers['Access-Control-Allow-Headers'] = 'Content-Type'
    response.headers['Access-Control-Max-Age'] = 86400
    return response


def create_app():
    app = Flask(__name__)
    CORS(app, supports_credentials=True)
    app.after_request(_access_control)
    # add register blueprint here...
    app.register_blueprint(apiv1, url_prefix='/api/v1')

    return app


app = create_app()
