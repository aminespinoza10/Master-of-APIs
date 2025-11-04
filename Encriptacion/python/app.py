from fastapi import FastAPI
from dotenv import load_dotenv
from routers import users, auth

load_dotenv()

app = FastAPI()

app.include_router(users.router)
app.include_router(auth.router)

