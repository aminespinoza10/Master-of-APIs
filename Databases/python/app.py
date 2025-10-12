from fastapi import FastAPI, Depends, HTTPException, status
from fastapi.responses import PlainTextResponse
from fastapi.security import OAuth2PasswordBearer, OAuth2PasswordRequestForm
from pydantic import BaseModel
import os
import asyncio
import asyncpg
from typing import List
from dotenv import load_dotenv

load_dotenv()

app = FastAPI()
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")


class User(BaseModel):
    id: int
    name: str
    username: str
    password: str


class CreateUser(BaseModel):
    name: str
    username: str
    password: str


class UserPublic(BaseModel):
    id: int
    name: str
    username: str


@app.post("/token")
async def login(form_data: OAuth2PasswordRequestForm = Depends()):
    if form_data.username == "aminespinoza@mail.com" and form_data.password == "polliTeAmo123":
        return {"access_token": "yourtoken", "token_type": "bearer"}
    raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid credentials")


@app.get("/okCode", response_class=PlainTextResponse, name="getOkCode")
async def get_ok_code(token: str = Depends(oauth2_scheme)):
    return PlainTextResponse("Everything is awesome!", status_code=200)


@app.get("/getUsers", response_model=List[User])
async def get_users(token: str = Depends(oauth2_scheme)):
    db_url = os.getenv("DATABASE_URL")
    if not db_url:
        raise HTTPException(status_code=500, detail="DATABASE_URL not set")

    try:
        conn = await asyncpg.connect(db_url)
        rows = await conn.fetch("SELECT * FROM users")
        await conn.close()
    except Exception as e:
        raise HTTPException(status_code=500, detail="DB error")

    users = [User(id=row["id"], name=row["name"], username=row["username"], password=row.get("password")) for row in rows]
    return users


@app.post("/createUser", response_model=UserPublic)
async def create_user(payload: CreateUser, token: str = Depends(oauth2_scheme)):
    db_url = os.getenv("DATABASE_URL")
    if not db_url:
        raise HTTPException(status_code=500, detail="DATABASE_URL not set")

    try:
        conn = await asyncpg.connect(db_url)
        row = await conn.fetchrow("INSERT INTO users (name, username, password) VALUES ($1, $2, $3) RETURNING id, name, username", payload.name, payload.username, payload.password)
        await conn.close()
    except Exception as e:
        raise HTTPException(status_code=500, detail="DB error")

    return UserPublic(id=row["id"], name=row["name"], username=row["username"])



