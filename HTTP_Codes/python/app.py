from fastapi import FastAPI
from fastapi.responses import PlainTextResponse

app = FastAPI()

@app.get("/okCode", response_class=PlainTextResponse, name="getOkCode")
async def get_ok_code():
    return PlainTextResponse("Everything is awesome!", status_code=200)

@app.get("/continueCode", response_class=PlainTextResponse, name="getContinueCode")
async def get_continue_code():
    return PlainTextResponse("Continue processing...", status_code=100)

@app.get("/movedPermanently", response_class=PlainTextResponse, name="getMovedPermanently")
async def get_moved_permanently():
    return PlainTextResponse("This resource has been moved permanently.", status_code=301)

@app.get("/badRequest", response_class=PlainTextResponse, name="getBadRequest")
async def get_bad_request():
    return PlainTextResponse("Bad request. Please check your input.", status_code=400)

@app.get("/forbidden", response_class=PlainTextResponse, name="getForbidden")
async def get_forbidden():
    return PlainTextResponse("Access forbidden. You don't have permission to access this resource.", status_code=403)

@app.get("/notFound", response_class=PlainTextResponse, name="getNotFound")
async def get_not_found():
    return PlainTextResponse("Resource not found.", status_code=404)

@app.get("/proxyRequired", response_class=PlainTextResponse, name="getProxyRequired")
async def get_proxy_required():
    return PlainTextResponse("Proxy authentication required.", status_code=407)
