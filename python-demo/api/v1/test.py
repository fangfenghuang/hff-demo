from flask import Blueprint

apiv1 = Blueprint('apiv1', __name__)


@apiv1.route("/test", methods=['GET'])
def test():
    return {"code": 0, "data": "", "message": "helloworld!"}
 